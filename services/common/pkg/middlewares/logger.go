package middlewares

import (
	"time"

	"github.com/google/uuid"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLoggerMiddleware(baseLogger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			// skip if health check
			if req.URL.Path == "/healthz" || req.URL.Path == "/ready" {
				return next(c)
			}

			// Generate a request ID (can be replaced by your own system)
			requestID := uuid.New().String()

			// Create request-scoped logger
			reqLogger := baseLogger.With(
				zap.String("request_id", requestID),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
			)

			// Put logger into context
			ctx := logger.WithContext(req.Context(), reqLogger)
			c.SetRequest(req.WithContext(ctx))

			// Execute handler
			err := next(c)

			// After handling
			latency := time.Since(start)

			status := res.Status
			size := res.Size

			fields := []zap.Field{
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.Int64("response_size", size),
			}

			if err != nil {
				// Let Echo handle transforming error → HTTP response
				c.Error(err)
				reqLogger.Error("request completed with error", append(fields, zap.Error(err))...)
			} else {
				reqLogger.Info("request completed", fields...)
			}

			return nil
		}
	}
}
