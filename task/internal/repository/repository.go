package repository

import (
	"context"
	"errors"

	"github.com/CP-Payne/taskflow/task/internal/model"
	"github.com/google/uuid"
)

// ErrNotFound is returned when a resource is not found, such as a UserID
var ErrNotFound = errors.New("resource not found")

type TaskRepository interface {
	GetByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	ListByAssignedUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error)
	ListUnassigned(ctx context.Context) ([]model.Task, error)
	// TODO: Update Task Status
	// TODO: Assign Unassigned task to user
}
