package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID
	Username string
	Password string // Hashed password
	Role     string
}

type Task struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	Duedate     time.Time
	Status      string
	UserID      primitive.ObjectID
}
