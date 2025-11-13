package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ApiServer struct {
	e *echo.Echo
	name string
	consumers []events.Consumer
	db *gorm.DB
	mq *amqp091.Connection
	wg sync.WaitGroup
	consumerCancel context.CancelFunc
}

func NewApiServer(
	e *echo.Echo,
	name string,
	consumers []events.Consumer,
	db *gorm.DB,
	mq *amqp091.Connection,
) *ApiServer {
	e.HideBanner = true
	return &ApiServer{
		e: e,
		name: name,
		consumers: consumers,
		db: db,
		mq: mq,
	}
}

func (a *ApiServer) Start(ctx context.Context, port int) error {
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
	log.Println("App shutting down...")

	if a.consumerCancel != nil {
		a.consumerCancel()
	}

	// drain http requests
	a.stopHttpServer(ctx)

	// close all connections
	a.closeConnections()
	log.Println("Application shutdown complete")
}

func (a *ApiServer) startHttpServer(port int) error {
	portStr := fmt.Sprintf(":%d", port)
	if err := a.e.Start(portStr); err != nil && !errors.Is(err ,http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *ApiServer) stopHttpServer(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	if err := a.e.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shut down gracefully")
	}

	waitChan := make(chan struct{})
	go func () {
		defer close(waitChan)
		// wait all consumers to be done
		a.wg.Wait()
	}()

	select {
	case <-waitChan:
		log.Println("All consumers stopped gracefully")
	case <-ctx.Done():
	log.Println("Consumers timed out, proceeding with shutdown")
	}
}

func (a *ApiServer) closeConnections() {
	if a.db != nil {
		sqlDb, _ := a.db.DB()
		err := sqlDb.Close()
		if err != nil {
			log.Printf("Cannot close database connection %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	if a.mq != nil {
		if err := a.mq.Close(); err != nil {
			log.Printf("Cannot close rabbitmq connection %v", err)
		} else {
			log.Println("RabbitMQ connection closed.")
		}
	}
}

func (a *ApiServer) startConsumers(ctx context.Context) {
	a.wg.Add(len(a.consumers))
	for _, c := range a.consumers {
		go func(consumer events.Consumer) {
			defer a.wg.Done()
			consumer.Start(ctx)
		}(c)
	}
}
