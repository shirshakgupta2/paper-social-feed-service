package postservice

import (
	"context"
	"log"
	"net"

	"github.com/paper-social/feed-service/model"
	"github.com/paper-social/feed-service/proto/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server implements the post service gRPC server
type Server struct {
	post.UnimplementedPostServiceServer
	db *model.Database
}

// NewServer creates a new post service server
func NewServer(db *model.Database) *Server {
	return &Server{db: db}
}

// ListPostsByUser implements the gRPC method to list posts by user
func (s *Server) ListPostsByUser(ctx context.Context, req *post.ListPostsRequest) (*post.ListPostsResponse, error) {
	log.Printf("Received request for posts of user: %s", req.UserId)

	// Get posts for the requested user
	posts := s.db.GetPostsByUserID(req.UserId)

	// Convert to proto posts
	pbPosts := make([]*post.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, &post.Post{
			Id:        p.ID,
			UserId:    p.UserID,
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Unix(),
		})
	}

	return &post.ListPostsResponse{Posts: pbPosts}, nil
}

// CreatePost implements the gRPC method to create a new post
func (s *Server) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.Post, error) {
	log.Printf("Creating post for user: %s", req.UserId)

	// Create the post in the database
	newPost, err := s.db.CreatePost(req.UserId, req.Content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return nil, err
	}

	// Convert to proto post
	pbPost := &post.Post{
		Id:        newPost.ID,
		UserId:    newPost.UserID,
		Content:   newPost.Content,
		CreatedAt: newPost.CreatedAt.Unix(),
	}

	return pbPost, nil
}

// UpdatePost implements the gRPC method to update an existing post
func (s *Server) UpdatePost(ctx context.Context, req *post.UpdatePostRequest) (*post.Post, error) {
	log.Printf("Updating post: %s", req.Id)

	// Update the post in the database
	updatedPost, err := s.db.UpdatePost(req.Id, req.Content)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return nil, err
	}

	// Convert to proto post
	pbPost := &post.Post{
		Id:        updatedPost.ID,
		UserId:    updatedPost.UserID,
		Content:   updatedPost.Content,
		CreatedAt: updatedPost.CreatedAt.Unix(),
	}

	return pbPost, nil
}

// DeletePost implements the gRPC method to delete a post
func (s *Server) DeletePost(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	log.Printf("Deleting post: %s", req.Id)

	// Delete the post from the database
	success, err := s.db.DeletePost(req.Id)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		return &post.DeletePostResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &post.DeletePostResponse{
		Success: success,
		Message: "Post deleted successfully",
	}, nil
}

// StartServer starts the gRPC server
func StartServer(db *model.Database, port string) error {
	// Create a TCP listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register our service
	server := NewServer(db)
	post.RegisterPostServiceServer(grpcServer, server)

	log.Printf("Post service gRPC server starting on %s (internal only)", port)

	// Start serving requests
	return grpcServer.Serve(listener)
}

// Client represents a client for the post service
type Client struct {
	client post.PostServiceClient
	conn   *grpc.ClientConn
}

// CreateClient creates a client to connect to the post service
func CreateClient(serverAddr string) *Client {
	// Set up a connection to the server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to post service: %v", err)
	}

	// Create a client
	client := post.NewPostServiceClient(conn)

	return &Client{
		client: client,
		conn:   conn,
	}
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListPostsByUser calls the post service to get posts for a user
func (c *Client) ListPostsByUser(ctx context.Context, req *post.ListPostsRequest) (*post.ListPostsResponse, error) {
	return c.client.ListPostsByUser(ctx, req)
}

// CreatePost calls the post service to create a new post
func (c *Client) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.Post, error) {
	return c.client.CreatePost(ctx, req)
}

// UpdatePost calls the post service to update an existing post
func (c *Client) UpdatePost(ctx context.Context, req *post.UpdatePostRequest) (*post.Post, error) {
	return c.client.UpdatePost(ctx, req)
}

// DeletePost calls the post service to delete a post
func (c *Client) DeletePost(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	return c.client.DeletePost(ctx, req)
}
