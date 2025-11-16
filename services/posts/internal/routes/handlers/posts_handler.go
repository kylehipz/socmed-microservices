package handlers

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"gorm.io/gorm"
)

type PostHandler struct {
	db        *gorm.DB
	publisher *events.Publisher
}

func NewPostHandler(db *gorm.DB, publisher *events.Publisher) *PostHandler {
	return &PostHandler{db: db, publisher: publisher}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
