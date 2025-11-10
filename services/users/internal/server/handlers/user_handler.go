package handlers

import "gorm.io/gorm"

type UserHandler struct {
	db *gorm.DB
	jwtSecret string
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewUserHandler(db *gorm.DB) *UserHandler{ 
	return &UserHandler{db: db}
}
