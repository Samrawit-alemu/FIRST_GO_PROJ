package controllers

import (
	"net/http"
	"taskmanager/delivery/dto"
	"taskmanager/domain"
	"taskmanager/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Promote(c *gin.Context)
}

type ITaskController interface {
	CreateTask(c *gin.Context)
	GetUserTasks(c *gin.Context)
	GetTaskByID(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

func toUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
	}
}

func toTaskResponse(task *domain.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.Duedate,
		Status:      task.Status,
		UserID:      task.UserID.Hex(),
	}
}

func toTasksResponse(tasks []domain.Task) []dto.TaskResponse {
	responses := make([]dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = toTaskResponse(&t)
	}
	return responses
}

// --- USER CONTROLLER ---
type UserController struct {
	userUsecase usecases.IUserUsecase
}

func NewUserController(userUsecase usecases.IUserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (uc *UserController) Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdUser, err := uc.userUsecase.Register(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map the result to our response DTO
	response := toUserResponse(createdUser)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": response})
}

func (uc *UserController) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, err := uc.userUsecase.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) Promote(c *gin.Context) {
	userID := c.Param("id")
	updatedUser, err := uc.userUsecase.Promote(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted", "user": toUserResponse(updatedUser)})
}

// --- TASK CONTROLLER ---
type TaskController struct {
	taskUsecase usecases.ITaskUsecase
}

func NewTaskController(taskUsecase usecases.ITaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var input dto.TaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	// Map the DTO to the Domain model
	domainTask := &domain.Task{
		Title:       input.Title,
		Description: input.Description,
		Duedate:     input.DueDate,
		Status:      input.Status,
	}

	createdTask, err := tc.taskUsecase.CreateTask(c.Request.Context(), domainTask, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	response := toTaskResponse(createdTask)
	c.JSON(http.StatusCreated, response)
}
func (tc *TaskController) GetUserTasks(c *gin.Context) {
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	tasks, err := tc.taskUsecase.GetUserTasks(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, toTasksResponse(tasks))
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	task, err := tc.taskUsecase.GetTaskByID(c.Request.Context(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toTaskResponse(task))
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var input dto.TaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	domainTask := &domain.Task{
		Title:       input.Title,
		Description: input.Description,
		Duedate:     input.DueDate,
		Status:      input.Status,
	}

	updatedTask, err := tc.taskUsecase.UpdateTask(c.Request.Context(), taskID, domainTask, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toTaskResponse(updatedTask))
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))
	err := tc.taskUsecase.DeleteTask(c.Request.Context(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
