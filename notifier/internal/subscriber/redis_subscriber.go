package subscriber

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CP-Payne/taskflow/notifier/internal/service"
	"github.com/CP-Payne/taskflow/pkg/events"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisSubscriber struct {
	rdb             *redis.Client
	logger          *zap.SugaredLogger
	notificationSrv *service.NotificationService
}

func NewRedisSubscriber(rdb *redis.Client, srv *service.NotificationService, logger *zap.SugaredLogger) *RedisSubscriber {
	return &RedisSubscriber{rdb: rdb, notificationSrv: srv, logger: logger}
}

// func (s *RedisSubscriber) SubscribeAndProcess(ctx context.Context) error {
// 	pubsub := s.rdb.Subscribe(ctx, events.ChannelTaskAssigned)
// 	_, err := pubsub.Receive(ctx)
// 	if err != nil {
// 		s.logger.Fatalf("Failed to subscribe to Redis channel %s: %v", events.ChannelTaskAssigned, err)
// 		return err
// 	}
//
// 	ch := pubsub.Channel()
// 	s.logger.Infof("Subscribe to %s. Waiting for messages...", events.ChannelTaskAssigned)
//
// 	for msg := range ch {
// 		s.logger.Infof("Received message on %s", msg.Channel)
//
// 		event, err := events.UnmarshalTaskAssignedEvent([]byte(msg.Payload))
// 		if err != nil {
// 			s.logger.Errorw("Failed to unmarshal TaskAssignedEvent", "error", err, "payload", msg.Payload)
// 			continue
// 		}
//
// 		userID, err := uuid.Parse(event.UserID)
// 		if err != nil {
// 			s.logger.Warnw("Failed to parse userID", "userID", event.UserID, "error", err)
// 		}
//
// 		taskID, err := uuid.Parse(event.TaskID)
// 		if err != nil {
// 			s.logger.Warnw("Failed to parse taskID", "taskID", event.TaskID, "error", err)
// 		}
//
// 		err = s.notificationSrv.NotifyUserToCompleteTask(ctx, userID, taskID)
// 		if err != nil {
// 			s.logger.Errorw("Failed to send notification", "TaskID", event.TaskID, "RecipientID", event.UserID, "error", err)
// 		} else {
// 			s.logger.Infof("Successfully processed notification for TaskID %s to UserID %s", event.TaskID, event.UserID)
// 		}
// 	}
// 	s.logger.Info("Subscription channel closed.")
//
// 	return nil
// }
//
//
//

func (s *RedisSubscriber) SubscribeAndProcess(ctx context.Context) error {
	pubsub := s.rdb.Subscribe(ctx, events.ChannelTaskAssigned)

	defer func() {
		if err := pubsub.Close(); err != nil {
			s.logger.Errorw("Failed to close Redis pubsub", "error", err)
		} else {
			s.logger.Info("Redis pubsub closed.")
		}
	}()

	receiveCtx, receiveCancel := context.WithTimeout(ctx, 5*time.Second)
	_, err := pubsub.Receive(receiveCtx)
	receiveCancel()
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			s.logger.Infow("Subscription cancelled during initial receive", "channel", events.ChannelTaskAssigned)
			return err
		}
		s.logger.Errorw("Failed to subscribe or receive confirmation", "channel", events.ChannelTaskAssigned, "error", err)
		return fmt.Errorf("failed to subscribe to Redis channel %s: %w", events.ChannelTaskAssigned, err)
	}

	ch := pubsub.Channel()
	s.logger.Infof("Subscribed to %s. Waiting for messages or cancellation...", events.ChannelTaskAssigned)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context cancelled, stopping Redis subscriber...")
			// Return the context error to signal cancellation
			return ctx.Err() // Typically context.Canceled

		case msg, ok := <-ch:
			if !ok {
				// Channel closed
				s.logger.Info("Redis pubsub channel closed.")
				return nil
			}

			s.logger.Infof("Received message on %s", msg.Channel)

			event, err := events.UnmarshalTaskAssignedEvent([]byte(msg.Payload))
			if err != nil {
				s.logger.Errorw("Failed to unmarshal TaskAssignedEvent", "error", err, "payload", msg.Payload)
				continue // Skip this message, continue loop
			}

			userID, err := uuid.Parse(event.UserID)
			if err != nil {
				s.logger.Warnw("Failed to parse userID, skipping notification", "userID", event.UserID, "error", err)
				continue // Skip if ID is invalid
			}

			taskID, err := uuid.Parse(event.TaskID)
			if err != nil {
				s.logger.Warnw("Failed to parse taskID, skipping notification", "taskID", event.TaskID, "error", err)
				continue // Skip if ID is invalid
			}
			err = s.notificationSrv.NotifyUserToCompleteTask(ctx, userID, taskID)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					s.logger.Warnw("Notification sending cancelled", "TaskID", event.TaskID, "RecipientID", event.UserID, "error", err)
				} else {
					s.logger.Errorw("Failed to send notification", "TaskID", event.TaskID, "RecipientID", event.UserID, "error", err)
				}
			} else {
				s.logger.Infof("Successfully processed notification for TaskID %s to UserID %s", event.TaskID, event.UserID)
			}
		}
	}
}

func (s *RedisSubscriber) Close() error {
	return nil
}
