package storage

import (
	"errors"

	"sync"

	"github.com/lomins/ozon-comments-graphql/pkg/models"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	posts    map[string]*models.Post
	comments map[string][]*models.Comment
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    make(map[string]*models.Post),
		comments: make(map[string][]*models.Comment),
	}
}

func (s *InMemoryStorage) GetPosts() ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *InMemoryStorage) GetPost(id string) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (s *InMemoryStorage) CreatePost(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryStorage) DisableComments(postID string) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[postID]
	if !exists {
		return nil, errors.New("post not found")
	}

	post.CommentsDisabled = true
	return post, nil
}

func (s *InMemoryStorage) GetComments(postID string, limit int, offset int) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments, exists := s.comments[postID]
	if !exists {
		return nil, nil
	}

	if offset >= len(comments) {
		return []*models.Comment{}, nil
	}

	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}

	return comments[offset:end], nil
}

func (s *InMemoryStorage) GetCommentCount(postID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments, exists := s.comments[postID]
	if !exists {
		return 0, nil
	}

	return len(comments), nil
}

func (s *InMemoryStorage) CreateComment(comment *models.Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.comments[comment.PostID] = append(s.comments[comment.PostID], comment)
	return nil
}
