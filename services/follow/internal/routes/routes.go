package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/routes/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEchoServer(log *zap.Logger, db *gorm.DB, publisher *events.Publisher) *echo.Echo {
	e := echo.New()

	followHandler := handlers.NewFollowHandler(db, publisher)
	loggerMiddleware := middlewares.RequestLoggerMiddleware(log)

	e.Use(loggerMiddleware)

	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)
	authenticatedRoutes := e.Group("", jwtMiddleware)

	authenticatedRoutes.POST("/:id", followHandler.FollowUser)
	authenticatedRoutes.DELETE("/:id", followHandler.UnfollowUser)

	e.GET("/healthz", followHandler.HealthCheck)

	return e
}
