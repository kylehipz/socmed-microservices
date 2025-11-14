package errors

import (
	"go.uber.org/zap"
)

func HandleFatalError(log *zap.Logger, err error) {
	if err != nil {
		log.Error("Error", zap.Error(err))
	}
}
