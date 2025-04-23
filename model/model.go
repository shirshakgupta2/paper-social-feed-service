package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// User represents a user in the system
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Follows  []string `json:"follows"` // IDs of users this user follows
}

// Post represents a post in the system
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetImageURLsFromContent extracts image URLs from post content
func (p *Post) GetImageURLsFromContent() []string {
	// Regular expression to match image URLs
	imageRegex := regexp.MustCompile(`https?:\/\/\S+\.(jpg|jpeg|png|gif|webp)(\?\S+)?`)
	return imageRegex.FindAllString(p.Content, -1)
}

// ContainsImages checks if the post content contains images
func (p *Post) ContainsImages() bool {
	return len(p.GetImageURLsFromContent()) > 0
}

// HasServerArchitectureDiagram checks if the post might contain server architecture diagrams
// This is a simple check based on content keywords
func (p *Post) HasServerArchitectureDiagram() bool {
	lowerContent := strings.ToLower(p.Content)
	keywords := []string{
		"architecture",
		"diagram",
		"server",
		"microservice",
		"infrastructure",
		"system design",
	}

	for _, keyword := range keywords {
		if strings.Contains(lowerContent, keyword) {
			return true
		}
	}

	return false
}

// Database is an in-memory database simulation
type Database struct {
	Users      map[string]*User
	Posts      map[string][]*Post // Posts indexed by user ID
	PostsByID  map[string]*Post   // Posts indexed by ID for faster lookups
	NextPostID int                // Used to generate unique post IDs
}

// NewDatabase creates a new in-memory database with mock data
func NewDatabase() *Database {
	db := &Database{
		Users:      make(map[string]*User),
		Posts:      make(map[string][]*Post),
		PostsByID:  make(map[string]*Post),
		NextPostID: 11, // Start after our initial posts
	}

	// Create mock users
	users := []User{
		{ID: "user1", Username: "alice", Follows: []string{"user2", "user3", "user4"}},
		{ID: "user2", Username: "bob", Follows: []string{"user1", "user5"}},
		{ID: "user3", Username: "charlie", Follows: []string{"user1", "user2"}},
		{ID: "user4", Username: "dave", Follows: []string{"user1", "user5"}},
		{ID: "user5", Username: "eve", Follows: []string{"user1", "user3"}},
	}

	for _, u := range users {
		userCopy := u
		db.Users[u.ID] = &userCopy
	}

	// Create mock posts
	now := time.Now()
	posts := []Post{
		{ID: "post1", UserID: "user1", Content: "Hello, world!", CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "post2", UserID: "user1", Content: "GraphQL is awesome", CreatedAt: now.Add(-2 * time.Hour)},
		{ID: "post3", UserID: "user2", Content: "gRPC is cool", CreatedAt: now.Add(-30 * time.Minute)},
		{ID: "post4", UserID: "user2", Content: "Check out this server architecture diagram: https://example.com/architecture-diagram.png", CreatedAt: now.Add(-3 * time.Hour)},
		{ID: "post5", UserID: "user3", Content: "Working on a new project", CreatedAt: now.Add(-4 * time.Hour)},
		{ID: "post6", UserID: "user3", Content: "Learning Go microservices with this system design diagram: https://example.com/microservices.jpg", CreatedAt: now.Add(-10 * time.Minute)},
		{ID: "post7", UserID: "user4", Content: "Just deployed my first service!", CreatedAt: now.Add(-1 * time.Minute)},
		{ID: "post8", UserID: "user4", Content: "Anyone else using Protocol Buffers?", CreatedAt: now.Add(-5 * time.Hour)},
		{ID: "post9", UserID: "user5", Content: "Server architecture examples: https://example.com/server-arch.png", CreatedAt: now},
		{ID: "post10", UserID: "user5", Content: "Microservices vs monoliths", CreatedAt: now.Add(-2 * time.Minute)},
	}

	for _, p := range posts {
		postCopy := p
		db.Posts[p.UserID] = append(db.Posts[p.UserID], &postCopy)
		db.PostsByID[p.ID] = &postCopy
	}

	return db
}

// GetUserByID retrieves a user by ID
func (db *Database) GetUserByID(id string) *User {
	return db.Users[id]
}

// GetPostsByUserID retrieves posts for a specific user
func (db *Database) GetPostsByUserID(userID string) []*Post {
	return db.Posts[userID]
}

// CreatePost creates a new post for a user and returns it
func (db *Database) CreatePost(userID string, content string) (*Post, error) {
	// Check if user exists
	if _, exists := db.Users[userID]; !exists {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}

	// Generate unique post ID
	postID := fmt.Sprintf("post%d", db.NextPostID)
	db.NextPostID++

	// Create the post
	post := &Post{
		ID:        postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// Add to user's posts
	db.Posts[userID] = append(db.Posts[userID], post)

	// Add to posts by ID map
	db.PostsByID[postID] = post

	return post, nil
}

// UpdatePost updates an existing post
func (db *Database) UpdatePost(postID string, content string) (*Post, error) {
	// Check if post exists
	post, exists := db.PostsByID[postID]
	if !exists {
		return nil, fmt.Errorf("post with ID %s not found", postID)
	}

	// Update the content
	post.Content = content

	return post, nil
}

// DeletePost deletes a post
func (db *Database) DeletePost(postID string) (bool, error) {
	// Check if post exists
	post, exists := db.PostsByID[postID]
	if !exists {
		return false, fmt.Errorf("post with ID %s not found", postID)
	}

	// Get the user's posts
	userPosts := db.Posts[post.UserID]

	// Find and remove the post from the user's posts
	for i, p := range userPosts {
		if p.ID == postID {
			// Remove this post from the slice
			db.Posts[post.UserID] = append(userPosts[:i], userPosts[i+1:]...)
			break
		}
	}

	// Remove from the post ID map
	delete(db.PostsByID, postID)

	return true, nil
}
