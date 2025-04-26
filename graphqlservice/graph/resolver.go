package graph

import (
	"github.com/paper-social/feed-service/graphqlservice"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the resolver for GraphQL queries and mutations
type Resolver struct {
	Service *graphqlservice.Service
}
