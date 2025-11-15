package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/constants"
	"github.com/kylehipz/socmed-microservices/common/pkg/db"
	err_utils "github.com/kylehipz/socmed-microservices/common/pkg/errors"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/logger"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/routes"
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

	publisher := events.NewPublisher(ch, constants.SocmedExchangeName)

	// init db
	gormDB, err := db.NewGormDB(log, config.DatabaseUrl)
	err_utils.HandleFatalError(log, err)

	// init echo and API Server
	e := routes.NewEchoServer(log, gormDB, publisher)
	apiServer := server.NewApiServer(
		log,
		e,
		config.ServiceName,
		nil,
		gormDB,
		rabbitMqConn,
	)

	go apiServer.Start(mainCtx, stop, config.HttpPort)

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
