package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a row in the posts.posts table.
type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null" json:"authorID"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"updatedAt"`

	// Optional associations if you want GORM to preload
	Author *User `gorm:"foreignKey:AuthorID;references:ID" json:"author,omitempty"`
}

// TableName overrides the default pluralized name (users)
// so GORM uses the users.users schema-qualified table.
func (Post) TableName() string {
	return "timeline.posts"
}

type PostWithAuthor struct {
	Post
	UserName string `gorm:"type:text;not null" json:"userName"`
	FirstName string `gorm:"type:text;not null" json:"firstName"`
	LastName string `gorm:"type:text;not null" json:"lastName"`
}

