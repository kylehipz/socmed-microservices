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

type FollowEventsHandler struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewFollowEventsHandler(log *zap.Logger, db *gorm.DB) *FollowEventsHandler {
	return &FollowEventsHandler{log: log, db: db}
}

func (u *FollowEventsHandler) HandleUserFollowedEvent(msg amqp091.Delivery) {
	extraLogFields := []zap.Field{
		zap.String("event_name", msg.RoutingKey),
		zap.String("message_id", msg.MessageId),
	}

	var follow models.Follow
	if err := json.Unmarshal(msg.Body, &follow); err != nil {
		u.log.Error("invalid payload", append(extraLogFields, zap.Error(err))...)
		msg.Nack(false, false)
		return
	}

	// Upsert user into local table
	if err := u.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&follow).Error; err != nil {
		u.log.Error("Failed to upsert follow", append(extraLogFields, zap.Error(err))...)
		requeue := !errors.Is(err, gorm.ErrInvalidData)
		msg.Nack(false, requeue)
		return
	}

	if err := msg.Ack(false); err != nil {
		u.log.Error("Failed to ack message", append(extraLogFields, zap.Error(err))...)
	}
}

func (u *FollowEventsHandler) HandleUserUnfollowedEvent(msg amqp091.Delivery) {
	extraLogFields := []zap.Field{
		zap.String("event_name", msg.RoutingKey),
		zap.String("message_id", msg.MessageId),
	}

	var follow models.Follow
	if err := json.Unmarshal(msg.Body, &follow); err != nil {
		u.log.Error("invalid payload", append(extraLogFields, zap.Error(err))...)
		msg.Nack(false, false)
		return
	}

	// Upsert user into local table
	if err := u.db.Delete(&follow).Error; err != nil {
		u.log.Error("Failed to delete follow", append(extraLogFields, zap.Error(err))...)
		requeue := !errors.Is(err, gorm.ErrInvalidData)
		msg.Nack(false, requeue)
		return
	}

	if err := msg.Ack(false); err != nil {
		u.log.Error("Failed to ack message", append(extraLogFields, zap.Error(err))...)
	}
}
