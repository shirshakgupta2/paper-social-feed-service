# Paper.Social Feed Service

This service handles the social feed functionality for Paper.Social, implementing both GraphQL and gRPC APIs for post management.

## Project Structure

```
.
├── graphqlservice/           # GraphQL API service
│   └── graph/
│       ├── model/           # GraphQL models
│       ├── schema.graphqls  # GraphQL schema
│       └── resolver.go      # GraphQL resolvers
├── proto/                   # gRPC service
│   └── post/
│       ├── post.proto      # Protocol Buffers definition
│       ├── post.pb.go      # Generated message types
│       └── post_grpc.pb.go # Generated gRPC service
└── README.md               # This file
```

## Prerequisites

- Go 1.22 or later
- Protocol Buffers compiler (protoc) v6.30.2 or later
- GraphQL code generator (gqlgen)
- gRPC tools

## Installation

1. **Install Go**
   - Download from [golang.org](https://golang.org/dl/)
   - Verify installation:
     ```bash
     go version
     ```

2. **Install Protocol Buffers Compiler**
   - Download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases)
   - Add to PATH
   - Verify installation:
     ```bash
     protoc --version
     ```

3. **Install Go tools**
   ```bash
   # Install protoc-gen-go (Protocol Buffers for Go)
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

   # Install protoc-gen-go-grpc (gRPC for Go)
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

   # Install gqlgen (GraphQL for Go)
   go install github.com/99designs/gqlgen@latest
   ```

4. **Clone the repository**
   ```bash
   git clone https://github.com/paper-social/feed-service.git
   cd feed-service
   ```

5. **Install dependencies**
   ```bash
   go mod download
   ```

## Code Generation

### GraphQL Code Generation

The GraphQL code is generated using `gqlgen` based on the schema defined in `graphqlservice/graph/schema.graphqls`.

1. **Configuration**
   - The generation is configured in `gqlgen.yml`
   - Key configuration includes:
     - Schema location
     - Output locations for generated code
     - Type mappings

2. **Generate code**
   ```bash
   go run github.com/99designs/gqlgen generate
   ```

3. **Generated files**
   - `graphqlservice/graph/model/models_gen.go`: Generated models
   - `graphqlservice/graph/generated/generated.go`: Generated GraphQL server code

### gRPC Code Generation

The gRPC code is generated using `protoc` based on the protocol buffer definition in `proto/post/post.proto`.

1. **Generate code**
   ```bash
   cd proto/post
   protoc --go_out=. --go-grpc_out=. post.proto
   ```

2. **Generated files**
   - `post.pb.go`: Message types and serialization code
   - `post_grpc.pb.go`: gRPC service definitions

## Development

### GraphQL Development

1. **Schema Updates**
   - Modify `graphqlservice/graph/schema.graphqls`
   - Run code generation
   - Implement new resolvers in `resolver.go`

2. **Run GraphQL server**
   ```bash
   go run ./cmd/server
   ```

3. **Access GraphQL Playground**
   - Open `http://localhost:8080/graphql` in your browser

### gRPC Development

1. **Proto Updates**
   - Modify `proto/post/post.proto`
   - Run code generation
   - Implement new service methods

2. **Run gRPC server**
   ```bash
   go run ./cmd/grpc
   ```

## API Documentation

### GraphQL API

The GraphQL API provides the following operations:

```graphql
type Query {
  getTimeline(userId: ID!): [Post!]!
}

type Mutation {
  createPost(userId: ID!, content: String!): Post!
  updatePost(id: ID!, content: String!): Post!
  deletePost(id: ID!): DeleteResponse!
}
```

### gRPC API

The gRPC service provides the following methods:

```protobuf
service PostService {
  rpc ListPostsByUser(ListPostsRequest) returns (ListPostsResponse);
  rpc CreatePost(CreatePostRequest) returns (Post);
  rpc UpdatePost(UpdatePostRequest) returns (Post);
  rpc DeletePost(DeletePostRequest) returns (DeletePostResponse);
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 