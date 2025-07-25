package router

import (
	"taskmanager/controllers"
	"taskmanager/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the routes for the entire application.
func SetupRouter() *gin.Engine {
	// Create a new Gin router with default middleware
	r := gin.Default()

	// --- Public Routes ---
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.RegisterUser)
		authRoutes.POST("/login", controllers.LoginUser)
	}

	// --- Protected Task Routes ---
	taskRoutes := r.Group("/tasks")
	taskRoutes.Use(middleware.AuthMiddleware())
	{
		// Only admins can create, update, and delete tasks.
		taskRoutes.POST("", middleware.RoleAuthMiddleware("admin"), controllers.CreateTask)
		taskRoutes.PUT("/:id", middleware.RoleAuthMiddleware("admin"), controllers.UpdateTask)
		taskRoutes.DELETE("/:id", middleware.RoleAuthMiddleware("admin"), controllers.DeleteTask)

		// Any authenticated user (admin or regular) can get tasks.
		taskRoutes.GET("", controllers.GetTasks)
		taskRoutes.GET("/:id", controllers.GetTaskByID)
	}

	// --- Admin-Only Routes ---
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RoleAuthMiddleware("admin"))
	{
		adminRoutes.PUT("/promote/:id", controllers.PromoteUser)
	}

	return r
}
