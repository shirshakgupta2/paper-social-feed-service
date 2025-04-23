package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	log.Println("Starting Paper.Social Feed Microservice System")

	// Get the current working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Determine if we're in the cmd directory or the project root
	dirName := filepath.Base(workDir)
	if dirName == "cmd" {
		// Move up one directory
		workDir = filepath.Dir(workDir)
		err = os.Chdir(workDir)
		if err != nil {
			log.Fatalf("Failed to change to project root directory: %v", err)
		}
	}

	// Start the post service (internal, not exposed externally)
	postServiceCmd := exec.Command("go", "run", filepath.Join(workDir, "postservice", "cmd", "main.go"))
	postServiceCmd.Stdout = os.Stdout
	postServiceCmd.Stderr = os.Stderr
	if err := postServiceCmd.Start(); err != nil {
		log.Fatalf("Failed to start post service: %v", err)
	}

	// Start the GraphQL service (external API)
	graphqlServiceCmd := exec.Command("go", "run", filepath.Join(workDir, "graphqlservice", "cmd", "main.go"))
	graphqlServiceCmd.Stdout = os.Stdout
	graphqlServiceCmd.Stderr = os.Stderr
	if err := graphqlServiceCmd.Start(); err != nil {
		log.Fatalf("Failed to start GraphQL service: %v", err)
	}

	log.Println("All services started successfully")
	log.Println("Post Service (internal) available at localhost:50051")
	log.Println("GraphQL API (public) available at http://localhost:8080/query")
	log.Println("GraphQL Playground available at http://localhost:8080")

	// Set up signal catching
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for Ctrl+C or kill signal
	<-signals
	log.Println("Shutting down services...")

	// Kill the services
	if err := postServiceCmd.Process.Kill(); err != nil {
		log.Printf("Error killing post service: %v", err)
	}
	if err := graphqlServiceCmd.Process.Kill(); err != nil {
		log.Printf("Error killing GraphQL service: %v", err)
	}

	log.Println("All services stopped")
}
