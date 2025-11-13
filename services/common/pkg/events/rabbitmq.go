package events

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConn(connStr string) (*amqp091.Connection, error) {
	if conn, err := amqp091.Dial(connStr); err != nil {
		return nil, err
	} else {
		log.Println("RabbitMQ connection established successfully.")

		return conn, nil
	}

}

