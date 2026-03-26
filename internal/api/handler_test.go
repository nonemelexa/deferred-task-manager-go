package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-scheduler/internal/domain"
)

type mockQueue struct {
	published bool
}

func (m *mockQueue) Publish(task domain.Task) error {
	m.published = true
	return nil
}

type mockStorage struct {
	saved bool
}

func (m *mockStorage) Save(task domain.Task) error {
	m.saved = true
	return nil
}

// 🧪 TEST: SUCCESS
func TestCreateTask_Immediate(t *testing.T) {
	q := &mockQueue{}
	s := &mockStorage{}

	handler := NewHandler(q, s)

	body := []byte(`{
		"type": "email",
		"payload": "test@test.com",
		"retries": 3,
		"delay": 0,
		"priority": 1
	}`)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if !q.published {
		t.Errorf("expected task to be published to queue")
	}
}

// 🧪 TEST: DELAYED
func TestCreateTask_Delayed(t *testing.T) {
	q := &mockQueue{}
	s := &mockStorage{}

	handler := NewHandler(q, s)

	body := []byte(`{
		"type": "email",
		"payload": "test@test.com",
		"retries": 3,
		"delay": 10,
		"priority": 1
	}`)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if !s.saved {
		t.Errorf("expected task to be saved in storage")
	}
}

// 🧪 TEST: INVALID JSON
func TestCreateTask_InvalidJSON(t *testing.T) {
	handler := NewHandler(&mockQueue{}, &mockStorage{})

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("bad json")))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// 🧪 TEST: INVALID TYPE
func TestCreateTask_InvalidType(t *testing.T) {
	handler := NewHandler(&mockQueue{}, &mockStorage{})

	body := []byte(`{
		"type": "unknown",
		"payload": "test",
		"delay": 0
	}`)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid type, got %d", w.Code)
	}
}
