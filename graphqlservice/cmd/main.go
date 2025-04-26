package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/paper-social/feed-service/graphqlservice"
	"github.com/paper-social/feed-service/graphqlservice/graph"
	"github.com/paper-social/feed-service/graphqlservice/graph/generated"
	"github.com/paper-social/feed-service/model"
)

func main() {
	// Create in-memory database with mock data
	db := model.NewDatabase()

	// Connect to the internal post service
	postServiceAddr := "localhost:50051"
	log.Printf("Connecting to internal post service at %s", postServiceAddr)

	// Create service
	service := graphqlservice.NewService(db, postServiceAddr)

	// Set up GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			Service: service,
		},
	}))

	// Set up CORS middleware
	setupCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	// Register the GraphQL playground at the root
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// Register the GraphQL query handler with CORS support
	http.Handle("/query", setupCORS(srv))

	// Start the public GraphQL service
	log.Println("Starting public GraphQL API...")
	log.Printf("GraphQL service starting on %s", ":8080")
	log.Printf("Connect to the GraphQL playground: http://localhost%s", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start GraphQL service: %v", err)
	}
}
