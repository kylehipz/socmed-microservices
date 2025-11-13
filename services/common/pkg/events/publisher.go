package events

import (
	"encoding/json"
	"log"
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel
	exchangeName string
}

func NewPublisher(ch *amqp091.Channel, exchangeName string) *Publisher {
	// Declare exchange on startup
	if err := ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	return &Publisher{ch: ch, exchangeName: exchangeName}
}

func (p *Publisher) PublishEvent(eventType string, payload any) error {
	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return p.ch.Publish(
		p.exchangeName,
		eventType,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
