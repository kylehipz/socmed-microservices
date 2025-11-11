package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"follower_id"`
	FolloweeID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"followee_id"`
	CreatedAt  time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`

	// Optional associations if you want GORM to preload
	Follower *User `gorm:"foreignKey:FollowerID;references:ID" json:"follower,omitempty"`
	Followee *User `gorm:"foreignKey:FolloweeID;references:ID" json:"followee,omitempty"`
}

type FollowResponse struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FolloweeID uuid.UUID `json:"followee_id"`

	FollowerFirstName string `json:"follower_first_name"`
	FolloweeFirstName string `json:"followee_first_name"`

	FollowerLastName string `json:"follower_last_name"`
	FolloweeLastName string `json:"followee_last_name"`

	FollowerUserName string `json:"follower_username"`
	FolloweeUserName string `json:"followee_username"`

	FollowerEmail string `json:"follower_email"`
	FolloweeEmail string `json:"followee_email"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Follow) TableName() string {
	return "follow.follow"
}

func (f *Follow) ToFollowResponse() *FollowResponse {
	return &FollowResponse{
		// Follower
		FollowerID: f.FollowerID,
		FollowerLastName: f.Follower.LastName,
		FollowerFirstName: f.Follower.FirstName,
		FollowerUserName: f.Follower.FirstName,
		FollowerEmail: f.Follower.Email,

		// Followee
		FolloweeID: f.FolloweeID,
		FolloweeLastName: f.Followee.LastName,
		FolloweeFirstName: f.Followee.FirstName,
		FolloweeUserName: f.Followee.FirstName,
		FolloweeEmail: f.Followee.Email,

		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
