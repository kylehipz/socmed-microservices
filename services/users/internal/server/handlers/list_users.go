package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/labstack/echo/v4"
)

func (u *UserHandler) ListUsers(c echo.Context) error {
	var users []models.User

	if err := u.db.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, users)
}
