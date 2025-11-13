package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/kylehipz/socmed-microservices/users/internal/types"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserHandler) RegisterUser(c echo.Context) error {
	var req types.RegisterUserRequest

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

	if err := u.publisher.PublishEvent(events.UserCreated, user); err != nil {
		// TODO: Handle error
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": user.ID})
}
