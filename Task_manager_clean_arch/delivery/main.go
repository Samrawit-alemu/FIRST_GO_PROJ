package main

import (
	"context"
	"log"
	"taskmanager/delivery/controllers"
	"taskmanager/delivery/routers"
	"taskmanager/infrastructure"
	"taskmanager/repositories"
	"taskmanager/usecases"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	// --- DATABASE CONNECTION ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)
	log.Println("Connected to MongoDB!")
	db := client.Database("taskmanager_clean")

	// --- DEPENDENCY INJECTION (WIRING THE LAYERS TOGETHER) ---
	// Layer 4: Infrastructure (The Tools)
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService()

	// Layer 3: Repositories (The Database Implementations)
	userRepo := repositories.NewUserRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	// Layer 2: Usecases (The Business Logic)
	userUsecase := usecases.NewUserUsecase(userRepo, passwordService, jwtService)
	taskUsecase := usecases.NewTaskUsecase(taskRepo)

	// Layer 1: Delivery (The HTTP Handlers)
	userController := controllers.NewUserController(userUsecase)
	taskController := controllers.NewTaskController(taskUsecase)

	// --- SETUP ROUTER AND START SERVER ---
	router := routers.SetupRouter(userController, taskController, jwtService)
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
