package main

import (
	"log"
	"taskmanager/router"
)

func main() {
	// setup the router
	r := router.SetupRouter()

	// start the server
	log.Println("Starting server on:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
