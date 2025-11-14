package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/kylehipz/socmed-microservices/users/internal/types"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserHandler) RegisterUser(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	var req types.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error("Invalid payload")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	log.Debug("Password hashed")

	user := models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Password:  string(hashed),
	}

	if err := u.db.Create(&user).Error; err != nil {
		errorMessage := "User already exists"
		log.Error(errorMessage, zap.Error(err))

		return c.JSON(http.StatusConflict, ErrorResponse{Message: errorMessage})
	}

	user.StripPassword()
	log.Debug("Password stripped")

	if err := u.publisher.PublishEvent(events.UserCreated, user); err != nil {
		log.With(zap.String("event_name", events.UserCreated)).Error("Failed to publish event")
	} else {
		log.With(zap.String("event_name", events.UserCreated)).Debug("Event published successfully")
	}

	return c.JSON(http.StatusCreated, echo.Map{"id": user.ID})
}
