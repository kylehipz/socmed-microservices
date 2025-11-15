package handlers

import "gorm.io/gorm"

type FollowHandler struct {
	db *gorm.DB
}

func NewFollowHandler(db *gorm.DB) *FollowHandler {
	return &FollowHandler{db: db}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
