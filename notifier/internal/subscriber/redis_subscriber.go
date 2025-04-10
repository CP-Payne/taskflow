package subscriber

import (
	"context"
	"fmt"

	"github.com/CP-Payne/taskflow/notifier/internal/notification"
	"github.com/CP-Payne/taskflow/pkg/events"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisSubscriber struct {
	rdb    *redis.Client
	sender notification.Sender
	logger *zap.SugaredLogger
}

func NewRedisSubscriber(rdb *redis.Client, sender notification.Sender, logger *zap.SugaredLogger) *RedisSubscriber {
	return &RedisSubscriber{rdb: rdb, sender: sender, logger: logger}
}

func (s *RedisSubscriber) SubscribeAndProcess(ctx context.Context) {
	pubsub := s.rdb.Subscribe(ctx, events.ChannelTaskAssigned)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		s.logger.Fatalf("Failed to subscribe to Redis channel %s: %v", events.ChannelTaskAssigned, err)
		return
	}

	ch := pubsub.Channel()
	s.logger.Infof("Subscribe to %s. Waiting for messages...", events.ChannelTaskAssigned)

	for msg := range ch {
		s.logger.Infof("Received message on %s", msg.Channel)

		event, err := events.UnmarshalTaskAssignedEvent([]byte(msg.Payload))
		if err != nil {
			s.logger.Errorw("Failed to unmarshal TaskAssignedEvent", "error", err, "payload", msg.Payload)
			continue
		}

		notificationMsg := fmt.Sprintf("Task '%s' has been assigned to you.", event.TaskID)

		err = s.sender.Send(ctx, event.UserID, notificationMsg)
		if err != nil {
			s.logger.Errorw("Failed to send notification", "TaskID", event.TaskID, "RecipientID", event.UserID, "error", err)
		} else {
			s.logger.Infof("Successfully processed notification for TaskID %s to UserID %s", event.TaskID, event.UserID)
		}
	}
	s.logger.Info("Subscription channel closed.")
}

func (s *RedisSubscriber) Close() error {
	return nil
}
