package models

import (
	"errors"
	"time"
)

const MaxCommentLength = 2000

type Comment struct {
	ID        string `gorm:"primary_key"`
	PostID    string
	ParentID  *string
	Content   string
	Children  []*Comment
	CreatedAt time.Time
}

func (c *Comment) Validate() error {
	if len(c.Content) > MaxCommentLength {
		return errors.New("comment content exceeds maximum length")
	}
	return nil
}
