package repository

import (
	"context"

	"github.com/CP-Payne/taskflow/user/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) error
}
