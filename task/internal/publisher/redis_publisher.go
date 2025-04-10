package publisher

import (
	"context"

	"github.com/CP-Payne/taskflow/pkg/events"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisPublisher struct {
	rdb    *redis.Client
	logger *zap.SugaredLogger
}

func NewRedisPublisher(rdb *redis.Client, logger *zap.SugaredLogger) *RedisPublisher {
	return &RedisPublisher{rdb: rdb, logger: logger}
}

func (p *RedisPublisher) PublishTaskAssigned(ctx context.Context, event *events.TaskAssignedEvent) error {
	payload, err := event.Marshal()
	if err != nil {
		p.logger.Errorw("Failed to marshal TaskAssignedEvent", "error", zap.Error(err))
		// log.Printf("ERROR: Failed to marshal TaskAssignedEvent: %v", err)
		return err
	}

	err = p.rdb.Publish(ctx, events.ChannelTaskAssigned, payload).Err()
	if err != nil {
		p.logger.Errorw("Failed to publish TaskAssignedEvent", "error", err, "channel", events.ChannelTaskAssigned)
		// log.Printf("ERROR: Failed to publish TaskAssignedEvent to %s: %v", events.ChannelTaskAssigned, err)
		return err
	}

	p.logger.Infow("Published TaskAssignedEvent", "channel", events.ChannelTaskAssigned, "event", event)
	// log.Printf("Published TaskAssignedEvent to %s: %+v", events.ChannelTaskAssigned, event)
	return nil
}
