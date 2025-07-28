package domain

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
	UserID      primitive.ObjectID `json:"user_id" bson: "user_id"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username,unique"`
	Password string             `json:"-" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}
