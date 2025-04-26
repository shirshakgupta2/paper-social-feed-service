package model

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

// Post represents a post in the GraphQL model
type Post struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
	ImageUrls []string `json:"imageUrls,omitempty"`
}

// DeleteResponse represents the response to a delete post operation
type DeleteResponse struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
}
