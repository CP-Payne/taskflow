package memory

import (
	"context"

	"github.com/CP-Payne/taskflow/user/internal/model"
	"github.com/CP-Payne/taskflow/user/internal/repository"
	"github.com/google/uuid"
)

type MemoryRepository struct {
	user map[uuid.UUID]*model.User
}

func NewInMemory() *MemoryRepository {
	return &MemoryRepository{
		user: make(map[uuid.UUID]*model.User),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, user *model.User) error {
	for _, existingUser := range r.user {
		if user.Email == existingUser.Email {
			return repository.ErrDuplicateEmail
		}

		if user.Username == existingUser.Username {
			return repository.ErrDuplicateUsername
		}
	}

	r.user[user.ID] = user
	return nil
}

func (r *MemoryRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, v := range r.user {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *MemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if v, ok := r.user[id]; ok {
		return v, nil
	}
	return nil, repository.ErrNotFound
}
