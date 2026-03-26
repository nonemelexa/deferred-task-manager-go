package api

import (
	"encoding/json"
	"net/http"
	"time"

	"task-scheduler/internal/domain"

	"github.com/google/uuid"
)

type Queue interface {
	Publish(task domain.Task) error
}

type Storage interface {
	Save(task domain.Task) error
}

type Handler struct {
	queue   Queue
	storage Storage
}

// NewHandler creates and returns a new instance of Handler with the provided queue and storage, which will be used to handle incoming API requests for creating tasks and managing their execution.
func NewHandler(q Queue, s Storage) *Handler {
	return &Handler{q, s}
}

// writeJSONError is a helper function to write JSON error responses with a consistent structure and appropriate HTTP status code.
func writeJSONError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

// isValidType is a helper function that checks if the provided task type is valid (i.e., "email", "payment", or "report") and returns true if it is, or false otherwise.
func isValidType(t string) bool {
	switch t {
	case "email", "payment", "report":
		return true
	default:
		return false
	}
}

// CreateTask is an HTTP handler that processes incoming requests to create new tasks. It validates the JSON payload, checks the task type, generates a unique ID, and either publishes the task to the queue for immediate execution or saves it in storage for delayed execution based on the specified delay.
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task domain.Task

	// check JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// check type
	if !isValidType(task.Type) {
		writeJSONError(w, "invalid task type", http.StatusBadRequest)
		return
	}

	// generate unique ID
	task.ID = uuid.New().String()

	task.ExecuteAt = time.Now().Add(time.Duration(task.Delay) * time.Second)

	if task.Delay == 0 {
		err := h.queue.Publish(task)
		if err != nil {
			writeJSONError(w, "failed to publish task", http.StatusInternalServerError)
			return
		}
	} else {

		err := h.storage.Save(task)
		if err != nil {
			writeJSONError(w, "failed to save task", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "task created",
		"id":     task.ID,
	})
}
