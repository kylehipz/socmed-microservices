package handlers

import (
	"github.com/kylehipz/socmed-microservices/users/internal/events"
	"gorm.io/gorm"
)

type UserHandler struct {
	db             *gorm.DB
	eventPublisher *events.Publisher
	jwtSecret      string
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewUserHandler(db *gorm.DB, eventPublisher *events.Publisher) *UserHandler {
	return &UserHandler{db: db, eventPublisher: eventPublisher}
}
