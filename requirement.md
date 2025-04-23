Objective: Build a simplified backend microservice for a social feed system like Paper.Social (a Twitter/X alternative)

📌 Project Overview
This task is designed to evaluate your real-world skills in building a backend microservice using Go, GraphQL, and gRPC — the core technologies at Paper.Social.
You’ll build a microservice that delivers a user’s timeline by aggregating posts from users they follow.

Here Posts are simple text and links in the content  and we can keep  any random image links from google related to server architecture diagrams(https://www.google.com/search?sca_esv=cc91aa7b516a412e&rlz=1C1ONGR_enIN1142IN1142&sxsrf=AHTn8zqdCQg8-f05dRJbH49xCR9qN2xH0Q:1745351669221&q=server+architecture+diagram&udm=2&fbs=ABzOT_CWdhQLP1FcmU5B0fn3xuWpA-dk4wpBWOGsoR7DG5zJBkzPWUS0OtApxR2914vrjk4ZqZZ4I2IkJifuoUeV0iQtlsVaSqiwnznvC1owt2z2tTdc23Auc6X4y2i7IIF0f-d_O-E9yXafSm5foej9KNb5dB5UNNsgm78dv2qEeljVjLTUov5wWn4x9of_4BNb8vF_2a_9-AxwH0UJGyfTMDuJ_sz_gg&sa=X&ved=2ahUKEwiK9YHSteyMAxXfcPUHHagtAi8QtKgLegQIGRAB&biw=1920&bih=911&dpr=1)


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
