package server

import (
	"github.com/kylehipz/socmed-microservices/follow/internal/middlewares"
	"github.com/kylehipz/socmed-microservices/follow/internal/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)


func NewEchoServer(db *gorm.DB) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	followHandler := handlers.NewFollowHandler(db)
	jwtMiddleware := middlewares.JWTAuth("")

	e.POST("/:id", jwtMiddleware(followHandler.FollowUser))
	e.DELETE("/:id", jwtMiddleware(followHandler.UnfollowUser))

	e.GET("/healthz", followHandler.HealthCheck)

	return e
}
