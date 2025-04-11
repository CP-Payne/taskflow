package notification

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, recipient string, message string) error
}
