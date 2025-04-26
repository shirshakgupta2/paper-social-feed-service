package main

import (
	"log"

	"github.com/paper-social/feed-service/model"
	"github.com/paper-social/feed-service/postservice"
)

func main() {
	// Create in-memory database with mock data
	db := model.NewDatabase()

	// Start the internal TCP server (not exposed externally)
	log.Println("Starting internal post service on port 50051...")
	if err := postservice.StartServer(db, ":50051"); err != nil {
		log.Fatalf("Failed to start post service: %v", err)
	}
}
