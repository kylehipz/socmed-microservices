package consumers

import (
	"context"
	"sync"
	"time"

	"github.com/kylehipz/socmed-microservices/common/pkg/constants"
	"github.com/kylehipz/socmed-microservices/common/pkg/events"
	"github.com/kylehipz/socmed-microservices/timeline/config"
	"github.com/kylehipz/socmed-microservices/timeline/internal/events/handlers"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FollowEventsConsumer struct {
	log               *zap.Logger
	followEventsHandler *handlers.FollowEventsHandler
	ch                *amqp091.Channel
	workerCount       int
	wg                sync.WaitGroup
}

func NewFollowEventsConsumer(
	log *zap.Logger,
	ch *amqp091.Channel,
	db *gorm.DB,
	workerCount int,
) *FollowEventsConsumer {
	consumerLog := log.With(zap.String("consumer_name", "timeline_follow_events"))
	followEventsHandler := handlers.NewFollowEventsHandler(consumerLog, db)

	return &FollowEventsConsumer{
		log:               consumerLog,
		ch:                ch,
		followEventsHandler: followEventsHandler,
		workerCount:       workerCount,
	}
}

func (f *FollowEventsConsumer) Start(ctx context.Context) error {
	// Dead letter Exchange and Queue
	if err := f.ch.ExchangeDeclare(
		constants.DeadLetterExchangeName,
		"fanout", // dead letters don’t need routing keys
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		f.log.Error("failed to declare DLX", zap.Error(err))
		return err
	}

	// Declare a dead-letter queue
	dlq, err := f.ch.QueueDeclare(
		constants.DeadLetterQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusiveFatalf
		false,
		nil,
	)

	if err != nil {
		f.log.Error("failed to declare DLQ", zap.Error(err))
		return err
	}

	// Bind DLQ to DLX
	if err := f.ch.QueueBind(dlq.Name, "", constants.DeadLetterExchangeName, false, nil); err != nil {
		f.log.Error("failed to bind DLQ", zap.Error(err))
		return err
	}

	// Declare exchange (must be identical to publisher)
	if err = f.ch.ExchangeDeclare(
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
		f.log.Error("failed to declare exchange", zap.Error(err))
		return err
	}

	// Declare queue
	q, err := f.ch.QueueDeclare(
		config.FollowEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		f.log.Error("failed to declare queue", zap.Error(err))
		return err
	}

	// Bind queue to exchange with pattern
	if err = f.ch.QueueBind(
		q.Name,
		events.AllFollowEvents, // captures user.created and user.updated
		constants.SocmedExchangeName,
		false,
		nil,
	); err != nil {
		f.log.Error("failed to bind queue", zap.Error(err))
		return err
	}

	f.ch.Qos(f.workerCount, 0, false)

	msgs, err := f.ch.Consume(config.FollowEventsQueue, "", false, false, false, false, nil)
	if err != nil {
		f.log.Error("failed to register consumer", zap.Error(err))
		return err
	}

	for workerID := 0; workerID < f.workerCount; workerID++ {
		f.wg.Add(1)

		go func(workerID int) {
			defer f.wg.Done()

			workerLog := f.log.With(zap.Int("worker_id", workerID))

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

				switch msg.RoutingKey {
					case events.UserFollowed:
						f.followEventsHandler.HandleUserFollowedEvent(msg)
					case events.UserUnfollowed:
						f.followEventsHandler.HandleUserUnfollowedEvent(msg)
					}

					workerLog.Info("Follow event handled successfully", extraLogFields...)
				}
			}
		}(workerID)
	}

	f.log.Info("All workers started")
	return nil
}

func (f *FollowEventsConsumer) Wait(ctx context.Context) {
	f.log.Info("Attempting event consumers graceful shutdown...")
	// setup consumer shutdown context: 5s
	workersShutdownCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// wait for all workers to be done
	waitWorkersChan := make(chan struct{})

	go func() {
		defer close(waitWorkersChan)
		f.wg.Wait()
	}()

	select {
	case <-waitWorkersChan:
		f.log.Info("All workers shutdown successfully")
	case <-workersShutdownCtx.Done():
		f.log.Info("Workers shutdown timed out. Shutdown forced")
	}
}
