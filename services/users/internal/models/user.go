package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a row in the users.users table.
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName string         `gorm:"type:text;not null"`
	LastName  string         `gorm:"type:text;not null"`
	Email     string         `gorm:"type:text;not null;uniqueIndex"`
	Username  string         `gorm:"type:text;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // optional soft delete
}

// TableName overrides the default pluralized name (users)
// so GORM uses the users.users schema-qualified table.
func (User) TableName() string {
	return "users.users"
}
