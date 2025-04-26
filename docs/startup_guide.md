# Startup Guide

This guide provides detailed instructions for starting the Paper.Social Feed Service.

## Prerequisites

- Go 1.22 or later
- Git
- A terminal or command prompt

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/paper-social/feed-service.git
cd feed-service
```

### 2. Install Dependencies

```bash
go mod download
```

## Running the Services

### Option 1: Using the Orchestrator (Recommended)

The orchestrator will start both services in the correct order and manage their lifecycle:

```bash
go run cmd/main.go
```

This will:
1. Start the Post Service on port 50051 (internal gRPC API)
2. Wait for it to initialize (2 seconds by default)
3. Start the GraphQL Service on port 8080 (public API)
4. Show the URLs for accessing the API and playground

**Success Output:**
```
Starting Paper.Social Feed Microservice System
Post service starting... waiting 2 seconds before starting GraphQL service
All services started successfully
Post Service (internal) available at localhost:50051
GraphQL API (public) available at http://localhost:8080/query
GraphQL Playground available at http://localhost:8080

----------------------------------------------------------
🚀 Access GraphQL Playground: http://localhost:8080
📋 Try the example queries in the README.md file
⚠️ Press Ctrl+C to stop all services
----------------------------------------------------------
```

### Option 2: Running Services Individually

For development or debugging, you can run each service separately.

#### 1. Start the Post Service (Terminal 1)

```bash
go run postservice/cmd/main.go
```

Output:
```
Starting internal post service on port 50051...
Post service gRPC server starting on :50051 (internal only)
```

#### 2. Start the GraphQL Service (Terminal 2)

```bash
go run graphqlservice/cmd/main.go
```

Output:
```
Connecting to internal post service at localhost:50051
Starting public GraphQL API...
GraphQL service starting on :8080
Connect to the GraphQL playground: http://localhost:8080
```

## Troubleshooting

### Port Already in Use

If you see an error like:

```
Port :50051 is already in use. Please stop any existing instances before running this command.
```

Use these commands to resolve:

**Windows:**
```bash
taskkill /F /IM go.exe
```

**Linux/Mac:**
```bash
pkill -f "go run"
```

Or find and kill the specific process:

**Windows:**
```bash
# Find the process
netstat -ano | findstr :50051
# Kill it (replace 12345 with the actual PID)
taskkill /F /PID 12345
```

**Linux/Mac:**
```bash
# Find the process
lsof -i :50051
# Kill it (replace 12345 with the actual PID)
kill -9 12345
```

### Connection Refused

If the GraphQL service can't connect to the Post service, ensure:

1. The Post service is running
2. It's available on localhost:50051
3. No firewall is blocking the connection

### Empty Timeline

If the timeline API returns no posts:

1. Check that the user ID exists (e.g., "user1" through "user5")
2. Verify the user has followers in the mock data
3. Ensure those users have posts in the mock data

## Accessing the API

After starting the services, you can:

1. Open the GraphQL Playground at http://localhost:8080
2. Send GraphQL queries to http://localhost:8080/query

## Example Query

Try this query in the GraphQL Playground:

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

## Stopping the Services

If using the orchestrator, press `Ctrl+C` in the terminal to stop all services.

If running services individually, press `Ctrl+C` in each terminal window. 