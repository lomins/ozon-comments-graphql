package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/lomins/ozon-comments-graphql/pkg/models"
)

type PostgresStorage struct {
	DB *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *PostgresStorage {
	return &PostgresStorage{DB: db}
}

func (s *PostgresStorage) GetPosts() ([]*models.Post, error) {
	var posts []*models.Post
	err := s.DB.Preload("Comments").Find(&posts).Error
	return posts, err
}

func (s *PostgresStorage) GetPost(id string) (*models.Post, error) {
	var post models.Post
	err := s.DB.Preload("Comments").First(&post, "id = ?", id).Error
	return &post, err
}

func (s *PostgresStorage) CreatePost(post *models.Post) error {
	return s.DB.Create(post).Error
}

func (s *PostgresStorage) DisableComments(postID string) (*models.Post, error) {
	var post models.Post
	if err := s.DB.First(&post, "id = ?", postID).Error; err != nil {
		return nil, err
	}

	post.CommentsDisabled = true
	return &post, s.DB.Save(&post).Error
}

func (s *PostgresStorage) GetComments(postID string) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := s.DB.Where("post_id = ?", postID).Order("created_at ASC")

	err := query.Find(&comments).Error
	return comments, err
}

func (s *PostgresStorage) CreateComment(comment *models.Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}
	return s.DB.Create(comment).Error
}
