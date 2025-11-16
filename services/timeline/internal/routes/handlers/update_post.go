package handlers

import (
	"errors"
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/posts/internal/models"
	"github.com/kylehipz/socmed-microservices/posts/internal/types"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (p *PostHandler) UpdatePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	authorID := c.Get("user_id").(string)
	id := c.Param("id")

	var req types.CreateOrUpdatePostRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	var post models.Post

	if err := p.db.First(&post, "id = ? AND author_id = ?", id, authorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Post not found"})
		}

		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	post.Content = req.Content

	if err := p.db.Save(&post).Error; err != nil {
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	log_event_field := zap.String("event_name", events.PostUpdated)
	if err := p.publisher.PublishEvent(events.PostUpdated, post); err != nil {
		log.Error("Failed to publish event", log_event_field, zap.Error(err))
	}

	return c.JSON(http.StatusOK, post)
}
