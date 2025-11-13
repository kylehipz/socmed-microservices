package events

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func NewRabbitMQConn(log *zap.Logger, connStr string) (*amqp091.Connection, error) {
	if conn, err := amqp091.Dial(connStr); err != nil {
		return nil, err
	} else {
		log.Info("RabbitMQ connection established successfully.")

		return conn, nil
	}

}

