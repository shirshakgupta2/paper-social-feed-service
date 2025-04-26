# Paper.Social Feed Microservice System

This project implements a simplified backend for a social feed system like Paper.Social. The architecture consists of two microservices:

1. **Post Service** (gRPC) - Internal service for post CRUD operations
2. **GraphQL Service** - Public API for accessing the social feed

## Project Structure

```
.
├── cmd/                  # Orchestrator to run both services
├── docs/                 # Documentation
│   ├── architecture.md   # System architecture overview
│   ├── code_generation.md # Code generation guide
│   ├── development.md    # Development workflow
│   ├── graphql_documentation.md # GraphQL API documentation
│   └── startup_guide.md  # Detailed startup instructions
├── graphqlservice/       # GraphQL service implementation
│   ├── cmd/              # GraphQL service entry point
│   ├── graph/            # GraphQL schema and resolvers
│   └── service.go        # Core service implementation
├── model/                # Shared data models
├── postservice/          # gRPC post service 
│   ├── cmd/              # Post service entry point
│   └── server.go         # gRPC server implementation
└── proto/                # Protocol Buffers definitions
    └── post/             # Post service protobuf
```

## Features

- **Timeline API**: Fetch the latest 20 posts from followed users, sorted by time
- **Post Management**: Create, read, update, and delete posts
- **Image URL Detection**: Automatically detect image URLs in post content
- **Parallel Processing**: Fetch posts from multiple users concurrently using goroutines
- **In-Memory Database**: Simulated database with mock data for testing
- **Orchestrated Startup**: Single command to start all services in the correct order
- **Port Availability Checking**: Prevents port conflicts when starting services

## Quick Start

### Option 1: Run both services with the orchestrator (Recommended)

```bash
go run cmd/main.go
```

This will start both the Post Service and the GraphQL Service in the correct order. The orchestrator:
- Checks if ports are available
- Starts the Post Service on port 50051 (internal)
- Waits for it to initialize
- Starts the GraphQL Service on port 8080 (public)
- Provides a nice URL display for accessing the playground

### Option 2: Run services individually

In one terminal, start the Post Service:

```bash
go run postservice/cmd/main.go
```

In another terminal, start the GraphQL Service:

```bash
go run graphqlservice/cmd/main.go
```

For detailed startup instructions, troubleshooting, and more, see [Startup Guide](docs/startup_guide.md).

## Accessing the API

- **GraphQL Playground**: http://localhost:8080
- **GraphQL Endpoint**: http://localhost:8080/query

## Example Queries

### Get Timeline

```graphql
query {
  getTimeline(userId: "user1") {
    id
    userId
    content
    createdAt
    imageUrls
  }
}
```

### Create Post

```graphql
mutation {
  createPost(userId: "user1", content: "This is a new post with an image: https://example.com/image.jpg") {
    id
    content
    createdAt
    imageUrls
  }
}
```

### Update Post

```graphql
mutation {
  updatePost(id: "post1", content: "Updated content") {
    id
    content
    imageUrls
  }
}
```

### Delete Post

```graphql
mutation {
  deletePost(id: "post1") {
    success
    message
  }
}
```

For complete API documentation, see [GraphQL Documentation](docs/graphql_documentation.md).

## Architecture

This project uses a microservices architecture with:

- **Orchestrator**: Manages service lifecycle
- **Post Service (gRPC)**: Handles post CRUD operations internally
- **GraphQL Service**: Provides public API and communicates with Post Service
- **Shared Model**: Common data structures and in-memory database

For a detailed architecture overview including flow diagrams, see [Architecture Documentation](docs/architecture.md).

## Development

For information on development workflow, code generation, and best practices, see:
- [Development Workflow](docs/development.md)
- [Code Generation Guide](docs/code_generation.md)

## Technical Implementation Details

- **Concurrency**: The timeline aggregation uses goroutines and mutex to fetch posts from followed users in parallel
- **Data Model**: In-memory simulation of a database with users, follower relationships, and posts
- **Error Handling**: Proper error propagation from the gRPC service to the GraphQL API
- **Image Detection**: Regular expressions to extract image URLs from post content
- **Port Management**: Dynamic port availability checking to prevent conflicts

## Further Improvements

- Add authentication and authorization
- Implement pagination for timeline results
- Add caching for frequently accessed timelines
- Create subscription for real-time updates to timeline
- Add test suite for both services
- Containerize services for production deployment 