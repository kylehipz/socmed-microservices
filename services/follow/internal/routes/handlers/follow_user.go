package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/follow/internal/models"
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

	err = f.db.Create(&follow).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				errorMessage := "User is already followed"
				log.Error(errorMessage, zap.Error(err))
				return c.JSON(http.StatusBadRequest, ErrorResponse{Message: errorMessage})
			}
		} 
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	var result models.Follow
	if err := f.db.
		Preload("Follower").Preload("Followee").
		First(&result, "follower_id = ? AND followee_id = ?", follow.FollowerID, follow.FolloweeID).
		Error; err != nil {
		log.Error("Load follow associations failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
	}

	return c.JSON(http.StatusCreated, result.ToFollowResponse())
}
