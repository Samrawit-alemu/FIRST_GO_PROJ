package data

import (
	"context"
	"errors"
	"taskmanager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// In-memory database

// GetAllTasks - retrieve all tasks from the MongoDB collection
func GetAllTasks() ([]models.Task, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTaskById - retrieve a task by its id
func GetTaskById(id string) (models.Task, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("Invalid task Id format")
	}

	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

// CreateTask - create a new task
func CreateTask(task models.Task) (models.Task, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}
	task.Id = result.InsertedID.(primitive.ObjectID)
	return task, nil
}

// UpdateTask - update an existing task
func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"due_date":    updatedTask.Duedate,
		"status":      updatedTask.Status,
	}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Task{}, err
	}
	if result.MatchedCount == 0 {
		return models.Task{}, errors.New("task not found")
	}
	updatedTask.Id = objectID
	return updatedTask, nil
}

// DeleteTask - removes a task by its id
func DeleteTask(id string) error {
	collection := GetCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
