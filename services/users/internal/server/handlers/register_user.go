package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	UserName  string `json:"userName" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	var req RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Password:  string(hashed),
	}

	if err := u.db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, ErrorResponse{Message: "User already exists"})
	}

	user.StripPassword()

	if err := u.eventPublisher.PublishUserEvent(config.Settings.UserCreatedEvent, user); err != nil {
		// TODO: Handle error
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": user.ID})
}
