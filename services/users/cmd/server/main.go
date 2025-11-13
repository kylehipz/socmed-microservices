package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/db"
	err_utils "github.com/kylehipz/socmed-microservices/common/pkg/errors"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/common/pkg/server"
	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/routes"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// init rabbitmq
	rabbitMqConn, err := events.NewRabbitMQConn(config.RabbitMqUrl)
	err_utils.HandleFatalError(err)

	ch, err := rabbitMqConn.Channel()
	err_utils.HandleFatalError(err)

	defer ch.Close()

	publisher := events.NewPublisher(ch, config.SocmedExchangeName)

	// init db
	gormDB, err := db.NewGormDB(config.DatabaseUrl)
	err_utils.HandleFatalError(err)

	// init echo and API Server
	e := routes.NewEchoServer(gormDB, publisher)
	apiServer := server.NewApiServer(e, config.ServiceName, nil, gormDB, rabbitMqConn)

	go func() {
		if err := apiServer.Start(mainCtx, config.HttpPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("API server error :%v", err)
		}
	}()

	// Wait for shutdown signal
	<-mainCtx.Done()

	// Call stop to remove the signal handler
	stop()
	log.Println("Shutdown signal received. Starting graceful shutdown...")

	// Drain API Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	apiServer.Stop(shutdownCtx)

	log.Println("Application shutdown gracefully...")
}
