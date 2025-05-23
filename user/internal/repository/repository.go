package repository

import (
	"context"
	"errors"

	"github.com/CP-Payne/taskflow/user/internal/model"
	"github.com/google/uuid"
)

var (
	ErrDuplicateUsername = errors.New("username already exist")
	ErrDuplicateEmail    = errors.New("email already exist")
	ErrNotFound          = errors.New("resource not found")
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) error
}
