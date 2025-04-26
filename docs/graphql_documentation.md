# Paper.Social GraphQL API Documentation

## Overview
This document provides detailed information about the available GraphQL queries and mutations in the Paper.Social API.

## Base URL
```
http://localhost:8080/query
```

## Headers
```
Content-Type: application/json
```

## Queries

### Get Timeline
Retrieves the timeline for a specific user, showing posts from users they follow.

```graphql
query GetTimeline($userId: ID!) {
  getTimeline(userId: $userId) {
    id
    userId
    content
    createdAt
    imageUrls
  }
}
```

#### Variables
```json
{
  "userId": "user1"
}
```

#### Example Response
```json
{
  "data": {
    "getTimeline": [
      {
        "id": "post1",
        "userId": "user2",
        "content": "Example post content",
        "createdAt": "2024-03-20T10:00:00Z",
        "imageUrls": ["https://example.com/image1.jpg"]
      }
    ]
  }
}
```

## Mutations

### Create Post
Creates a new post for a specific user.

```graphql
mutation CreatePost($userId: ID!, $content: String!) {
  createPost(userId: $userId, content: $content) {
    id
    userId
    content
    createdAt
  }
}
```

#### Variables
```json
{
  "userId": "user1",
  "content": "This is a new post!"
}
```

#### Example Response
```json
{
  "data": {
    "createPost": {
      "id": "post1",
      "userId": "user1",
      "content": "This is a new post!",
      "createdAt": "2024-03-20T10:00:00Z"
    }
  }
}
```

### Update Post
Updates the content of an existing post.

```graphql
mutation UpdatePost($id: ID!, $content: String!) {
  updatePost(id: $id, content: $content) {
    id
    userId
    content
    createdAt
  }
}
```

#### Variables
```json
{
  "id": "post1",
  "content": "Updated content for this post"
}
```

#### Example Response
```json
{
  "data": {
    "updatePost": {
      "id": "post1",
      "userId": "user1",
      "content": "Updated content for this post",
      "createdAt": "2024-03-20T10:00:00Z"
    }
  }
}
```

### Delete Post
Deletes a specific post.

```graphql
mutation DeletePost($id: ID!) {
  deletePost(id: $id) {
    success
    message
  }
}
```

#### Variables
```json
{
  "id": "post1"
}
```

#### Example Response
```json
{
  "data": {
    "deletePost": {
      "success": true,
      "message": "Post deleted successfully"
    }
  }
}
```

## Types

### Post
```graphql
type Post {
  id: ID!
  userId: ID!
  content: String!
  createdAt: String!
  imageUrls: [String!]
}
```

### DeletePostResponse
```graphql
type DeletePostResponse {
  success: Boolean!
  message: String!
}
```

## Error Handling
The API may return errors in the following format:

```json
{
  "errors": [
    {
      "message": "Error message description",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": ["queryName"],
      "extensions": {
        "code": "ERROR_CODE"
      }
    }
  ]
}
```

## Best Practices
1. Always include the `Content-Type: application/json` header
2. Use variables for dynamic values instead of hardcoding them in the query
3. Request only the fields you need to minimize response size
4. Handle errors appropriately in your client application 