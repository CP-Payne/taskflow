package subscriber

import (
	"context"
	"fmt"
	"log"

	"github.com/CP-Payne/taskflow/notifier/internal/notification"
	"github.com/CP-Payne/taskflow/pkg/events"
	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	rdb    *redis.Client
	sender notification.Sender
}

func NewRedisSubscriber(rdb *redis.Client, sender notification.Sender) *RedisSubscriber {
	return &RedisSubscriber{rdb: rdb, sender: sender}
}

func (s *RedisSubscriber) SubscribeAndProcess(ctx context.Context) {
	pubsub := s.rdb.Subscribe(ctx, events.ChannelTaskAssigned)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Failed to subscribe to Redis channel %s: %v", events.ChannelTaskAssigned, err)
		return
	}

	ch := pubsub.Channel()
	log.Printf("Subscribe to %s. Waiting for messages...\n", events.ChannelTaskAssigned)

	for msg := range ch {
		log.Printf("Received message on %s\n", msg.Channel)
		event, err := events.UnmarshalTaskAssignedEvent([]byte(msg.Payload))
		if err != nil {
			log.Printf("ERROR: Failed to unmarshal TaskAssignedEvent: %v. Payload: %s", err, msg.Payload)
			continue
		}

		notificationMsg := fmt.Sprintf("Task '%s' has been assigned to you.", event.TaskID)

		err = s.sender.Send(ctx, event.UserID, notificationMsg)
		if err != nil {
			log.Printf("ERROR: Failed to send notification for TaskID %s to UserID %s: %v", event.TaskID, event.UserID, err)
		} else {
			log.Printf("Successfully processed notification for TaskID %s to UserID %s", event.TaskID, event.UserID)
		}
	}
	log.Println("Subscription channel closed.")
}

func (s *RedisSubscriber) Close() error {
	return nil
}
