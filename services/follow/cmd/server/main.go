package main

import (
	"log"

	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/db"
	"github.com/kylehipz/socmed-microservices/follow/internal/events"
	"github.com/kylehipz/socmed-microservices/follow/internal/queue"
	"github.com/kylehipz/socmed-microservices/follow/internal/server"
)

func main() {
	// init queue
	rabbitMqConn := queue.NewRabbitMQConnection(config.Settings.RabbitMqUrl)
	defer rabbitMqConn.Close()
	ch, err := rabbitMqConn.Channel()

	if err != nil {
		log.Fatalf("Failed to create rabbitmq channel: %v", err)
	}

	// init db
	gormDB := db.NewGormDB(config.Settings.DatabaseUrl)

	// init echo server
	e := server.NewEchoServer(gormDB)

	// start consumer
	go events.StartUserSyncConsumer(gormDB, ch)

	log.Println("Starting follow service...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
