package main

import (
	"log"

	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/db"
	"github.com/kylehipz/socmed-microservices/users/internal/events"
	"github.com/kylehipz/socmed-microservices/users/internal/queue"
	"github.com/kylehipz/socmed-microservices/users/internal/server"
)

func main() {
	// init queue
	rabbitMqConn := queue.NewRabbitMQConnection(config.Settings.RabbitMqUrl)
	defer rabbitMqConn.Close()
	ch, err := rabbitMqConn.Channel()

	if err != nil {
		log.Fatalf("Failed to create rabbitmq channel: %v", err)
	}

	eventPublisher := events.NewPublisher(ch)

	// init db
	gormDB := db.NewGormDB(config.Settings.DatabaseUrl)

	// init echo server
	e := server.NewEchoServer(gormDB, eventPublisher)

	log.Println("Starting users service...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
