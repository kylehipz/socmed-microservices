package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/timeline/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEchoServer(log *zap.Logger, db *gorm.DB) *echo.Echo {
	e := echo.New()

	// middlewares
	loggerMiddleware := middlewares.RequestLoggerMiddleware(log)


	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)
	e.Use(loggerMiddleware, jwtMiddleware)
	// authenticatedRoutes := e.Group("", jwtMiddleware)

	// authenticatedRoutes.POST("/", postHandler.CreatePost)
	//
	// authenticatedRoutes.PATCH("/:id", postHandler.UpdatePost)
	// authenticatedRoutes.DELETE("/:id", postHandler.DeletePost)

	e.GET("/healthz", server.HealthCheck)

	return e
}
