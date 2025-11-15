package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/follow/internal/utils"
	"github.com/labstack/echo/v4"
)

func (f *FollowHandler) UnfollowUser(c echo.Context) error {
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

	return c.NoContent(http.StatusNoContent)
}
