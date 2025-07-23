package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task-struct : the data structure for a task
// The json tags are used by Gin to map map the json keys from requests to the struct fields
type Task struct {
	// Use primitive.ObjectID for the _id field and add bson tags
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson: "title"`
	Description string             `json:"description" bson:"description"`
	Duedate     time.Time          `json:"due_date" bson:"due_date"`
	Status      string             `json:"status" bson:"status"`
}
