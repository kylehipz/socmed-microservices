package errors

import (
	"go.uber.org/zap"
)

func HandleFatalError(log *zap.Logger, err error) {
	if err != nil {
		log.Fatal("Fatal error", zap.Error(err))
	}
}
