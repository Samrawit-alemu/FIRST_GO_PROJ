package models

import (
	"time"
)

// Task-struct : the data structure for a task
// The json tags are used by Gin to map map the json keys from requests to the struct fields
type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duedate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}
