package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson: "title"`
	Description string             `json:"description" bson:"description"`
	Duedate     time.Time          `json:"due_date" bson:"due_date"`
	Status      string             `json:"status" bson:"status"`
	//Id of the user who created the task
	UserID primitive.ObjectID `json:"user_id" bson: "user_id"`
}
