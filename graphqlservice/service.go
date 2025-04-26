package graphqlservice

import (
	"context"
	"log"
	"sort"
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
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
	ImageURLs []string `json:"imageUrls,omitempty"`
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

				// Extract image URLs from content
				modelPost := &model.Post{
					Content: p.Content,
				}

				posts = append(posts, &Post{
					ID:        p.Id,
					UserID:    p.UserId,
					Content:   p.Content,
					CreatedAt: t.Format(time.RFC3339),
					ImageURLs: modelPost.GetImageURLsFromContent(),
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

	// Extract image URLs from content
	modelPost := &model.Post{
		Content: resp.Content,
	}

	return &Post{
		ID:        resp.Id,
		UserID:    resp.UserId,
		Content:   resp.Content,
		CreatedAt: t.Format(time.RFC3339),
		ImageURLs: modelPost.GetImageURLsFromContent(),
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

	// Extract image URLs from content
	modelPost := &model.Post{
		Content: resp.Content,
	}

	return &Post{
		ID:        resp.Id,
		UserID:    resp.UserId,
		Content:   resp.Content,
		CreatedAt: t.Format(time.RFC3339),
		ImageURLs: modelPost.GetImageURLsFromContent(),
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
