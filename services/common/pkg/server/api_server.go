package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	err_utils "github.com/kylehipz/socmed-microservices/common/pkg/errors"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApiServer struct {
	log            *zap.Logger
	e              *echo.Echo
	name           string
	consumers      []events.Consumer
	db             *gorm.DB
	mq             *amqp091.Connection
	wg             sync.WaitGroup
	consumerCancel context.CancelFunc
}

func NewApiServer(
	log *zap.Logger,
	e *echo.Echo,
	name string,
	consumers []events.Consumer,
	db *gorm.DB,
	mq *amqp091.Connection,
) *ApiServer {
	e.HideBanner = true
	e.HidePort = true

	return &ApiServer{
		log:       log,
		e:         e,
		name:      name,
		consumers: consumers,
		db:        db,
		mq:        mq,
	}
}

func (a *ApiServer) Start(ctx context.Context, port string) error {
	// Consumer context
	consumerCtx, cancel := context.WithCancel(ctx)
	a.consumerCancel = cancel

	// Start consumers
	a.startConsumers(consumerCtx)

	// Start API Server
	a.startHttpServer(port)

	return nil
}

func (a *ApiServer) Stop(ctx context.Context) {
	a.log.Info("App shutting down...")

	if a.consumerCancel != nil {
		a.consumerCancel()
	}

	// drain http requests
	a.stopHttpServer(ctx)

	// close all connections
	a.closeConnections()
	a.log.Info("Application shutdown complete")
}

func (a *ApiServer) startHttpServer(port string) error {
	a.e.HideBanner = true
	a.e.HidePort = true

	portStr := fmt.Sprintf(":%s", port)

	a.log.Info(fmt.Sprintf("%s started on port %s", a.name, port))
	if err := a.e.Start(portStr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *ApiServer) stopHttpServer(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	if err := a.e.Shutdown(shutdownCtx); err != nil {
		a.log.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		a.log.Info("HTTP server shut down gracefully")
	}

	waitChan := make(chan struct{})
	go func() {
		defer close(waitChan)
		// wait all consumers to be done
		a.wg.Wait()
	}()

	select {
	case <-waitChan:
		a.log.Info("All consumers stopped gracefully")
	case <-ctx.Done():
		a.log.Info("Consumers timed out, proceeding with shutdown")
	}
}

func (a *ApiServer) closeConnections() {
	if a.db != nil {
		sqlDb, _ := a.db.DB()
		err := sqlDb.Close()
		if err != nil {
			a.log.Error("Cannot close database connection", zap.Error(err))
		} else {
			a.log.Info("Database connection closed")
		}
	}

	if a.mq != nil {
		if err := a.mq.Close(); err != nil {
			a.log.Error("Cannot close rabbitmq connection", zap.Error(err))
		} else {
			a.log.Info("RabbitMQ connection closed.")
		}
	}
}

func (a *ApiServer) startConsumers(ctx context.Context) {
	a.wg.Add(len(a.consumers))
	for _, c := range a.consumers {
		go func(consumer events.Consumer) {
			defer a.wg.Done()
			if err := consumer.Start(ctx); err != nil {
				err_utils.HandleFatalError(a.log, err)
			}
		}(c)
	}
}
