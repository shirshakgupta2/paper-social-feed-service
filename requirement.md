Objective: Build a simplified backend microservice for a social feed system like Paper.Social (a Twitter/X alternative)

📌 Project Overview
This task is designed to evaluate your real-world skills in building a backend microservice using Go, GraphQL, and gRPC — the core technologies at Paper.Social.
You’ll build a microservice that delivers a user’s timeline by aggregating posts from users they follow.

Here Posts are simple text and links in the content 

🎯 Task Requirements
1. GraphQL API
Create a GraphQL API with the following query:
graphql
getTimeline(userId: ID!): [Post]

Should return the 20 most recent posts from users that the given userId follows.
Posts must be sorted in reverse chronological order.

2. Simulated gRPC Post Service
Create a simple gRPC service that simulates fetching posts for a user:
proto
rpc ListPostsByUser(ListPostsRequest) returns (ListPostsResponse)

Populate the service with mocked data (5 users, a few posts each).
The feed service should communicate with this gRPC service to get post data.

3. Data Model
Use in-memory data structures to simulate:
Users
Follower relationships
Posts (timestamp, content, author)
You do not need to persist data or connect to a real database.

4. Feed Aggregation Logic
Fetch posts of all followed users.
Combine and sort them by timestamp.
Return the latest 20 posts.
Use Go’s concurrency features (e.g., goroutines) to fetch in parallel where appropriate.

✅ Bonus Points
Efficient concurrency and batching logic.
Clean project structure and idiomatic Go.
Basic unit tests for timeline generation.
README explaining your decisions.
