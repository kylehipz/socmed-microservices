package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/constants"
	"github.com/kylehipz/socmed-microservices/common/pkg/db"
	err_utils "github.com/kylehipz/socmed-microservices/common/pkg/errors"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/events/consumers"
	"github.com/kylehipz/socmed-microservices/follow/internal/routes"
	"go.uber.org/zap"
)

func main() {
	// init logger
	log := logger.NewLogger(config.Environment, config.LogLevel)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// init rabbitmq
	rabbitMqConn, err := events.NewRabbitMQConn(log, config.RabbitMqUrl)
	err_utils.HandleFatalError(log, err)

	ch, err := rabbitMqConn.Channel()
	err_utils.HandleFatalError(log, err)

	defer ch.Close()

	publisher := events.NewPublisher(ch, constants.SocmedExchangeName)

	// init db
	gormDB, err := db.NewGormDB(log, config.DatabaseUrl)
	err_utils.HandleFatalError(log, err)
	
	// init consumers
	userEventsConsumer := consumers.NewUserEventsConsumer(log, ch, gormDB, 10)
	consumers := []events.Consumer{userEventsConsumer}

	// init echo and API Server
	e := routes.NewEchoServer(log, gormDB, publisher)
	apiServer := server.NewApiServer(
		log,
		e,
		config.ServiceName,
		consumers,
		gormDB,
		rabbitMqConn,
	)

	go func() {
		if err := apiServer.Start(mainCtx, config.HttpPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("API server error", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-mainCtx.Done()

	// Call stop to remove the signal handler
	stop()
	log.Info("Shutdown signal received. Starting graceful shutdown...")

	// Drain API Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	apiServer.Stop(shutdownCtx)

	log.Info("Application shutdown gracefully...")
}
