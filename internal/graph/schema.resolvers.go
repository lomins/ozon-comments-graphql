package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lomins/ozon-comments-graphql/internal/graph/model"
	"github.com/lomins/ozon-comments-graphql/pkg/models"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsDisabled bool) (*model.Post, error) {
	post := &model.Post{
		ID:               fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
	}

	postGorm := models.ToGORMPost(post)
	err := r.Storage.CreatePost(postGorm)
	return post, err
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	postGorm, err := r.Storage.GetPost(postID)
	if err != nil {
		return nil, err
	}

	if postGorm.CommentsDisabled {
		return nil, errors.New("comments are disabled for this post")
	}

	commentQL := &model.Comment{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	commentGorm := models.ToGORMComment(commentQL)
	err = r.Storage.CreateComment(commentGorm)

	r.mu.Lock()
	if ch, ok := r.Comments[postID]; ok {
		ch <- commentQL
	}
	r.mu.Unlock()

	return commentQL, err
}

// DisableComments is the resolver for the disableComments field.
func (r *mutationResolver) DisableComments(ctx context.Context, postID string) (*model.Post, error) {
	postGorm, err := r.Storage.DisableComments(postID)
	if err != nil {
		return nil, err
	}
	commentsQL, err := r.Storage.GetComments(postID)
	if err != nil {
		return nil, err
	}

	postQL := models.ToGraphQLPost(postGorm, commentsQL)

	return postQL, err
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	postsGorm, err := r.Storage.GetPosts()
	if err != nil {
		return nil, err
	}

	postsQL := make([]*model.Post, len(postsGorm))

	for i, postGorm := range postsGorm {
		comments, err := r.Storage.GetComments(postGorm.ID)
		if err != nil {
			return nil, err
		}
		postsQL[i] = models.ToGraphQLPost(postGorm, comments)
	}

	return postsQL, err
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	postGorm, err := r.Storage.GetPost(id)
	if err != nil {
		return nil, err
	}

	comments, err := r.Storage.GetComments(postGorm.ID)
	if err != nil {
		return nil, err
	}

	postQL := models.ToGraphQLPost(postGorm, comments)
	return postQL, err
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string) ([]*model.Comment, error) {
	commentsGorm, err := r.Storage.GetComments(postID)
	if err != nil {
		return nil, err
	}

	commentsQL := make([]*model.Comment, len(commentsGorm))

	for i, commentGorm := range commentsGorm {
		commentsQL[i] = models.ToGraphQLComment(commentGorm)
	}

	return commentsQL, err
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Comments[postID]; !ok {
		r.Comments[postID] = make(chan *model.Comment)
	}
	return r.Comments[postID], nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
// func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
// 	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
// }
// func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
// 	panic(fmt.Errorf("not implemented: Todos - todos"))
// }