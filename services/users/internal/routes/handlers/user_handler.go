package handlers

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"gorm.io/gorm"
)

type UserHandler struct {
	db        *gorm.DB
	publisher *events.Publisher
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewUserHandler(db *gorm.DB, publisher *events.Publisher) *UserHandler {
	return &UserHandler{db: db, publisher: publisher}
}
