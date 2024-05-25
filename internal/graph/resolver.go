package graph

import (
	"sync"

	"github.com/lomins/ozon-comments-graphql/internal/graph/model"
	"github.com/lomins/ozon-comments-graphql/pkg/storage"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage  storage.Storage
	Comments map[string]chan *model.Comment
	mu       sync.Mutex
}
