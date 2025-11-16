package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/posts/config"
	"github.com/kylehipz/socmed-microservices/posts/internal/routes/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEchoServer(log *zap.Logger, db *gorm.DB, publisher *events.Publisher) *echo.Echo {
	e := echo.New()
	// init handler
	postHandler := handlers.NewPostHandler(db, publisher)

	// middlewares
	loggerMiddleware := middlewares.RequestLoggerMiddleware(log)

	e.Use(loggerMiddleware)

	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)
	authenticatedRoutes := e.Group("", jwtMiddleware)

	authenticatedRoutes.POST("/", postHandler.CreatePost)

	authenticatedRoutes.PATCH("/:id", postHandler.UpdatePost)
	authenticatedRoutes.DELETE("/:id", postHandler.DeletePost)

	e.GET("/healthz", server.HealthCheck)

	return e
}
