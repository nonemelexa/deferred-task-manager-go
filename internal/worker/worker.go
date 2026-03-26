package worker

import (
	"fmt"
	"time"

	"task-scheduler/internal/domain"
)

type Queue interface {
	Consume() <-chan domain.Task
}

type Storage interface {
	Save(task domain.Task) error
	SaveFailed(task domain.Task) error
}

type Processor interface {
	Process(task domain.Task) error
}

type Worker struct {
	id        int
	queue     Queue
	storage   Storage
	processor Processor
}

// NewWorker creates and returns a new Worker instance with the given ID, queue, storage, and processor.
func NewWorker(id int, q Queue, s Storage, p Processor) *Worker {
	return &Worker{
		id:        id,
		queue:     q,
		storage:   s,
		processor: p,
	}
}

// Start begins consuming tasks from the queue and processing them in an infinite loop, handling success, retries with backoff, and failures according to the defined logic.
func (w *Worker) Start() {
	tasks := w.queue.Consume()

	for task := range tasks {
		fmt.Printf("[Worker %d] Processing task %s (type=%s)\n", w.id, task.ID, task.Type)
		w.handleTask(task)
	}
}

// handleTask processes a single task, implementing the logic for success, retries with exponential backoff, and marking tasks as failed after exceeding retry attempts.
func (w *Worker) handleTask(task domain.Task) {
	err := w.processor.Process(task)

	if err != nil {
		fmt.Printf("[Worker %d] Error: %s\n", w.id, err.Error())

		task.Attempt++

		if task.Attempt < task.Retries {
			delay := time.Duration(1<<task.Attempt) * time.Second
			task.ExecuteAt = time.Now().Add(delay)

			fmt.Printf("[Worker %d] Retry task %s in %v\n", w.id, task.ID, delay)

			w.storage.Save(task)
		} else {
			fmt.Printf("[Worker %d] Dead task %s\n", w.id, task.ID)

			w.storage.SaveFailed(task)
		}

		return
	}

	fmt.Printf("[Worker %d] Success task %s\n", w.id, task.ID)
}
