package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/labstack/echo/v4"
)

func (u *UserHandler) CreateUser(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid payload"})
	}

	if err := u.db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, user)
}
