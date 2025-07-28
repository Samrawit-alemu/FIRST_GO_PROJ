package repositories

import (
	"context"
	"taskmanager/domain"
	datamodels "taskmanager/repositories/models" // Aliased import

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskRepository interface definition remains the same.
type ITaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// mongoTaskRepository is the concrete implementation.
type mongoTaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository is the constructor.
func NewTaskRepository(db *mongo.Database) ITaskRepository {
	return &mongoTaskRepository{collection: db.Collection("tasks")}
}

// toBsonTask converts a Domain Task to a BSON Task model.
func toBsonTask(task *domain.Task) *datamodels.Task {
	return &datamodels.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.Duedate,
		Status:      task.Status,
		UserID:      task.UserID,
	}
}

// toDomainTask converts a BSON Task model to a Domain Task.
func toDomainTask(task *datamodels.Task) *domain.Task {
	return &domain.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Duedate:     task.DueDate,
		Status:      task.Status,
		UserID:      task.UserID,
	}
}

// toDomainTasks converts a slice of BSON Task models to a slice of Domain Tasks.
func toDomainTasks(tasks []datamodels.Task) []domain.Task {
	domainTasks := make([]domain.Task, len(tasks))
	for i, t := range tasks {
		domainTasks[i] = *toDomainTask(&t)
	}
	return domainTasks
}

func (r *mongoTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	bsonTask := toBsonTask(task)
	result, err := r.collection.InsertOne(ctx, bsonTask)
	if err != nil {
		return err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mongoTaskRepository) GetAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error) {
	var bsonTasks []datamodels.Task
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bsonTasks); err != nil {
		return nil, err
	}
	return toDomainTasks(bsonTasks), nil
}

func (r *mongoTaskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var bsonTask datamodels.Task
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bsonTask)
	if err != nil {
		return nil, err
	}
	return toDomainTask(&bsonTask), nil
}

func (r *mongoTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	bsonTask := toBsonTask(task)
	filter := bson.M{"_id": bsonTask.ID}
	update := bson.M{"$set": bsonTask}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoTaskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
