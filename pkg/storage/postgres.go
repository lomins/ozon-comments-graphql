package storage

import (
	"fmt"

	"github.com/lomins/ozon-comments-graphql/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresStorage struct {
	DB *gorm.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	err = db.AutoMigrate(models.Post{}, models.Comment{})
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить миграцию: %w", err)
	}
	return &PostgresStorage{DB: db}, nil
}

func (s *PostgresStorage) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("не удалось получить объект базы данных: %w", err)
	}
	return sqlDB.Close()
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

func (s *PostgresStorage) GetComments(postID string, limit int, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := s.DB.Where("post_id = ?", postID).Order("created_at ASC").Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}

func (s *PostgresStorage) GetCommentCount(postID string) (int, error) {
	var count int64
	err := s.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}

func (s *PostgresStorage) CreateComment(comment *models.Comment) error {
	if comment.ParentID != nil {
		var parentComment models.Comment
		if err := s.DB.First(&parentComment, "id = ?", *comment.ParentID).Error; err != nil {
			return fmt.Errorf("родительский комментарий не найден: %w", err)
		}
	}

	if err := comment.Validate(); err != nil {
		return err
	}
	return s.DB.Create(comment).Error
}
