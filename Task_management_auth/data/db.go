package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

const DBNAME = "taskmanagerdb"
const TASKS_COLLECTION = "tasks"
const USERS_COLLECTION = "users"

// MONGO_URI is the connection string.
const MONGO_URI = "mongodb://localhost:27017"

// ConnectDB establishes a connection to the MongoDB server and initializes the client.
// This function should be called once when the application starts.
func ConnectDB() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	log.Println("Successfully connected to MongoDB!")
}

// GetTasksCollection is a helper function that returns a handle to the "tasks" collection.
func GetTasksCollection() *mongo.Collection {
	return mongoClient.Database(DBNAME).Collection(TASKS_COLLECTION)
}

// GetUsersCollection is a helper function that returns a handle to the "users" collection.
func GetUsersCollection() *mongo.Collection {
	return mongoClient.Database(DBNAME).Collection(USERS_COLLECTION)
}
