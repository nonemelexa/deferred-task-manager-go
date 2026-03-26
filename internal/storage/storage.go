package storage

import (
	"sync"
	"time"

	"task-scheduler/internal/domain"
)

// Storage is an in-memory storage for tasks, using a mutex to ensure thread safety when saving and retrieving tasks.
type Storage struct {
	mu    sync.Mutex
	tasks []domain.Task
}

// NewStorage creates and returns a new instance of Storage with an initialized task slice.
func NewStorage() *Storage {
	return &Storage{
		tasks: []domain.Task{},
	}
}

// Save adds a task to the storage, locking the mutex to ensure thread safety while modifying the task slice.
func (s *Storage) Save(task domain.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = append(s.tasks, task)
}

// GetReadyTasks retrieves tasks that are ready to be executed based on the current time, removing them from storage and returning them as a slice.
func (s *Storage) GetReadyTasks(now time.Time) []domain.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var ready []domain.Task
	var remaining []domain.Task

	for _, t := range s.tasks {
		if t.ExecuteAt.Before(now) || t.ExecuteAt.Equal(now) {
			ready = append(ready, t)
		} else {
			remaining = append(remaining, t)
		}
	}

	s.tasks = remaining
	return ready
}
