package router

import (
	"taskmanager/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter - configures the routes
func SetupRouter() *gin.Engine {
	// a Gin router with dafault middleware
	r := gin.Default()
	api := r.Group("/tasks")
	{
		api.GET("", controllers.GetTasks)
		api.GET("/:id", controllers.GetTaskById)
		api.POST("", controllers.CreateTask)
		api.PUT("/:id", controllers.UpdateTask)
		api.DELETE("/:id", controllers.DeleteTask)
	}
	return r
}
