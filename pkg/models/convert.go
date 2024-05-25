// pkg/models/convert.go
package models

import (
	"time"

	"github.com/lomins/ozon-comments-graphql/internal/graph/model"
)

func ToGraphQLPost(post *Post, comments []*Comment) *model.Post {
	gqlComments := make([]*model.Comment, len(comments))
	for i, c := range comments {
		gqlComments[i] = ToGraphQLComment(c)
	}

	return &model.Post{
		ID:               post.ID,
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
		Comments:         gqlComments,
	}
}

func ToGraphQLComment(comment *Comment) *model.Comment {
	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}
}

func ToGORMPost(post *model.Post) *Post {
	return &Post{
		ID:               post.ID,
		Title:            post.Title,
		Content:          post.Content,
		CommentsDisabled: post.CommentsDisabled,
	}
}

func ToGORMComment(comment *model.Comment) *Comment {
	return &Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: time.Now(), // you can adjust this according to your needs
	}
}
