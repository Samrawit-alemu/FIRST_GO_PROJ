package usecases

import (
	"context"
	"errors"
	"taskmanager/domain"
	"taskmanager/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITaskUsecase interface {
	CreateTask(ctx context.Context, task *domain.Task, userID primitive.ObjectID) (*domain.Task, error)
	GetUserTasks(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error)
	GetTaskByID(ctx context.Context, taskID string, userID primitive.ObjectID) (*domain.Task, error)
	UpdateTask(ctx context.Context, taskID string, updatedTask *domain.Task, userID primitive.ObjectID) (*domain.Task, error)
	DeleteTask(ctx context.Context, taskID string, userID primitive.ObjectID) error
}

type taskUsecase struct {
	taskRepo repositories.ITaskRepository
}

func NewTaskUsecase(repo repositories.ITaskRepository) ITaskUsecase {
	return &taskUsecase{taskRepo: repo}
}

func (uc *taskUsecase) CreateTask(ctx context.Context, task *domain.Task, userID primitive.ObjectID) (*domain.Task, error) {
	task.UserID = userID
	err := uc.taskRepo.Create(ctx, task)
	return task, err
}

func (uc *taskUsecase) GetUserTasks(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error) {
	return uc.taskRepo.GetAllByUserID(ctx, userID)
}

func (uc *taskUsecase) GetTaskByID(ctx context.Context, taskID string, userID primitive.ObjectID) (*domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, errors.New("invalid task ID format")
	}

	task, err := uc.taskRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	if task.UserID != userID {
		return nil, errors.New("task not found")
	}

	return task, nil
}

func (uc *taskUsecase) UpdateTask(ctx context.Context, taskID string, updatedTask *domain.Task, userID primitive.ObjectID) (*domain.Task, error) {
	taskToUpdate, err := uc.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	taskToUpdate.Title = updatedTask.Title
	taskToUpdate.Description = updatedTask.Description
	taskToUpdate.Duedate = updatedTask.Duedate
	taskToUpdate.Status = updatedTask.Status

	err = uc.taskRepo.Update(ctx, taskToUpdate)
	if err != nil {
		return nil, err
	}
	return taskToUpdate, nil
}

func (uc *taskUsecase) DeleteTask(ctx context.Context, taskID string, userID primitive.ObjectID) error {
	taskToDelete, err := uc.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	return uc.taskRepo.Delete(ctx, taskToDelete.Id)
}
