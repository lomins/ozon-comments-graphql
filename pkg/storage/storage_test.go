package storage_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("could not connect to the database: %v", err)
	}
	defer db.Close()

	store := storage.NewPostgresStorage(db)

	db.AutoMigrate(&models.Post{}, &models.Comment{})

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

	db.Delete(&models.Post{}, "id = ?", "1")
	db.Delete(&models.Comment{}, "post_id = ?", "1")
}
