package scheduler

import (
	"log"
	"time"

	"task-scheduler/internal/queue"
	"task-scheduler/internal/storage"
)

type Scheduler struct {
	storage *storage.RedisStorage
	queue   *queue.RabbitMQ
}

// NewScheduler creates a new Scheduler with the given storage and queue, which will be used to manage delayed tasks and publish them when they are ready for execution.
func NewScheduler(s *storage.RedisStorage, q *queue.RabbitMQ) *Scheduler {
	return &Scheduler{s, q}
}

// Start runs the scheduler in an infinite loop, periodically checking for tasks that are ready to be executed (i.e., their scheduled execution time has arrived) and publishing them to the queue for workers to consume.
func (s *Scheduler) Start() {
	for {
		tasks, err := s.storage.GetReadyTasks(time.Now())
		if err != nil {
			log.Println("Scheduler error:", err)
			continue
		}

		for _, task := range tasks {
			s.queue.Publish(task)
		}

		time.Sleep(2 * time.Second)
	}
}
