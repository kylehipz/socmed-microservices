package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/timeline/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (t *TimelineHandler) GetTimeline(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	userID := c.Get("user_id").(string)
	cursor := c.QueryParam("cursor")
	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		log.Warn("Invalid limit, defaulting to 10")
		limit = 10
	}

	if cursor == "" {
		cursor = time.Now().UTC().Format(time.RFC3339Nano)
	} else {
		// convert cursor to the right format
		// Try to parse the cursor the client sent
    parsed, err := time.Parse(time.RFC3339Nano, cursor)
    if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{ Message: "Invalid cursor format", })
    }
		cursor = parsed.UTC().Format(time.RFC3339Nano)
	}

	// query _args
	queryArgs := []any{userID, cursor, limit}

	query := `
		SELECT p.*, u.id as author_id, u.first_name, u.last_name, u.user_name
		FROM timeline.posts p
		JOIN timeline.follow f ON p.author_id = f.followee_id
		JOIN timeline.users u ON p.author_id = u.id
		WHERE f.follower_id = ? AND p.created_at < ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	result := []models.PostWithAuthor{}

	if err := t.db.Raw(query, queryArgs...).Scan(&result).Error; err != nil {
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	return c.JSON(http.StatusOK, result)
}
