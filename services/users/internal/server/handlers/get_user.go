package handlers

import (
	"errors"
	"net/http"

	"github.com/kylehipz/socmed-microservices/users/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (u *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
	}

	return c.JSON(http.StatusOK, user)
}
