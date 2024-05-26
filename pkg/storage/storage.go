package storage

import "github.com/lomins/ozon-comments-graphql/pkg/models"

type Storage interface {
	GetPosts() ([]*models.Post, error)
	GetPost(id string) (*models.Post, error)
	CreatePost(post *models.Post) error
	DisableComments(postID string) (*models.Post, error)
	GetComments(postID string, limit int, offset int) ([]*models.Comment, error)
	GetCommentCount(postID string) (int, error)
	CreateComment(comment *models.Comment) error
	Close() error
}
