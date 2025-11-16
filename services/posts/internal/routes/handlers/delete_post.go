package handlers

import (
	"errors"
	"net/http"

	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/posts/internal/models"
	"github.com/kylehipz/socmed-microservices/posts/internal/types"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (p *PostHandler) DeletePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	authorID := c.Get("user_id").(string)
	id := c.Param("id")

	var post models.Post

	if err := p.db.First(&post, "id = ? AND author_id = ?", id, authorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
        return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Post not found"})
    }

		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	if err := p.db.Delete(&post).Error; err != nil {
		errorMessage := "Internal server error"
		log.Error(errorMessage, zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: errorMessage})
	}

	return c.JSON(http.StatusOK, post)
}
