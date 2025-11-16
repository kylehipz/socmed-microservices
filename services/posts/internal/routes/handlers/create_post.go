package handlers

import (
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/posts/internal/models"
	"github.com/kylehipz/socmed-microservices/posts/internal/types"
	"github.com/kylehipz/socmed-microservices/posts/internal/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (p *PostHandler) CreatePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	authorID := c.Get("user_id").(string)

	var req types.CreateOrUpdatePostRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid payload"})
	}

	post := models.Post{
		AuthorID: utils.ToUUID(authorID),
		Content:  req.Content,
	}

	if err := p.db.Create(&post).Error; err != nil {
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	log_event_field := zap.String("event_name", events.PostCreated)
	if err := p.publisher.PublishEvent(events.PostCreated, post); err != nil {
		log.Error("Failed to publish event", log_event_field, zap.Error(err))
	}

	return c.JSON(http.StatusCreated, post)
}
