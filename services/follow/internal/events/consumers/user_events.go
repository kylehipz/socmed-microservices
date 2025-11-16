package consumers

import (
	"context"
	"sync"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/constants"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/events/handlers"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserEventsConsumer struct {
	log               *zap.Logger
	userEventsHandler *handlers.UserEventsHandler
	ch                *amqp091.Channel
	workerCount       int
	wg                sync.WaitGroup
}

func NewUserEventsConsumer(
	log *zap.Logger,
	ch *amqp091.Channel,
	db *gorm.DB,
	workerCount int,
) *UserEventsConsumer {
	consumerLog := log.With(zap.String("consumer_name", "follow_user_events"))
	userEventsHandler := handlers.NewUserEventsHandler(consumerLog, db)

	return &UserEventsConsumer{
		log:               consumerLog,
		ch:                ch,
		userEventsHandler: userEventsHandler,
		workerCount:       workerCount,
	}
}

func (u *UserEventsConsumer) Start(ctx context.Context) error {
	// Dead letter Exchange and Queue
	if err := u.ch.ExchangeDeclare(
		constants.DeadLetterExchangeName,
		"fanout", // dead letters don’t need routing keys
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		u.log.Error("failed to declare DLX", zap.Error(err))
		return err
	}

	// Declare a dead-letter queue
	dlq, err := u.ch.QueueDeclare(
		constants.DeadLetterQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusiveFatalf
		false,
		nil,
	)

	if err != nil {
		u.log.Error("failed to declare DLQ", zap.Error(err))
		return err
	}

	// Bind DLQ to DLX
	if err := u.ch.QueueBind(dlq.Name, "", constants.DeadLetterExchangeName, false, nil); err != nil {
		u.log.Error("failed to bind DLQ", zap.Error(err))
		return err
	}

	// Declare exchange (must be identical to publisher)
	if err = u.ch.ExchangeDeclare(
		constants.SocmedExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange": constants.DeadLetterExchangeName,
		},
	); err != nil {
		u.log.Error("failed to declare exchange", zap.Error(err))
		return err
	}

	// Declare queue
	q, err := u.ch.QueueDeclare(
		config.UserEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		u.log.Error("failed to declare queue", zap.Error(err))
		return err
	}

	// Bind queue to exchange with pattern
	if err = u.ch.QueueBind(
		q.Name,
		events.AllUserEvents, // captures user.created and user.updated
		constants.SocmedExchangeName,
		false,
		nil,
	); err != nil {
		u.log.Error("failed to bind queue", zap.Error(err))
		return err
	}

	u.ch.Qos(u.workerCount, 0, false)

	msgs, err := u.ch.Consume(config.UserEventsQueue, "", false, false, false, false, nil)
	if err != nil {
		u.log.Error("failed to register consumer", zap.Error(err))
		return err
	}

	for workerID := 0; workerID < u.workerCount; workerID++ {
		u.wg.Add(1)

		go func(workerID int) {
			defer u.wg.Done()

			workerLog := u.log.With(zap.Int("worker_id", workerID))

			for {
				select {
				case <-ctx.Done():
					workerLog.Info("Application shutdown signal received. Worker shutting down...")
					return
				case msg, ok := <-msgs:
					if !ok {
						workerLog.Info("Message channel closed. Worker shutting down...")
						return
					}

					extraLogFields := []zap.Field{
						zap.String("event_name", msg.RoutingKey),
						zap.String("message_id", msg.MessageId),
					}

					u.userEventsHandler.HandleUserCreatedOrUpdatedEvent(msg)

					workerLog.Info("User event handled successfully", extraLogFields...)
				}
			}
		}(workerID)
	}

	u.log.Info("All workers started")
	return nil
}

func (u *UserEventsConsumer) Wait(ctx context.Context) {
	u.log.Info("Attempting event consumers graceful shutdown...")
	// setup consumer shutdown context: 5s
	workersShutdownCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// wait for all workers to be done
	waitWorkersChan := make(chan struct{})

	go func() {
		defer close(waitWorkersChan)
		u.wg.Wait()
	}()

	select {
	case <-waitWorkersChan:
		u.log.Info("All workers shutdown successfully")
	case <-workersShutdownCtx.Done():
		u.log.Info("Workers shutdown timed out. Shutdown forced")
	}
}
