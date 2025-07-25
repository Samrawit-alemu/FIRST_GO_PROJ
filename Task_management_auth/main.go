package main

import (
	"log"
	"taskmanager/data"
	"taskmanager/router"
)

func main() {
	// Connect to the database when the application starts
	data.ConnectDB()

	data.EnsureUsernameUnique()

	// setup the router
	r := router.SetupRouter()

	log.Println("Application setup complete. (Router to be added next)")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
