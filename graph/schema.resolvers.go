package graph

import (
	"context"
	"time"

	"fmt"
	"test-graphql/graph/model"
)

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// CreateUser is the resolver for the createUser mutation.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := &model.User{
		ID:   input.ID,
		Name: input.Name,
	}

	// Save the user to the database
	if err := r.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// CreateVideo is the resolver for the createVideo mutation.
func (r *mutationResolver) CreateVideo(ctx context.Context, input model.NewVideo) (*model.Video, error) {
	// Retrieve the User from the database
	var author model.User
	if err := r.DB.First(&author, "id = ?", input.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Create the video entry
	video := &model.Video{
		ID:       fmt.Sprintf("T%d", time.Time.UnixMilli(time.Now())),
		Title:    input.Title,
		URL:      input.URL,
		AuthorID: input.UserID, // Set the AuthorID
		Author:   &author,      // Optionally set the Author field
	}

	if err := r.DB.Create(video).Error; err != nil {
		return nil, fmt.Errorf("failed to create video: %v", err)
	}

	return video, nil
}

// UpdateVideo is the resolver for the updateVideo mutation.
func (r *mutationResolver) UpdateVideo(ctx context.Context, id string, input model.NewVideo) (*model.Video, error) {
	video := &model.Video{}
	if err := r.DB.First(video, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("video not found: %v", err)
	}

	video.Title = input.Title
	video.URL = input.URL

	if err := r.DB.Save(video).Error; err != nil {
		return nil, fmt.Errorf("failed to update video: %v", err)
	}

	return video, nil
}

// DeleteVideo is the resolver for the deleteVideo mutation.
func (r *mutationResolver) DeleteVideo(ctx context.Context, id string) (bool, error) {
	if err := r.DB.Delete(&model.Video{}, "id = ?", id).Error; err != nil {
		return false, fmt.Errorf("failed to delete video: %v", err)
	}
	return true, nil
}

// Videos is the resolver for the videos query.
func (r *queryResolver) Videos(ctx context.Context) ([]*model.Video, error) {
	var videos []*model.Video
	if err := r.DB.Find(&videos).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}
	return videos, nil
}

// User is the resolver for the user query.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	if err := r.DB.First(user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
