package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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

func (a *ApiServer) Start(ctx context.Context,  port string) error {
	// Start consumers
	if err := a.startConsumers(ctx); err != nil {
		a.log.Error("Critical consumer failed to start", zap.Error(err))
		return err
	}

	// Start API Server
	if err := a.startHttpServer(port); err != nil {
		a.log.Error("Http server failed to start", zap.Error(err))
		return err
	}

	return nil
}

func (a *ApiServer) Wait(ctx context.Context) {
	a.log.Info("App shutting down...")

	var wg sync.WaitGroup
	wg.Add(2)

	// drain all consumers
	go func() {
		defer wg.Done()
		a.stopConsumers(ctx)
	}()

	// drain http requests
	go func() {
		defer wg.Done()
		a.stopHttpServer(ctx)
	}()

	wg.Wait()

	// close all connections
	a.closeConnections()
	a.log.Info("API server shutdown complete")
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
	a.log.Info("Attempting HTTP server graceful shutdown...")
	httpServerShutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.e.Shutdown(httpServerShutdownCtx); err != nil {
		a.log.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		a.log.Info("HTTP server shutdown successfully")
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

func (a *ApiServer) startConsumers(ctx context.Context) error {
	consumerCount := len(a.consumers)

	if consumerCount == 0 {
		a.log.Warn("No registered consumers")
	}

	for _, c := range a.consumers {
			if err := c.Start(ctx); err != nil {
				return err
			}
	}

	return nil
}

func (a *ApiServer) stopConsumers(ctx context.Context) {
	a.log.Info("Attempting event consumers graceful shutdown...")
	// setup consumer shutdown context: 5s
	consumersShutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// wait
	var wg sync.WaitGroup
	wg.Add(len(a.consumers))

	for _, consumer := range a.consumers {
		go func(c events.Consumer) {
			defer wg.Done()
			c.Wait(consumersShutdownCtx)
		}(consumer)
	}

	// wait for all consumers to be done
	waitConsumersChan := make(chan struct{})

	go func() {
		defer close(waitConsumersChan)
		wg.Wait()
	}()

	select {
	case <-waitConsumersChan:
		a.log.Info("Consumers shutdown successfully")
	case <-consumersShutdownCtx.Done():
		a.log.Info("Consumers shutdown timed out. Shutdown forced")
	}
}
