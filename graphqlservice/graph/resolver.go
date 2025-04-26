package graph

import (
	"regexp"

	"github.com/paper-social/feed-service/graphqlservice"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Service *graphqlservice.Service
}

// Helper function to extract image URLs from content
func getImageURLs(content string) []string {
	// Use regex to extract image URLs from content
	// This is a simplified version that looks for common image extensions
	re := regexp.MustCompile(`https?:\/\/\S+\.(jpg|jpeg|png|gif|webp)(\?\S+)?`)
	return re.FindAllString(content, -1)
}
