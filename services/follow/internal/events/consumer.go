// follow/internal/events/consumer.go
package events

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartUserSyncConsumer(db *gorm.DB, ch *amqp091.Channel) {
	// Dead letter Exchange and Queue
	err := ch.ExchangeDeclare(
		config.Settings.DeadLetterExchangeName,
		"fanout", // dead letters don’t need routing keys
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare DLX: %v", err)
	}

	// Declare a dead-letter queue
	dlq, err := ch.QueueDeclare(
		config.Settings.DeadLetterQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusiveFatalf
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare DLQ: %v", err)
	}

	// Bind DLQ to DLX
	if err := ch.QueueBind(dlq.Name, "", config.Settings.DeadLetterExchangeName, false, nil); err != nil {
		log.Fatalf("failed to bind DLQ: %v", err)
	}

	// Declare exchange (must be identical to publisher)
	err = ch.ExchangeDeclare(
		config.Settings.SocmedExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange": config.Settings.DeadLetterExchangeName,
		},
	)
	if err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		config.Settings.UserEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	// Bind queue to exchange with pattern
	err = ch.QueueBind(
		q.Name,
		config.Settings.AllUserEvents,       // captures user.created and user.updated
		config.Settings.SocmedExchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	ch.Qos(10, 0, false)

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var user models.User
			if err := json.Unmarshal(msg.Body, &user); err != nil {
				log.Printf("invalid payload: %v", err)
				msg.Nack(false, false)
				continue
			}

			// Upsert user into local table
			if err := db.Clauses(clause.OnConflict{ UpdateAll: true }).Create(&user).Error; err != nil {
				log.Printf("Failed to upsert user %s: %v", user.ID, err)
				requeue := !errors.Is(err, gorm.ErrInvalidData)
				msg.Nack(false, requeue)
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}()
}
