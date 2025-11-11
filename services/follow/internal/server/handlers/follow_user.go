package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/follow/internal/utils"
	"github.com/labstack/echo/v4"
)

func (f *FollowHandler) FollowUser(c echo.Context) error {
	followerVal := c.Get("user_id")
	followeeStr := c.Param("id")

	follow, err := utils.ParseFollow(followerVal, followeeStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	if err := f.db.Create(follow).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, *follow)
}
