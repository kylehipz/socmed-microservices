package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/routes/handlers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewEchoServer(db *gorm.DB, publisher *events.Publisher) *echo.Echo {
	e := echo.New()

	userHandler := handlers.NewUserHandler(db, publisher)
	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)

	authenticatedRoutes := e.Group("")
	authenticatedRoutes.Use(jwtMiddleware)

	authenticatedRoutes.GET("/", userHandler.ListUsers)
	authenticatedRoutes.GET("/me", userHandler.GetUser)

	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)
	e.GET("/healthz", userHandler.HealthCheck)

	return e
}
