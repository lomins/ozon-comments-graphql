package graph

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/lomins/ozon-comments-graphql/pkg/models"
	"github.com/lomins/ozon-comments-graphql/pkg/storage"
)

func benchmarkCreatePost(b *testing.B, store storage.Storage) {
	for i := 0; i < b.N; i++ {
		post := &models.Post{
			ID:               fmt.Sprintf("%d", time.Now().UnixNano()+rand.Int63()),
			Title:            "Benchmark Post",
			Content:          "This is a benchmark post",
			CommentsDisabled: false,
		}
		store.CreatePost(post)
	}
}

func benchmarkCreateComment(b *testing.B, store storage.Storage) {
	postID := "1"
	store.CreatePost(&models.Post{
		ID:               postID,
		Title:            "Post for Benchmark Comments",
		Content:          "This post is used for benchmarking comments",
		CommentsDisabled: false,
	})

	for i := 0; i < b.N; i++ {
		comment := &models.Comment{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			PostID:    postID,
			Content:   "This is a benchmark comment",
			CreatedAt: time.Now(),
		}
		store.CreateComment(comment)
	}
}

func BenchmarkCreatePostInMemory(b *testing.B) {
	store := storage.NewInMemoryStorage()
	benchmarkCreatePost(b, store)
}

func BenchmarkCreateCommentInMemory(b *testing.B) {
	store := storage.NewInMemoryStorage()
	benchmarkCreateComment(b, store)
}

func BenchmarkCreatePostPostgres(b *testing.B) {
	dsn := "host=localhost port=5432 user=postgres dbname=ozon-comments-test password=7070 sslmode=disable"

	store, err := storage.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Postgres: %s", err)
	}
	defer store.Close()
	benchmarkCreatePost(b, store)
}

func BenchmarkCreateCommentPostgres(b *testing.B) {
	dsn := "host=localhost port=5432 user=postgres dbname=ozon-comments-test password=7070 sslmode=disable"

	store, err := storage.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Postgres: %s", err)
	}
	defer store.Close()
	benchmarkCreateComment(b, store)
}
