package memory_test

import (
	"context"
	"testing"

	"github.com/CP-Payne/taskflow/user/internal/model"
	"github.com/CP-Payne/taskflow/user/internal/repository/memory"
	"github.com/google/uuid"
)

func TestMemoryRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo *memory.MemoryRepository) *model.User
		expectErr error
	}{
		{
			name: "Successfully create a user",
			setup: func(repo *memory.MemoryRepository) *model.User {
				return &model.User{
					ID:       uuid.New(),
					Email:    "test@example.com",
					Username: "testuser",
				}
			},
			expectErr: nil,
		},
		{
			name: "Fail to create user with duplicate email",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "duplicate@example.com",
					Username: "user1",
				}
				_ = repo.Create(context.Background(), user) // Pre-add user
				return &model.User{
					ID:       uuid.New(),
					Email:    "duplicate@example.com", // Same email, different ID
					Username: "user2",
				}
			},
			expectErr: memory.ErrDuplicateEmail,
		},
		{
			name: "Fail to create user with duplicate username",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "user1@example.com",
					Username: "duplicateUser",
				}
				_ = repo.Create(context.Background(), user) // Pre-add user
				return &model.User{
					ID:       uuid.New(),
					Email:    "user2@example.com",
					Username: "duplicateUser", // Same username, different ID
				}
			},
			expectErr: memory.ErrDuplicateUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := memory.NewInMemory()
			ctx := context.Background()
			user := tt.setup(repo)

			// Act
			err := repo.Create(ctx, user)

			// Assert
			if err != tt.expectErr {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}

			// Verify user creation only if no error is expected
			if err == nil {
				storedUser, err := repo.GetByID(ctx, user.ID)
				if err != nil {
					t.Errorf("GetByID() failed: %v", err)
				}
				if storedUser.Email != user.Email || storedUser.Username != user.Username {
					t.Errorf("Stored user does not match the original: got %+v, want %+v", storedUser, user)
				}
			}
		})
	}
}

func TestMemoryRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo *memory.MemoryRepository) *model.User
		expectErr error
	}{
		{
			name: "Successfully get user by ID",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "test@example.com",
					Username: "testuser",
				}

				_ = repo.Create(context.Background(), user)

				return user
			},
			expectErr: nil,
		},
		{
			name: "Fail to get user with NotFound",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "test2@example.com",
					Username: "testuser2",
				}

				_ = repo.Create(context.Background(), user)

				return &model.User{
					ID:       uuid.New(),
					Email:    "someemail@example.com",
					Username: "someUsername",
				}
			},
			expectErr: memory.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewInMemory()
			ctx := context.Background()
			user := tt.setup(repo)

			userRetrieved, err := repo.GetByID(ctx, user.ID)

			if err != tt.expectErr {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}

			if err == nil {
				if userRetrieved.Email != user.Email || userRetrieved.Username != user.Username {
					t.Errorf("Retrieved user does not match the original: got %+v, want %+v", userRetrieved, user)
				}
			}
		})
	}
}

func TestMemoryRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo *memory.MemoryRepository) *model.User
		expectErr error
	}{
		{
			name: "Successfully get user by email",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "test@example.com",
					Username: "testuser",
				}

				_ = repo.Create(context.Background(), user)

				return user
			},
			expectErr: nil,
		},
		{
			name: "Fail to get user with NotFound",
			setup: func(repo *memory.MemoryRepository) *model.User {
				user := &model.User{
					ID:       uuid.New(),
					Email:    "test2@example.com",
					Username: "testuser2",
				}

				_ = repo.Create(context.Background(), user)

				return &model.User{
					ID:       uuid.New(),
					Email:    "someemail@example.com",
					Username: "someUsername",
				}
			},
			expectErr: memory.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewInMemory()
			ctx := context.Background()
			user := tt.setup(repo)

			userRetrieved, err := repo.GetByEmail(ctx, user.Email)

			if err != tt.expectErr {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}

			if err == nil {
				if userRetrieved.ID != user.ID || userRetrieved.Username != user.Username {
					t.Errorf("Retrieved user does not match the original: got %+v, want %+v", userRetrieved, user)
				}
			}
		})
	}
}
