package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a row in the users.users table.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName string    `gorm:"type:text;not null" json:"firstName"`
	LastName  string    `gorm:"type:text;not null" json:"lastName"`
	Email     string    `gorm:"type:text;not null;uniqueIndex" json:"email"`
	UserName  string    `gorm:"type:text;not null;uniqueIndex" json:"userName"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"updatedAt"`
}

// TableName overrides the default pluralized name (users)
// so GORM uses the users.users schema-qualified table.
func (User) TableName() string {
	return "follow.users"
}
