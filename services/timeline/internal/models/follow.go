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

func (Follow) TableName() string {
	return "timeline.follow"
}


