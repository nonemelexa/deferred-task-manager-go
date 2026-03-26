package worker

import (
	"fmt"
	"testing"
	"time"

	"task-scheduler/internal/domain"
)

//
// 🔧 MOCK STORAGE
//

type mockStorage struct {
	saved       bool
	savedFailed bool
	lastTask    domain.Task
}

func (m *mockStorage) Save(task domain.Task) error {
	m.saved = true
	m.lastTask = task
	return nil
}

func (m *mockStorage) SaveFailed(task domain.Task) error {
	m.savedFailed = true
	m.lastTask = task
	return nil
}

//
// 🔧 MOCK PROCESSOR
//

type mockProcessor struct {
	shouldFail bool
}

func (m *mockProcessor) Process(task domain.Task) error {
	if m.shouldFail {
		return fmt.Errorf("forced error")
	}
	return nil
}

//
// 🧪 TEST: SUCCESS
//

func TestWorker_HandleTask_Success(t *testing.T) {
	ms := &mockStorage{}
	mp := &mockProcessor{shouldFail: false}

	w := &Worker{
		id:        1,
		storage:   ms,
		processor: mp,
	}

	task := domain.Task{
		ID:      "1",
		Type:    "email",
		Retries: 3,
		Attempt: 0,
	}

	w.handleTask(task)

	if ms.saved {
		t.Errorf("expected no retry save on success")
	}

	if ms.savedFailed {
		t.Errorf("expected no failed save on success")
	}
}

//
// 🧪 TEST: RETRY
//

func TestWorker_HandleTask_Retry(t *testing.T) {
	ms := &mockStorage{}
	mp := &mockProcessor{shouldFail: true}

	w := &Worker{
		id:        1,
		storage:   ms,
		processor: mp,
	}

	task := domain.Task{
		ID:      "1",
		Type:    "email",
		Retries: 3,
		Attempt: 0,
	}

	w.handleTask(task)

	if !ms.saved {
		t.Errorf("expected task to be saved for retry")
	}

	if ms.savedFailed {
		t.Errorf("did not expect task to be marked as failed")
	}

	if ms.lastTask.Attempt != 1 {
		t.Errorf("expected attempt to increment, got %d", ms.lastTask.Attempt)
	}
}

//
// 🧪 TEST: DEAD TASK
//

func TestWorker_HandleTask_Dead(t *testing.T) {
	ms := &mockStorage{}
	mp := &mockProcessor{shouldFail: true}

	w := &Worker{
		id:        1,
		storage:   ms,
		processor: mp,
	}

	task := domain.Task{
		ID:      "1",
		Type:    "email",
		Retries: 1,
		Attempt: 1,
	}

	w.handleTask(task)

	if !ms.savedFailed {
		t.Errorf("expected task to be marked as failed")
	}

	if ms.saved {
		t.Errorf("did not expect retry save")
	}
}

//
// 🧪 TEST: BACKOFF
//

func TestWorker_HandleTask_Backoff(t *testing.T) {
	ms := &mockStorage{}
	mp := &mockProcessor{shouldFail: true}

	w := &Worker{
		id:        1,
		storage:   ms,
		processor: mp,
	}

	task := domain.Task{
		ID:      "1",
		Type:    "email",
		Retries: 3,
		Attempt: 1,
	}

	before := time.Now()

	w.handleTask(task)

	if !ms.saved {
		t.Errorf("expected retry")
	}

	if ms.lastTask.ExecuteAt.Before(before) {
		t.Errorf("expected ExecuteAt to be in the future")
	}
}
