package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DBName = "taskmanager"
const CollectionName = "tasks"

// ConnectToDB - establishes a connection to the MongoDB server
func ConnectToDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("could not ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB")
}

// GetCollection - returns the tasks collection
func GetCollection() *mongo.Collection {
	return client.Database(DBName).Collection(CollectionName)

}
