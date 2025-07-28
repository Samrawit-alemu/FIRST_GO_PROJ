package controllers

import (
	"net/http"
	"taskmanager/domain"
	"taskmanager/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserController handles HTTP requests for user-related actions.
type UserController struct {
	userUsecase usecases.IUserUsecase
}

func NewUserController(userUsecase usecases.IUserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

// Register handles the POST /auth/register request.
func (uc *UserController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: username and password are required"})
		return
	}

	user, err := uc.userUsecase.Register(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

// Login handles the POST /auth/login request.
func (uc *UserController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	token, err := uc.userUsecase.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote handles the PUT /admin/promote/:id request.
func (uc *UserController) Promote(c *gin.Context) {
	userID := c.Param("id")
	updatedUser, err := uc.userUsecase.Promote(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted successfully", "user": updatedUser})
}

// TaskController handles HTTP requests for task-related actions.
type TaskController struct {
	taskUsecase usecases.ITaskUsecase
}

func NewTaskController(taskUsecase usecases.ITaskUsecase) *TaskController {
	return &TaskController{taskUsecase: taskUsecase}
}

// CreateTask handles the POST /tasks request.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var input domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	task, err := tc.taskUsecase.CreateTask(c.Request.Context(), &input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// GetUserTasks handles the GET /tasks request.
func (tc *TaskController) GetUserTasks(c *gin.Context) {
	// The middleware has already validated the token. We get the user's ID from the context.
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	// Call the use case to get the tasks for this user.
	tasks, err := tc.taskUsecase.GetUserTasks(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID handles the GET /tasks/:id request.
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	task, err := tc.taskUsecase.GetTaskByID(c.Request.Context(), taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UpdateTask handles the PUT /tasks/:id request.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var input domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userIDHex, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDHex.(string))

	updatedTask, err := tc.taskUsecase.UpdateTask(c.Request.Context(), taskID, &input, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask handles the DELETE /tasks/:id request.
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
