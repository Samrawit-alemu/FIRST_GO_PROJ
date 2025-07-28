package routers

import (
	"taskmanager/delivery/controllers"
	"taskmanager/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *controllers.UserController, taskController *controllers.TaskController, jwtService infrastructure.IJWTService) *gin.Engine {
	r := gin.Default()

	// Public routes for authentication
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userController.Register)
		authRoutes.POST("/login", userController.Login)
	}

	// Protected routes that require a valid token
	protected := r.Group("")
	protected.Use(infrastructure.AuthMiddleware(jwtService))
	{
		// Task routes, accessible to all logged-in users
		taskRoutes := protected.Group("/tasks")
		{
			taskRoutes.GET("", taskController.GetUserTasks)
			taskRoutes.GET("/:id", taskController.GetTaskByID)

			// Admin-only task routes
			taskRoutes.POST("", infrastructure.RoleAuthMiddleware("admin"), taskController.CreateTask)
			taskRoutes.PUT("/:id", infrastructure.RoleAuthMiddleware("admin"), taskController.UpdateTask)
			taskRoutes.DELETE("/:id", infrastructure.RoleAuthMiddleware("admin"), taskController.DeleteTask)
		}

		// Admin-only management routes
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(infrastructure.RoleAuthMiddleware("admin"))
		{
			adminRoutes.PUT("/promote/:id", userController.Promote)
		}
	}

	return r
}
