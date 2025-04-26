# Code Generation Guide

This document provides detailed information about the code generation processes used in the Paper.Social Feed Service.

## Table of Contents

1. [GraphQL Code Generation](#graphql-code-generation)
2. [gRPC Code Generation](#grpc-code-generation)
3. [Troubleshooting](#troubleshooting)

## GraphQL Code Generation

### Overview

The GraphQL code generation is handled by `gqlgen`, which generates Go code from your GraphQL schema. This ensures type-safe operations and maintains consistency between your schema and implementation.

### Configuration

The configuration is defined in `gqlgen.yml`:

```yaml
# Schema location
schema:
  - graphqlservice/graph/*.graphqls

# Generated server code
exec:
  filename: graphqlservice/graph/generated/generated.go
  package: generated

# Generated models
model:
  filename: graphqlservice/graph/model/models_gen.go
  package: model

# Resolver implementation
resolver:
  layout: follow-schema
  dir: graphqlservice/graph
  package: graph
  filename_template: "{name}.resolvers.go"
```

### Schema Definition

The GraphQL schema is defined in `graphqlservice/graph/schema.graphqls`:

```graphql
type Post {
  id: ID!
  userId: ID!
  content: String!
  createdAt: String!
  imageUrls: [String!]
}

type Query {
  getTimeline(userId: ID!): [Post!]!
}

type Mutation {
  createPost(userId: ID!, content: String!): Post!
  updatePost(id: ID!, content: String!): Post!
  deletePost(id: ID!): DeleteResponse!
}
```

### Generated Files

1. **Models (`models_gen.go`)**
   - Location: `graphqlservice/graph/model/models_gen.go`
   - Contains Go structs that match GraphQL types
   - Example:
     ```go
     type Post struct {
         ID        string   `json:"id"`
         UserID    string   `json:"userId"`
         Content   string   `json:"content"`
         CreatedAt string   `json:"createdAt"`
         ImageUrls []string `json:"imageUrls,omitempty"`
     }
     ```

2. **Generated Server Code (`generated.go`)**
   - Location: `graphqlservice/graph/generated/generated.go`
   - Contains:
     - Schema validation
     - Query/mutation execution logic
     - Type system definitions

### Generation Process

1. **Install gqlgen**
   ```bash
   go install github.com/99designs/gqlgen@latest
   ```

2. **Generate code**
   ```bash
   go run github.com/99designs/gqlgen generate
   ```

3. **Implement resolvers**
   - Edit `graphqlservice/graph/resolver.go`
   - Implement the generated interfaces

## gRPC Code Generation

### Overview

The gRPC code generation uses `protoc` (Protocol Buffers Compiler) to generate Go code from your protocol buffer definitions.

### Prerequisites

1. **Install Protocol Buffers Compiler**
   - Download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases)
   - Add to PATH
   - Verify: `protoc --version`

2. **Install Go plugins**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

### Proto Definition

The service definition is in `proto/post/post.proto`:

```protobuf
syntax = "proto3";

package post;

option go_package = "github.com/paper-social/feed-service/proto/post";

service PostService {
  rpc ListPostsByUser(ListPostsRequest) returns (ListPostsResponse);
  rpc CreatePost(CreatePostRequest) returns (Post);
  rpc UpdatePost(UpdatePostRequest) returns (Post);
  rpc DeletePost(DeletePostRequest) returns (DeletePostResponse);
}

message Post {
  string id = 1;
  string user_id = 2;
  string content = 3;
  int64 created_at = 4;
}
```

### Generated Files

1. **Message Types (`post.pb.go`)**
   - Contains:
     - Go structs for messages
     - Serialization code
     - Helper methods

2. **gRPC Service (`post_grpc.pb.go`)**
   - Contains:
     - Service interface definitions
     - Client implementation
     - Server interface
     - Registration functions

### Generation Process

1. **Navigate to proto directory**
   ```bash
   cd proto/post
   ```

2. **Generate code**
   ```bash
   protoc --go_out=. --go-grpc_out=. post.proto
   ```

## Troubleshooting

### Common Issues

1. **protoc-gen-go: program not found**
   - Solution: Ensure `$GOPATH/bin` is in your PATH
   - Check installation: `which protoc-gen-go`

2. **gqlgen generation fails**
   - Check `gqlgen.yml` configuration
   - Ensure schema is valid
   - Run with verbose output: `go run github.com/99designs/gqlgen generate -v`

3. **Import path issues**
   - Verify `go.mod` module name
   - Check `go_package` option in `.proto` files
   - Ensure correct package names in generated code

### Best Practices

1. **Version Control**
   - Commit generated code
   - Include generation commands in build scripts
   - Document any manual modifications

2. **Code Organization**
   - Keep proto files in dedicated directory
   - Separate generated and hand-written code
   - Use consistent package naming

3. **Maintenance**
   - Regenerate code after schema changes
   - Update resolvers when adding fields
   - Keep tool versions in sync across team

### Verification Steps

1. **GraphQL**
   ```bash
   # Verify schema
   gqlgen validate
   
   # Check generated types
   go run ./... -v
   ```

2. **gRPC**
   ```bash
   # Verify proto syntax
   protoc --experimental_allow_proto3_optional --validate_out="lang=go:." post.proto
   
   # Check generated code
   go build ./...
   ``` 