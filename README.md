# Paper.Social Feed Microservice

A feed microservice for a social media platform that displays users' timelines using Go, GraphQL, and gRPC.

## Architecture

This project implements a backend microservice system with two main components:

1. **Post Service (Internal)**: A gRPC service that manages posts
   - Not directly exposed to clients
   - Handles CRUD operations for posts
   - Uses an in-memory database for storing posts

2. **GraphQL Service (Public-facing)**: A GraphQL API that serves as the entry point for clients
   - Exposed at http://localhost:8080/query
   - Provides a GraphQL playground at http://localhost:8080
   - Communicates with the Post Service to fetch and manage post data

The architecture follows a unidirectional flow where:
- Clients interact only with the GraphQL API
- The GraphQL service communicates with the Post service
- The Post service manages the data store

## Features

- Get a user's timeline (aggregated posts from followed users)
- Create new posts
- Update existing posts
- Delete posts
- Concurrent post fetching using goroutines
- Post content with image URL detection
- Server architecture diagram detection

## Getting Started

### Prerequisites

- Go 1.20 or later
- [Protocol Buffers compiler](https://grpc.io/docs/protoc-installation/) (for development)

### Running the Application

```bash
go run cmd/main.go
```

This starts:
- Post Service on `localhost:50051` (internal)
- GraphQL API at `http://localhost:8080/query`
- GraphQL Playground at `http://localhost:8080`

## API Usage

### GraphQL Playground

The easiest way to interact with the API is through the GraphQL Playground at http://localhost:8080.

### Using with Postman

To use the API with Postman:

1. Create a new request to `http://localhost:8080/query`
2. Set the method to `POST`
3. Set the Content-Type header to `application/json`
4. In the body tab, select "raw" and "JSON"
5. Enter a GraphQL query:

```json
{
  "query": "{ getTimeline(userId: \"user1\") { id userId content createdAt } }",
  "variables": {}
}
```

#### Example Queries

**Get Timeline**
```json
{
  "query": "{ getTimeline(userId: \"user1\") { id userId content createdAt } }",
  "variables": { "userId": "user1" }
}
```

**Create Post**
```json
{
  "query": "mutation { createPost(userId: \"user1\", content: \"This is a new post!\") { id userId content createdAt } }",
  "variables": { "userId": "user1", "content": "This is a new post!" }
}
```

**Update Post**
```json
{
  "query": "mutation { updatePost(id: \"post1\", content: \"Updated content\") { id userId content createdAt } }",
  "variables": { "id": "post1", "content": "Updated content" }
}
```

**Delete Post**
```json
{
  "query": "mutation { deletePost(id: \"post1\") { success message } }",
  "variables": { "id": "post1" }
}
```

## Data Model

The system uses an in-memory database with mock data:

- 5 users: alice, bob, charlie, dave, and eve
- 10 initial posts with various content
- Follow relationships between users

## Implementation Details

### Post Service (gRPC)

The Post Service implements the following RPC methods:
- `ListPostsByUser`: Retrieves posts for a specific user
- `CreatePost`: Creates a new post
- `UpdatePost`: Updates an existing post
- `DeletePost`: Deletes a post

### GraphQL Service

The GraphQL service provides:
- A query for retrieving a user's timeline
- Mutations for creating, updating, and deleting posts
- A computed field for extracting image URLs from post content
- A computed field for detecting server architecture diagrams

## Design Decisions

1. **Microservice Architecture**: Separates concerns between post management and API exposure
2. **Unidirectional Communication Flow**: GraphQL → gRPC → Database
3. **Concurrent Post Fetching**: Uses goroutines for parallel data fetching
4. **In-Memory Database**: Simulates a data store for demonstration purposes 