package events

import (
	"encoding/json"
	"log"

	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel
}

func NewPublisher(ch *amqp091.Channel) *Publisher {
	// Declare exchange on startup
	if err := ch.ExchangeDeclare(
		config.Settings.SocmedExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	return &Publisher{ch: ch}
}

func (p *Publisher) PublishUserEvent(eventType string, payload any) error {
	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return p.ch.Publish(
		config.Settings.SocmedExchangeName,
		eventType,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
