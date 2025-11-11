// follow/internal/events/consumer.go
package events

import (
	"encoding/json"
	"log"

	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartUserSyncConsumer(db *gorm.DB, ch *amqp091.Channel) {
	// Declare exchange (must be identical to publisher)
	err := ch.ExchangeDeclare(
		config.Settings.SocmedExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
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

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var user models.User
			if err := json.Unmarshal(msg.Body, &user); err != nil {
				log.Printf("invalid payload: %v", err)
				msg.Nack(false, true)
				continue
			}

			// Upsert user into local table
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&user)

			msg.Ack(false)
		}
	}()
}
