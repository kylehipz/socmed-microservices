package server

import (
	"github.com/kylehipz/socmed-microservices/users/internal/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/users/internal/server/handlers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewEchoServer(db *gorm.DB, eventPublisher *events.Publisher) *echo.Echo {
	e := echo.New()

	userHandler := handlers.NewUserHandler(db, eventPublisher)
	jwtMiddleware := middlewares.JWTAuth("")

	e.GET("/", jwtMiddleware(userHandler.ListUsers))
	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)
	e.GET("/me", jwtMiddleware(userHandler.GetUser))
	e.GET("/healthz", userHandler.HealthCheck)

	return e
}
