package server

import (
	"github.com/kylehipz/socmed-microservices/users/internal/middlewares"
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
	jwtMiddleware := middlewares.JWTAuth("")

	e.GET("/", jwtMiddleware(userHandler.ListUsers))
	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)
	e.GET("/me", jwtMiddleware(userHandler.GetUser))
	e.GET("/healthz", userHandler.HealthCheck)

	return e
}
