package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/follow/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (f *FollowHandler) FollowUser(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	followerVal := c.Get("user_id")
	followeeStr := c.Param("id")

	follow, err := utils.ParseFollow(followerVal, followeeStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	if err := f.db.Create(follow).Error; err != nil {
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	return c.JSON(http.StatusCreated, *follow)
}
