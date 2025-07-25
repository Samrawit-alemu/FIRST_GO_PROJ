package data

import (
	"context"
	"errors"
	"taskmanager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTask requires the ID of the user creating the task.
func CreateTask(task models.Task, userID primitive.ObjectID) (models.Task, error) {
	collection := GetTasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.UserID = userID

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	task.Id = result.InsertedID.(primitive.ObjectID)
	return task, nil
}

// GetAllTasks only gets tasks for a specific user.
func GetAllTasks(userID primitive.ObjectID) ([]models.Task, error) {
	collection := GetTasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}

	cursor, err := collection.Find(ctx, filter)
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

// GetTaskByID must check for ownership.
func GetTaskByID(taskID string, userID primitive.ObjectID) (models.Task, error) {
	collection := GetTasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	// Filter by BOTH the task ID and the user ID to ensure ownership.
	filter := bson.M{"_id": objectID, "userid": userID}

	var task models.Task
	err = collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.Task{}, errors.New("task not found or not owned by user")
	}
	return task, nil
}

// UpdateTask must check for ownership.
func UpdateTask(taskID string, updatedTask models.Task, userID primitive.ObjectID) (models.Task, error) {
	collection := GetTasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	// Filter by both ID and owner.
	filter := bson.M{"_id": objectID, "userid": userID}
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
		return models.Task{}, errors.New("task not found or not owned by user")
	}

	updatedTask.Id = objectID
	updatedTask.UserID = userID
	return updatedTask, nil
}

// DeleteTask must also check for ownership.
func DeleteTask(taskID string, userID primitive.ObjectID) error {
	collection := GetTasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	// Filter by both ID and owner.
	filter := bson.M{"_id": objectID, "userid": userID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found or not owned by user")
	}
	return nil
}
