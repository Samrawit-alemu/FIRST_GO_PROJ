package usecases

import (
	"context"
	"taskmanager/domain"
	"taskmanager/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask_Success(t *testing.T) {
	mockTaskRepo := new(mocks.ITaskRepository)
	userID := primitive.NewObjectID()

	taskToCreate := &domain.Task{
		Title:  "New Task",
		Status: "Pending",
	}

	mockTaskRepo.On("Create", mock.Anything, mock.MatchedBy(func(task *domain.Task) bool {
		return task.UserID == userID && task.Title == "New Task"
	})).Return(nil)

	usecase := NewTaskUsecase(mockTaskRepo)
	createdTask, err := usecase.CreateTask(context.Background(), taskToCreate, userID)

	// --- ASSERT ---
	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.Equal(t, userID, createdTask.UserID)
	mockTaskRepo.AssertExpectations(t)
}

func TestGetTaskByID_Success_OwnerMatch(t *testing.T) {
	mockTaskRepo := new(mocks.ITaskRepository)
	taskID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	fakeTask := &domain.Task{
		ID:     taskID,
		Title:  "My Task",
		UserID: userID,
	}

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(fakeTask, nil)

	usecase := NewTaskUsecase(mockTaskRepo)
	foundTask, err := usecase.GetTaskByID(context.Background(), taskID.Hex(), userID)

	// --- ASSERT ---
	assert.NoError(t, err)
	assert.NotNil(t, foundTask)
	assert.Equal(t, fakeTask.Title, foundTask.Title)
	mockTaskRepo.AssertExpectations(t)
}

func TestGetTaskByID_Failure_OwnerMismatch(t *testing.T) {
	mockTaskRepo := new(mocks.ITaskRepository)
	taskID := primitive.NewObjectID()
	ownerUserID := primitive.NewObjectID()
	requesterUserID := primitive.NewObjectID()

	fakeTask := &domain.Task{
		ID:     taskID,
		Title:  "Someone Else's Task",
		UserID: ownerUserID,
	}

	mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(fakeTask, nil)
	usecase := NewTaskUsecase(mockTaskRepo)
	foundTask, err := usecase.GetTaskByID(context.Background(), taskID.Hex(), requesterUserID)

	// --- ASSERT ---
	assert.Error(t, err)
	assert.Nil(t, foundTask)
	assert.Equal(t, "task not found", err.Error())
	mockTaskRepo.AssertExpectations(t)
}
