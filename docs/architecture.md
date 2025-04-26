# Architecture Overview

This document provides a comprehensive overview of the Paper.Social Feed Microservice architecture.

## System Architecture

Paper.Social Feed Service implements a microservices architecture with two main services:

```
                  ┌─────────────────────────┐
                  │      Orchestrator       │
                  │     (cmd/main.go)       │
                  └───────────┬─────────────┘
                              │
                              │ Manages
                              ▼
         ┌────────────────────┴───────────────────┐
         │                                        │
┌────────▼───────────┐                 ┌──────────▼────────┐
│                    │                 │                    │
│   Post Service     │ ◄────gRPC─────► │  GraphQL Service   │
│  (Internal API)    │                 │   (Public API)     │
│                    │                 │                    │
└────────┬───────────┘                 └──────────┬────────┘
         │                                        │
         │                                        │
┌────────▼───────────┐                 ┌──────────▼────────┐
│                    │                 │                    │
│   In-Memory DB     │                 │  GraphQL Playground│
│                    │                 │                    │
└────────────────────┘                 └────────────────────┘
```

### Key Components

1. **Orchestrator**
   - **Purpose**: Manages the lifecycle of both services
   - **Implementation**: `cmd/main.go`
   - **Features**:
     - Starts the Post Service first
     - Waits for it to initialize
     - Starts the GraphQL Service
     - Handles graceful shutdown

2. **Post Service (gRPC)**
   - **Purpose**: Internal service for post CRUD operations
   - **Implementation**: `postservice/server.go` and `postservice/cmd/main.go`
   - **Features**:
     - Provides gRPC endpoints for post management
     - Accesses the in-memory database
     - Not exposed externally (localhost only)
     - Runs on port 50051

3. **GraphQL Service**
   - **Purpose**: Public API for the feed service
   - **Implementation**: `graphqlservice/service.go` and `graphqlservice/cmd/main.go`
   - **Features**:
     - Provides GraphQL API for client applications
     - Communicates with the Post Service via gRPC
     - Exposed externally on port 8080
     - Includes GraphQL Playground for testing

4. **Shared Model**
   - **Purpose**: Common data structures and business logic
   - **Implementation**: `model/model.go`
   - **Features**:
     - User and Post data structures
     - In-memory database simulation
     - Image URL detection utility

## Request Flow

### Getting Timeline Example

```
┌─────────┐      ┌───────────────┐      ┌─────────────┐      ┌──────────┐
│ Client  │      │ GraphQL       │      │ Post        │      │ In-Memory│
│         │      │ Service       │      │ Service     │      │ Database │
└────┬────┘      └───────┬───────┘      └──────┬──────┘      └─────┬────┘
     │                   │                      │                   │
     │  GraphQL Request  │                      │                   │
     │  getTimeline()    │                      │                   │
     │  (userId: "...")  │                      │                   │
     │───────────────────►                      │                   │
     │                   │                      │                   │
     │                   │  1. Get user         │                   │
     │                   │─────────────────────────────────────────►│
     │                   │                      │                   │
     │                   │  2. Return user data │                   │
     │                   │◄─────────────────────────────────────────│
     │                   │                      │                   │
     │                   │  For each followed user:                 │
     │                   │  3. ListPostsByUser  │                   │
     │                   │─────────────────────►│                   │
     │                   │                      │  4. Get posts     │
     │                   │                      │──────────────────►│
     │                   │                      │                   │
     │                   │                      │  5. Return posts  │
     │                   │                      │◄──────────────────│
     │                   │  6. Posts            │                   │
     │                   │◄─────────────────────│                   │
     │                   │                      │                   │
     │                   │  7. Sort by time     │                   │
     │                   │  8. Limit to 20      │                   │
     │                   │  9. Extract images   │                   │
     │                   │                      │                   │
     │  GraphQL Response │                      │                   │
     │◄───────────────────                      │                   │
     │                   │                      │                   │
```

## Service Communication

1. **Client to GraphQL Service**
   - Protocol: HTTP/JSON
   - Endpoint: `/query`
   - Authentication: None (could be added as an improvement)

2. **GraphQL Service to Post Service**
   - Protocol: gRPC
   - Methods: ListPostsByUser, CreatePost, UpdatePost, DeletePost
   - Connection: localhost:50051

## Data Flow

1. **Timeline Aggregation**
   - GraphQL service fetches the user's "follows" list
   - For each followed user, it concurrently requests posts via gRPC
   - A mutex protects the aggregated posts collection
   - After all goroutines complete, posts are sorted by time
   - The most recent 20 posts are returned

2. **Post Operations**
   - Create/Update/Delete operations flow from GraphQL to gRPC service
   - The gRPC service interacts with the database
   - Results are converted back to GraphQL types

## Deployment Considerations

1. **Development**
   - Run using the orchestrator: `go run cmd/main.go`
   - Services can also be run individually for debugging

2. **Production**
   - Each service could be containerized separately
   - Service discovery would replace hardcoded addresses
   - Add proper authentication and authorization

## Future Architecture Extensions

1. **Authentication Service**
   - Handle user authentication and token generation
   - Integrate with GraphQL service via middleware

2. **Notification Service**
   - Push real-time updates to clients
   - Use GraphQL subscriptions for delivery

3. **Media Service**
   - Handle image uploads and processing
   - Store and serve media files 