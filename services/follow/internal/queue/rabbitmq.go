package queue

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(connectionString string) *amqp091.Connection {
	conn, err := amqp091.Dial(connectionString)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("RabbitMQ connection established successfully.")

	return conn
}
