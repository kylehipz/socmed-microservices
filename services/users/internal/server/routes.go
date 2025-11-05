package server

import (
	"github.com/kylehipz/socmed-microservices/users/internal/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)


func NewEchoServer(db *gorm.DB) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userHandler := handlers.NewUserHandler(db)

	e.GET("/", userHandler.ListUsers)
	e.POST("/", userHandler.CreateUser)
	e.GET("/healthz", userHandler.HealthCheck)

	return e
}
