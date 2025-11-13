package routes

import (
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/middlewares"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/routes/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewEchoServer(log *zap.Logger, db *gorm.DB, publisher *events.Publisher) *echo.Echo {
	e := echo.New()

	userHandler := handlers.NewUserHandler(db, publisher)

	// middlewares
	loggerMiddleware := middlewares.RequestLoggerMiddleware(log)
	jwtMiddleware := middlewares.JWTAuth(config.JwtSecret)

	e.Use(loggerMiddleware)
	authenticatedRoutes := e.Group("")
	authenticatedRoutes.Use(jwtMiddleware)

	// Routes
	authenticatedRoutes.GET("/", userHandler.ListUsers)
	authenticatedRoutes.GET("/me", userHandler.GetUser)

	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser)

	e.GET("/healthz", server.HealthCheck)

	return e
}
