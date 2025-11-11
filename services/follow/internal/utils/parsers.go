package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kylehipz/socmed-microservices/follow/internal/models"
)

func ParseFollow(followerVal any, followeeStr string) (*models.Follow, error) {
	// Get follower ID from context
	followerStr, ok := followerVal.(string)
	if !ok {
		return nil, errors.New("Invalid user ID")
	}

	followerID, err := uuid.Parse(followerStr)
	if err != nil {
		return nil, err
	}

	// Get followee ID from route param
	followeeID, err := uuid.Parse(followeeStr)
	if err != nil {
		return nil, err
	}

	return &models.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}, nil
}
