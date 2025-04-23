package graph

import (
	"context"
	"regexp"
	"strings"

	"github.com/paper-social/feed-service/graphqlservice"
	"github.com/paper-social/feed-service/graphqlservice/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the resolver root
type Resolver struct {
	Service *graphqlservice.Service
}

// Query returns the query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation returns the mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// QueryResolver is the interface for Query resolvers
type QueryResolver interface {
	GetTimeline(ctx context.Context, userID string) ([]*model.Post, error)
}

// MutationResolver is the interface for Mutation resolvers
type MutationResolver interface {
	CreatePost(ctx context.Context, userID string, content string) (*model.Post, error)
	UpdatePost(ctx context.Context, id string, content string) (*model.Post, error)
	DeletePost(ctx context.Context, id string) (*model.DeleteResponse, error)
}

// queryResolver is the implementation of the Query interface
type queryResolver struct {
	*Resolver
}

// mutationResolver is the implementation of the Mutation interface
type mutationResolver struct {
	*Resolver
}

// GetTimeline implements the GraphQL query to get a user's timeline
func (r *queryResolver) GetTimeline(ctx context.Context, userID string) ([]*model.Post, error) {
	posts, err := r.Service.GetTimeline(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert graphqlservice.Post to model.Post with additional fields
	result := make([]*model.Post, len(posts))
	for i, p := range posts {
		result[i] = &model.Post{
			ID:                           p.ID,
			UserID:                       p.UserID,
			Content:                      p.Content,
			CreatedAt:                    p.CreatedAt,
			ImageUrls:                    getImageURLs(p.Content),
			HasServerArchitectureDiagram: hasServerArchitectureDiagram(p.Content),
		}
	}

	return result, nil
}

// CreatePost implements the GraphQL mutation to create a new post
func (r *mutationResolver) CreatePost(ctx context.Context, userID string, content string) (*model.Post, error) {
	post, err := r.Service.CreatePost(ctx, userID, content)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:                           post.ID,
		UserID:                       post.UserID,
		Content:                      post.Content,
		CreatedAt:                    post.CreatedAt,
		ImageUrls:                    getImageURLs(post.Content),
		HasServerArchitectureDiagram: hasServerArchitectureDiagram(post.Content),
	}, nil
}

// UpdatePost implements the GraphQL mutation to update an existing post
func (r *mutationResolver) UpdatePost(ctx context.Context, id string, content string) (*model.Post, error) {
	post, err := r.Service.UpdatePost(ctx, id, content)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:                           post.ID,
		UserID:                       post.UserID,
		Content:                      post.Content,
		CreatedAt:                    post.CreatedAt,
		ImageUrls:                    getImageURLs(post.Content),
		HasServerArchitectureDiagram: hasServerArchitectureDiagram(post.Content),
	}, nil
}

// DeletePost implements the GraphQL mutation to delete a post
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*model.DeleteResponse, error) {
	response, err := r.Service.DeletePost(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.DeleteResponse{
		Success: response.Success,
		Message: &response.Message,
	}, nil
}

// Helper functions to implement the additional fields
func getImageURLs(content string) []string {
	// Use regex to extract image URLs from content
	// This is a simplified version that looks for common image extensions
	re := regexp.MustCompile(`https?:\/\/\S+\.(jpg|jpeg|png|gif|webp)(\?\S+)?`)
	return re.FindAllString(content, -1)
}

func hasServerArchitectureDiagram(content string) bool {
	// Simple implementation to check for server architecture diagram keywords
	lowerContent := strings.ToLower(content)
	keywords := []string{"architecture", "diagram", "server", "microservice", "infrastructure", "system design"}

	for _, keyword := range keywords {
		if strings.Contains(lowerContent, keyword) {
			return true
		}
	}

	return false
}
