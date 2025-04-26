# Development Workflow

This document outlines the development workflow for the Paper.Social Feed Service.

## Development Environment Setup

### Required Tools

1. **Go**
   - Version: 1.22 or later
   - Installation: [golang.org/dl](https://golang.org/dl/)
   - Environment variables:
     ```bash
     # Add to your shell profile (.bashrc, .zshrc, etc.)
     export GOPATH=$HOME/go
     export PATH=$PATH:$GOPATH/bin
     ```

2. **Protocol Buffers**
   - Version: v6.30.2 or later
   - Installation:
     - Windows: Download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases)
     - macOS: `brew install protobuf`
     - Linux: `apt-get install protobuf-compiler`

3. **Go Tools**
   ```bash
   # Protocol Buffers
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

   # GraphQL
   go install github.com/99designs/gqlgen@latest
   ```

### IDE Setup

1. **VS Code**
   - Install Go extension
   - Install GraphQL extension
   - Install Proto3 extension
   - Recommended settings:
     ```json
     {
       "go.useLanguageServer": true,
       "go.lintTool": "golangci-lint",
       "editor.formatOnSave": true
     }
     ```

2. **GoLand**
   - Enable Go modules integration
   - Install GraphQL plugin
   - Install Protocol Buffers plugin

## Development Workflow

### 1. Making Changes to GraphQL Schema

1. **Edit Schema**
   - Modify `graphqlservice/graph/schema.graphqls`
   - Add new types, queries, or mutations

2. **Generate Code**
   ```bash
   go run github.com/99designs/gqlgen generate
   ```

3. **Implement Resolvers**
   - Edit `graphqlservice/graph/resolver.go`
   - Implement new resolver methods
   - Add tests for new functionality

4. **Test Changes**
   ```bash
   # Run GraphQL server
   go run ./cmd/server

   # Access GraphQL Playground
   open http://localhost:8080/graphql
   ```

### 2. Making Changes to gRPC Service

1. **Edit Proto File**
   - Modify `proto/post/post.proto`
   - Add new messages or RPC methods

2. **Generate Code**
   ```bash
   cd proto/post
   protoc --go_out=. --go-grpc_out=. post.proto
   ```

3. **Implement Service**
   - Create/update service implementation
   - Add tests for new functionality

4. **Test Changes**
   ```bash
   # Run gRPC server
   go run ./cmd/grpc

   # Use grpcurl for testing
   grpcurl -plaintext localhost:50051 list
   ```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./graphqlservice/...
go test ./proto/...

# Run with coverage
go test -cover ./...
```

### Integration Tests

```bash
# Start services
docker-compose up -d

# Run integration tests
go test -tags=integration ./...
```

### Load Testing

```bash
# GraphQL endpoint
hey -n 1000 -c 50 -m POST -H "Content-Type: application/json" \
  -d '{"query":"{ getTimeline(userId:\"user1\") { id content } }"}' \
  http://localhost:8080/query

# gRPC endpoint
ghz --insecure \
  --proto proto/post/post.proto \
  --call post.PostService.ListPostsByUser \
  -d '{"user_id":"user1"}' \
  localhost:50051
```

## Code Quality

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Code Generation Verification

```bash
# Verify GraphQL schema
gqlgen validate

# Verify proto syntax
protoc --experimental_allow_proto3_optional --validate_out="lang=go:." post.proto
```

### Pre-commit Checks

1. **Install pre-commit hooks**
   ```bash
   # Install pre-commit
   pip install pre-commit

   # Install hooks
   pre-commit install
   ```

2. **Run checks manually**
   ```bash
   pre-commit run --all-files
   ```

## Deployment

### Local Development

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up graphql-service

# View logs
docker-compose logs -f
```

### Production Deployment

1. **Build Images**
   ```bash
   # Build all services
   docker-compose build

   # Build specific service
   docker-compose build graphql-service
   ```

2. **Push Images**
   ```bash
   docker-compose push
   ```

3. **Deploy**
   ```bash
   # Using kubectl
   kubectl apply -f k8s/

   # Using helm
   helm upgrade --install feed-service ./helm
   ```

## Monitoring and Debugging

### Metrics

- GraphQL metrics available at `/metrics`
- gRPC metrics exposed via Prometheus
- Default port: 2112

### Tracing

- OpenTelemetry integration
- Jaeger UI available at `http://localhost:16686`

### Logging

- Structured logging using zerolog
- Log levels: debug, info, warn, error
- JSON format for production

## Common Tasks

### Adding a New Feature

1. Plan changes needed in both GraphQL and gRPC
2. Update proto definition if needed
3. Update GraphQL schema if needed
4. Generate code for both services
5. Implement business logic
6. Add tests
7. Update documentation

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/99designs/gqlgen

# Tidy go.mod
go mod tidy
```

### Regenerating All Code

```bash
# GraphQL
go run github.com/99designs/gqlgen generate

# gRPC
cd proto/post && \
protoc --go_out=. --go-grpc_out=. post.proto
``` 