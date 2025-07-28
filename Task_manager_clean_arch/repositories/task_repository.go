package repositories

import (
	"context"
	"taskmanager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type mongoTaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) ITaskRepository {
	return &mongoTaskRepository{collection: db.Collection("tasks")}
}

func (r *mongoTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	task.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoTaskRepository) GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []domain.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *mongoTaskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var task domain.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *mongoTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	filter := bson.M{"_id": task.Id}
	update := bson.M{"$set": bson.M{
		"title":       task.Title,
		"description": task.Description,
		"due_date":    task.Duedate,
		"status":      task.Status,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoTaskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
