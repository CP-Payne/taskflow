package subscriber

import "context"

type Subscriber interface {
	SubscribeAndProcess(ctx context.Context) error
	// Close cleans up any resources if needed
	Close() error
}
