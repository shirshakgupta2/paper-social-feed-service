package graphqlservice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/paper-social/feed-service/model"
	"github.com/paper-social/feed-service/postservice"
	"github.com/paper-social/feed-service/proto/post"
)

// Service represents the GraphQL service
type Service struct {
	db         *model.Database
	postClient *postservice.Client
}

// Post represents a post in the GraphQL schema
type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// DeleteResponse represents the response to a delete operation
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewService creates a new GraphQL service
func NewService(db *model.Database, postServiceAddr string) *Service {
	// Create a client to connect to the post service
	client := postservice.CreateClient(postServiceAddr)

	return &Service{
		db:         db,
		postClient: client,
	}
}

// GetTimeline retrieves timeline posts for a user
func (s *Service) GetTimeline(ctx context.Context, userID string) ([]*Post, error) {
	// Get the user
	user := s.db.GetUserByID(userID)
	if user == nil {
		return nil, nil
	}

	// Prepare to fetch posts for all followed users
	var wg sync.WaitGroup
	var mu sync.Mutex
	allPosts := make([]*Post, 0)

	// For each followed user, fetch their posts
	for _, followedUserID := range user.Follows {
		wg.Add(1)
		go func(followedID string) {
			defer wg.Done()

			// Call the post service to get posts for this user
			resp, err := s.postClient.ListPostsByUser(ctx, &post.ListPostsRequest{UserId: followedID})
			if err != nil {
				log.Printf("Error fetching posts for user %s: %v", followedID, err)
				return
			}

			// Convert proto posts to our Post type
			posts := make([]*Post, 0, len(resp.Posts))
			for _, p := range resp.Posts {
				// Convert Unix timestamp to RFC3339 format
				t := time.Unix(p.CreatedAt, 0)
				posts = append(posts, &Post{
					ID:        p.Id,
					UserID:    p.UserId,
					Content:   p.Content,
					CreatedAt: t.Format(time.RFC3339),
				})
			}

			// Lock and update the allPosts slice
			mu.Lock()
			allPosts = append(allPosts, posts...)
			mu.Unlock()
		}(followedUserID)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Sort posts by creation time (newest first)
	sort.Slice(allPosts, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, allPosts[i].CreatedAt)
		timeJ, _ := time.Parse(time.RFC3339, allPosts[j].CreatedAt)
		return timeI.After(timeJ)
	})

	// Return at most 20 posts
	if len(allPosts) > 20 {
		return allPosts[:20], nil
	}
	return allPosts, nil
}

// CreatePost creates a new post
func (s *Service) CreatePost(ctx context.Context, userID string, content string) (*Post, error) {
	// Call the post service to create a post
	resp, err := s.postClient.CreatePost(ctx, &post.CreatePostRequest{
		UserId:  userID,
		Content: content,
	})
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return nil, err
	}

	// Convert proto post to our Post type
	t := time.Unix(resp.CreatedAt, 0)
	return &Post{
		ID:        resp.Id,
		UserID:    resp.UserId,
		Content:   resp.Content,
		CreatedAt: t.Format(time.RFC3339),
	}, nil
}

// UpdatePost updates an existing post
func (s *Service) UpdatePost(ctx context.Context, id string, content string) (*Post, error) {
	// Call the post service to update a post
	resp, err := s.postClient.UpdatePost(ctx, &post.UpdatePostRequest{
		Id:      id,
		Content: content,
	})
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return nil, err
	}

	// Convert proto post to our Post type
	t := time.Unix(resp.CreatedAt, 0)
	return &Post{
		ID:        resp.Id,
		UserID:    resp.UserId,
		Content:   resp.Content,
		CreatedAt: t.Format(time.RFC3339),
	}, nil
}

