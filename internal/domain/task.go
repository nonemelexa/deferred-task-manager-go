package domain

import "time"

// Task represents a unit of work to be executed, containing fields for ID, type, payload, retries, attempt count, delay, priority, and execution time.
type Task struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Retries   int       `json:"retries"`
	Attempt   int       `json:"attempt"`
	Delay     int       `json:"delay"`
	Priority  int       `json:"priority"`
	ExecuteAt time.Time `json:"execute_at"`
}
