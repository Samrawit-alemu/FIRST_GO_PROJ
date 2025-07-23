package controllers

import (
	"net/http"
	"taskmanager/data"
	"taskmanager/models"

	"github.com/gin-gonic/gin"
)

// GetTasks handle GET requests to fetch all requests
func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTaskById handle GET requests to fetch a task by its id
func GetTaskById(c *gin.Context) {
	idstr := c.Param("id")

	task, err := data.GetTaskById(idstr)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)

}

// CreateTask handles POST requests to create a new task
func CreateTask(c *gin.Context) {
	var newTask models.Task

	// Bind the JSON to the struct
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	createdTask, err := data.CreateTask(newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask handles PUT requests to update a task.
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	task, err := data.UpdateTask(idStr, updatedTask)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE requests to remove a task
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")

	if err := data.DeleteTask(idStr); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