// DeletePost deletes a post
func (s *Service) DeletePost(ctx context.Context, id string) (*DeleteResponse, error) {
	// Call the post service to delete a post
	resp, err := s.postClient.DeletePost(ctx, &post.DeletePostRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		return &DeleteResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &DeleteResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// TimelineQuery represents a GraphQL query for a user's timeline
type TimelineQuery struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName,omitempty"`
}

// GraphQLResponse represents a response from the GraphQL API
type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

// StartService starts the GraphQL service with a simplified HTTP handler
func StartService(db *model.Database, postServiceAddr, httpAddr string) error {
	service := NewService(db, postServiceAddr)

	// Create a handler for GraphQL queries
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only handle POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Read the request body for debugging
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			resp := GraphQLResponse{
				Errors: []string{"Error reading request: " + err.Error()},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Log the raw request body for debugging
		log.Printf("Raw request body: %s", string(body))

		// Create a new reader from the body for json.Decoder
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Decode the request
		var query TimelineQuery
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			log.Printf("Error decoding request: %v", err)
			resp := GraphQLResponse{
				Errors: []string{"Invalid request: " + err.Error()},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		log.Printf("Received GraphQL query: %s", query.Query)
		log.Printf("Variables: %+v", query.Variables)

		// Clean up the query - remove escaped newlines and unnecessary whitespace
		cleanQuery := strings.ReplaceAll(query.Query, "\\n", " ")
		cleanQuery = strings.TrimSpace(cleanQuery)
		log.Printf("Cleaned query: %s", cleanQuery)

		// Check if this is a getTimeline query
		if strings.Contains(cleanQuery, "getTimeline") {
			// Get the timeline
			userId := extractUserIdFromQuery(cleanQuery)
			log.Printf("Extracted userId for timeline: %s", userId)

			// If userId wasn't found in query, try variables
			if userId == "" && query.Variables != nil {
				userId = getUserIDFromVariables(query.Variables)
				log.Printf("Got userId from variables: %s", userId)
			}

			timeline, err := service.GetTimeline(r.Context(), userId)
			if err != nil {
				resp := GraphQLResponse{
					Errors: []string{err.Error()},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// Return the timeline
			resp := GraphQLResponse{
				Data: map[string]interface{}{
					"getTimeline": timeline,
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Check for mutations
		if strings.Contains(cleanQuery, "mutation") {
			// Handle create post
			if strings.Contains(cleanQuery, "createPost") {
				// Extract userId and content using regex for more reliable parsing
				userIdMatch := regexp.MustCompile(`userId\s*:\s*"([^"]*)"`)
				contentMatch := regexp.MustCompile(`content\s*:\s*"([^"]*)"`)

				userIdMatches := userIdMatch.FindStringSubmatch(cleanQuery)
				contentMatches := contentMatch.FindStringSubmatch(cleanQuery)

				var userId, content string
				if len(userIdMatches) >= 2 {
					userId = userIdMatches[1]
				} else if query.Variables != nil {
					userId = getUserIDFromVariables(query.Variables)
				}

				if len(contentMatches) >= 2 {
					content = contentMatches[1]
				}

				log.Printf("Creating post - userId: %s, content: %s", userId, content)

				if userId == "" {
					resp := GraphQLResponse{
						Errors: []string{"Missing required parameter: userId"},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				post, err := service.CreatePost(r.Context(), userId, content)
				if err != nil {
					resp := GraphQLResponse{
						Errors: []string{err.Error()},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				resp := GraphQLResponse{
					Data: map[string]interface{}{
						"createPost": post,
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// Handle update post (similar approach with regex)
			if strings.Contains(cleanQuery, "updatePost") {
				idMatch := regexp.MustCompile(`id\s*:\s*"([^"]*)"`)
				contentMatch := regexp.MustCompile(`content\s*:\s*"([^"]*)"`)

				idMatches := idMatch.FindStringSubmatch(cleanQuery)
				contentMatches := contentMatch.FindStringSubmatch(cleanQuery)

				var id, content string
				if len(idMatches) >= 2 {
					id = idMatches[1]
				} else if query.Variables != nil && query.Variables["id"] != nil {
					id = query.Variables["id"].(string)
				}

				if len(contentMatches) >= 2 {
					content = contentMatches[1]
				}

				log.Printf("Updating post - id: %s, content: %s", id, content)

				if id == "" {
					resp := GraphQLResponse{
						Errors: []string{"Missing required parameter: id"},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				post, err := service.UpdatePost(r.Context(), id, content)
				if err != nil {
					resp := GraphQLResponse{
						Errors: []string{err.Error()},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				resp := GraphQLResponse{
					Data: map[string]interface{}{
						"updatePost": post,
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// Handle delete post
			if strings.Contains(cleanQuery, "deletePost") {
				idMatch := regexp.MustCompile(`id\s*:\s*"([^"]*)"`)
				idMatches := idMatch.FindStringSubmatch(cleanQuery)

				var id string
				if len(idMatches) >= 2 {
					id = idMatches[1]
				} else if query.Variables != nil && query.Variables["id"] != nil {
					id = query.Variables["id"].(string)
				}

				log.Printf("Deleting post - id: %s", id)

				if id == "" {
					resp := GraphQLResponse{
						Errors: []string{"Missing required parameter: id"},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				deleteResp, err := service.DeletePost(r.Context(), id)
				if err != nil {
					resp := GraphQLResponse{
						Errors: []string{err.Error()},
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

				resp := GraphQLResponse{
					Data: map[string]interface{}{
						"deletePost": deleteResp,
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		// Return a proper JSON error response
		resp := GraphQLResponse{
			Errors: []string{"Unsupported query type"},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	})

	// Create a simple GraphQL playground
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		playground := `
<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <meta charset="utf-8">
    <style>
        body { margin: 0; padding: 20px; font-family: Arial, sans-serif; }
        h1 { color: #333; }
        #query { width: 100%; height: 200px; padding: 10px; margin-bottom: 10px; font-family: monospace; }
        #variables { width: 100%; height: 100px; padding: 10px; margin-bottom: 10px; font-family: monospace; }
        #response { width: 100%; height: 300px; padding: 10px; background-color: #f5f5f5; font-family: monospace; overflow: auto; }
        button { padding: 10px 20px; background-color: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background-color: #45a049; }
        .container { max-width: 800px; margin: 0 auto; }
        .tabs { display: flex; margin-bottom: 10px; }
        .tab { padding: 10px 20px; cursor: pointer; border: 1px solid #ccc; background: #f5f5f5; }
        .tab.active { background: #4CAF50; color: white; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GraphQL Playground</h1>
        
        <div class="tabs">
            <div class="tab active" onclick="setQuery('timeline')">Timeline</div>
            <div class="tab" onclick="setQuery('create')">Create Post</div>
            <div class="tab" onclick="setQuery('update')">Update Post</div>
            <div class="tab" onclick="setQuery('delete')">Delete Post</div>
        </div>
        
        <h3>Query:</h3>
        <textarea id="query">{
  getTimeline(userId: "user1") {
    id
    userId
    content
    createdAt
  }
}</textarea>
        <h3>Variables:</h3>
        <textarea id="variables">{
  "userId": "user1"
}</textarea>
        <button onclick="executeQuery()">Execute Query</button>
        <h3>Response:</h3>
        <pre id="response"></pre>
    </div>

    <script>
        function executeQuery() {
            const query = document.getElementById('query').value;
            const variables = JSON.parse(document.getElementById('variables').value || '{}');
            
            fetch('/query', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query, variables })
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById('response').textContent = JSON.stringify(data, null, 2);
            })
            .catch(error => {
                document.getElementById('response').textContent = 'Error: ' + error;
            });
        }
        
        function setQuery(type) {
            // Set all tabs inactive
            document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
            
            // Set the clicked tab active
            if (event) {
                event.target.classList.add('active');
            } else {
                document.querySelector('.tab[onclick="setQuery(\'' + type + '\')"]').classList.add('active');
            }
            
            let query = '';
            let variables = {};
            
            switch(type) {
                case 'timeline':
                    query = '{\\n  getTimeline(userId: \\"user1\\") {\\n    id\\n    userId\\n    content\\n    createdAt\\n  }\\n}';
                    variables = { userId: "user1" };
                    break;
                case 'create':
                    query = 'mutation {\\n  createPost(userId: \\"user1\\", content: \\"This is a new post!\\") {\\n    id\\n    userId\\n    content\\n    createdAt\\n  }\\n}';
                    variables = {};
                    break;
                case 'update':
                    query = 'mutation {\\n  updatePost(id: \\"post1\\", content: \\"Updated content for this post\\") {\\n    id\\n    userId\\n    content\\n    createdAt\\n  }\\n}';
                    variables = {};
                    break;
                case 'delete':
                    query = 'mutation {\\n  deletePost(id: \\"post1\\") {\\n    success\\n    message\\n  }\\n}';
                    variables = {};
                    break;
            }
            
            document.getElementById('query').value = query;
            document.getElementById('variables').value = JSON.stringify(variables, null, 2);
        }
    </script>
</body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(playground))
	})

	log.Printf("GraphQL service starting on %s", httpAddr)
	log.Printf("Connect to the GraphQL playground: http://localhost%s", httpAddr)

	return http.ListenAndServe(httpAddr, nil)
}

// Helper function to get userId from variables
func getUserIDFromVariables(variables map[string]interface{}) string {
	if variables == nil {
		return ""
	}

	// Check for userId directly
	if userId, ok := variables["userId"].(string); ok {
		return userId
	}

	// Check for ID field that might be used instead
	if id, ok := variables["id"].(string); ok {
		return id
	}

	return ""
}

// Helper function to extract userId from a query
func extractUserIdFromQuery(query string) string {
	// Find the userId in the query
	userIdMatch := regexp.MustCompile(`userId\s*:\s*"([^"]*)"`)
	userIdMatches := userIdMatch.FindStringSubmatch(query)

	var userId string
	if len(userIdMatches) >= 2 {
		userId = userIdMatches[1]
	}
	return userId
}
