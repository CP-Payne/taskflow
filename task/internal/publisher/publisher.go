package publisher

import (
	"context"

	"github.com/CP-Payne/taskflow/pkg/events"
)

type Publisher interface {
	PublishTaskAssigned(ctx context.Context, event *events.TaskAssignedEvent) error
}
