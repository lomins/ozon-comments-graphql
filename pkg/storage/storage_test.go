package storage_test

import (
	"log"
	"testing"

	"github.com/lomins/ozon-comments-graphql/pkg/models"
	"github.com/lomins/ozon-comments-graphql/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage(t *testing.T) {
	store := storage.NewInMemoryStorage()

	post := &models.Post{
		ID:    "1",
		Title: "Test Post",
	}

	err := store.CreatePost(post)
	assert.NoError(t, err)

	fetchedPost, err := store.GetPost("1")
	assert.NoError(t, err)
	assert.Equal(t, post, fetchedPost)

	err = store.CreateComment(&models.Comment{
		ID:      "1",
		PostID:  "1",
		Content: "Test Comment",
	})
	assert.NoError(t, err)

	comments, err := store.GetComments("1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
}

func TestPostgresStorage(t *testing.T) {

	dsn := "host=localhost port=5432 user=postgres dbname=ozon-comments password=7070 sslmode=disable"

	store, err := storage.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Postgres: %s", err)
	}
	defer store.Close()

	post := &models.Post{
		ID:       "1",
		Title:    "Test Post",
		Comments: []*models.Comment{},
	}

	err = store.CreatePost(post)
	assert.NoError(t, err)

	fetchedPost, err := store.GetPost("1")
	assert.NoError(t, err)
	assert.Equal(t, post, fetchedPost)

	err = store.CreateComment(&models.Comment{
		ID:      "1",
		PostID:  "1",
		Content: "Test Comment",
	})
	assert.NoError(t, err)

	comments, err := store.GetComments("1", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
}
