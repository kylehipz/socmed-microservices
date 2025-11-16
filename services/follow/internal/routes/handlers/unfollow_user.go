package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/follow/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (f *FollowHandler) UnfollowUser(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	followerVal := c.Get("user_id")
	followeeStr := c.Param("id")

	follow, err := utils.ParseFollow(followerVal, followeeStr)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	tx := f.db.Delete(follow, "follower_id = ? AND followee_id = ?", follow.FollowerID, follow.FolloweeID)
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
	}

	if tx.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Follow relationship not found"})
	}

	log_event_field := zap.String("event_name", events.UserUnfollowed)
	if err := f.publisher.PublishEvent(events.UserUnfollowed, follow); err != nil {
		log.Error("Failed to publish event", log_event_field, zap.Error(err))
	}

	return c.NoContent(http.StatusNoContent)
}
