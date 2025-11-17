package handlers

import "gorm.io/gorm"

type TimelineHandler struct {
	db        *gorm.DB
}

func NewTimelineHandler(db *gorm.DB) *TimelineHandler {
	return &TimelineHandler{db: db}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
