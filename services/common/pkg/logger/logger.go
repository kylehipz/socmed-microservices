package logger

import (
	"context"
	"strings"

	"github.com/kylehipz/socmed-microservices/common/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const loggerKey = contextKey("logger")

func NewLogger(env string, logLevel string) *zap.Logger {
	var cfg zap.Config

	// Pick base config (dev has pretty console logs, prod has JSON logs)
	if env == constants.Development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Resolve the desired level
	level := resolveLogLevel(env, logLevel)
	cfg.Level = zap.NewAtomicLevelAt(level)

	cfg.DisableCaller = false
	cfg.DisableStacktrace = false

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := cfg.Build()

	return logger
}

func resolveLogLevel(env string, override string) zapcore.Level {
	// Manual override takes priority
	if override != "" {
		switch strings.ToLower(override) {
		case "debug":
			return zap.DebugLevel
		case "info":
			return zap.InfoLevel
		case "warn":
			return zap.WarnLevel
		case "error":
			return zap.ErrorLevel
		}
	}

	// Default behavior based on environment
	if env == constants.Development {
		return zap.DebugLevel
	}

	return zap.InfoLevel
}

func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	if log, ok := ctx.Value(loggerKey).(*zap.Logger); ok && log != nil {
		return log
	}

	return zap.L()
}
