package handlers

import (
	"encoding/json"
	"errors"

	"github.com/kylehipz/socmed-microservices/timeline/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserEventsHandler struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewUserEventsHandler(log *zap.Logger, db *gorm.DB) *UserEventsHandler {
	return &UserEventsHandler{log: log, db: db}
}

func (u *UserEventsHandler) HandleUserCreatedOrUpdatedEvent(msg amqp091.Delivery) {
	extraLogFields := []zap.Field{
		zap.String("event_name", msg.RoutingKey),
		zap.String("message_id", msg.MessageId),
	}

	var user models.User
	if err := json.Unmarshal(msg.Body, &user); err != nil {
		u.log.Error("invalid payload", append(extraLogFields, zap.Error(err))...)
		msg.Nack(false, false)
		return
	}

	// Upsert user into local table
	if err := u.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error; err != nil {
		u.log.Error("Failed to upsert user", append(extraLogFields, zap.Error(err))...)
		requeue := !errors.Is(err, gorm.ErrInvalidData)
		msg.Nack(false, requeue)
		return
	}

	if err := msg.Ack(false); err != nil {
		u.log.Error("Failed to ack message", append(extraLogFields, zap.Error(err))...)
	}
}
