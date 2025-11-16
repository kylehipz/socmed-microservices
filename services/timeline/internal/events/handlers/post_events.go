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

type PostEventsHandler struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewPostEventsHandler(log *zap.Logger, db *gorm.DB) *PostEventsHandler {
	return &PostEventsHandler{log: log, db: db}
}

func (u *PostEventsHandler) HandlePostCreatedOrUpdatedEvent(msg amqp091.Delivery) {
	extraLogFields := []zap.Field{
		zap.String("event_name", msg.RoutingKey),
		zap.String("message_id", msg.MessageId),
	}

	var post models.Post
	if err := json.Unmarshal(msg.Body, &post); err != nil {
		u.log.Error("invalid payload", append(extraLogFields, zap.Error(err))...)
		msg.Nack(false, false)
		return
	}

	// Upsert user into local table
	if err := u.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&post).Error; err != nil {
		u.log.Error("Failed to upsert post", append(extraLogFields, zap.Error(err))...)
		requeue := !errors.Is(err, gorm.ErrInvalidData)
		msg.Nack(false, requeue)
		return
	}

	if err := msg.Ack(false); err != nil {
		u.log.Error("Failed to ack message", append(extraLogFields, zap.Error(err))...)
	}
}

func (u *PostEventsHandler) HandlePostDeletedEvent(msg amqp091.Delivery) {
	extraLogFields := []zap.Field{
		zap.String("event_name", msg.RoutingKey),
		zap.String("message_id", msg.MessageId),
	}

	var post models.Post
	if err := json.Unmarshal(msg.Body, &post); err != nil {
		u.log.Error("invalid payload", append(extraLogFields, zap.Error(err))...)
		msg.Nack(false, false)
		return
	}

	// Upsert user into local table
	if err := u.db.Delete(&post).Error; err != nil {
		u.log.Error("Failed to delete post", append(extraLogFields, zap.Error(err))...)
		requeue := !errors.Is(err, gorm.ErrInvalidData)
		msg.Nack(false, requeue)
		return
	}

	if err := msg.Ack(false); err != nil {
		u.log.Error("Failed to ack message", append(extraLogFields, zap.Error(err))...)
	}
}

