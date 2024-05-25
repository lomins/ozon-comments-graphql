package graph_test

import (
	"context"
	"testing"

	"github.com/lomins/ozon-comments-graphql/internal/graph"
	"github.com/lomins/ozon-comments-graphql/pkg/models"
	"github.com/lomins/ozon-comments-graphql/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post, err := resolver.Mutation().CreatePost(ctx, "Title", "Content", false)
	assert.NoError(t, err)
	assert.Equal(t, "Title", post.Title)
	assert.Equal(t, "Content", post.Content)
}

func TestCreateCommentResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post := &models.Post{ID: "1", Title: "Title", Content: "Content"}
	resolver.Storage.CreatePost(post)

	comment, err := resolver.Mutation().CreateComment(ctx, "1", nil, "Content")
	assert.NoError(t, err)
	assert.Equal(t, "1", comment.PostID)
	assert.Equal(t, "Content", comment.Content)
}

func TestDisableCommentsResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post := &models.Post{ID: "1", Title: "Title", Content: "Content"}
	resolver.Storage.CreatePost(post)

	updatedPost, err := resolver.Mutation().DisableComments(ctx, "1")
	assert.NoError(t, err)
	assert.True(t, updatedPost.CommentsDisabled)
}

func TestPostsResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post := &models.Post{ID: "1", Title: "Title", Content: "Content"}
	resolver.Storage.CreatePost(post)

	posts, err := resolver.Query().Posts(ctx)
	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Title", posts[0].Title)
}

func TestPostResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post := &models.Post{ID: "1", Title: "Title", Content: "Content"}
	resolver.Storage.CreatePost(post)

	fetchedPost, err := resolver.Query().Post(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, "Title", fetchedPost.Title)
}

func TestCommentsResolver(t *testing.T) {
	resolver := &graph.Resolver{Storage: storage.NewInMemoryStorage()}
	ctx := context.Background()

	post := &models.Post{ID: "1", Title: "Title", Content: "Content"}
	resolver.Storage.CreatePost(post)
	resolver.Storage.CreateComment(&models.Comment{ID: "1", PostID: "1", Content: "Content"})

	comments, err := resolver.Query().Comments(ctx, "1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, comments.Comments, 1)
	assert.Equal(t, "Content", comments.Comments[0].Content)
}
