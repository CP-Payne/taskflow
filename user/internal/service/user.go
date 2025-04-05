package service

import (
	"context"
	"errors"
	"time"

	"github.com/CP-Payne/taskflow/user/internal/auth"
	"github.com/CP-Payne/taskflow/user/internal/model"
	"github.com/CP-Payne/taskflow/user/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrNotFound        = errors.New("resource not found")
	ErrInternal        = errors.New("internal server error")
	ErrInvalidPassword = errors.New("invalid password hash")
)

type Service struct {
	repo          repository.UserRepository
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

func New(repo repository.UserRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user *model.User) error {
	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AuthenticateUser(ctx context.Context, user *model.User) (string, error) {
	userDB, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", ErrNotFound
		default:
			return "", ErrInternal
		}
	}

	if err := userDB.Password.Compare(*user.Password.Text); err != nil {
		return "", ErrInvalidPassword
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": "taskflow-user-service",
		"aud": "taskflow-api",
	}
	token, err := s.authenticator.GenerateToken(claims)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}
