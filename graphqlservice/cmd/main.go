package main

import (
	"log"

	"github.com/paper-social/feed-service/graphqlservice"
	"github.com/paper-social/feed-service/model"
)

func main() {
	// Create in-memory database with mock data
	db := model.NewDatabase()

	// Connect to the internal post service
	postServiceAddr := "localhost:50051"
	log.Printf("Connecting to internal post service at %s", postServiceAddr)

	// Start the public GraphQL service
	log.Println("Starting public GraphQL API...")
	if err := graphqlservice.StartService(db, postServiceAddr, ":8080"); err != nil {
		log.Fatalf("Failed to start GraphQL service: %v", err)
	}
}
