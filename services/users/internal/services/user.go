package services

import (
	"github.com/kylehipz/socmed-microservices/users/internal/events"
	"gorm.io/gorm"
)

type UsersService struct {
	db *gorm.DB
	publisher *events.Publisher
}

func NewUsersService(db *gorm.DB, publisher *events.Publisher) *UsersService {
	return &UsersService{db: db, publisher: publisher}
}
