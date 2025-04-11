package subscriber

import (
	"context"

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

		userID, err := uuid.Parse(event.UserID)
		if err != nil {
			s.logger.Warnw("Failed to parse userID", "userID", event.UserID, "error", err)
		}

		taskID, err := uuid.Parse(event.TaskID)
		if err != nil {
			s.logger.Warnw("Failed to parse taskID", "taskID", event.TaskID, "error", err)
		}

		err = s.notificationSrv.NotifyUserToCompleteTask(ctx, userID, taskID)
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
