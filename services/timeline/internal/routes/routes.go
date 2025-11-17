package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/timeline/config"
	"github.com/kylehipz/socmed-microservices/timeline/internal/routes/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEchoServer(log *zap.Logger, db *gorm.DB) *echo.Echo {
	e := echo.New()

	timelineHandler := handlers.NewTimelineHandler(db)

	// middlewares
	loggerMiddleware := middlewares.RequestLoggerMiddleware(log)
	e.Use(loggerMiddleware)

	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)
	authenticatedRoutes := e.Group("", jwtMiddleware)

	authenticatedRoutes.GET("/", timelineHandler.GetTimeline)
	e.GET("/healthz", server.HealthCheck)

	return e
}
