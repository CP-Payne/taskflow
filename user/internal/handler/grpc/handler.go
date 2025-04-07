package grpc

import (
	"context"
	"errors"

	"github.com/CP-Payne/taskflow/user/internal/model"
	api "github.com/CP-Payne/taskflow/user/internal/proto/api/v1"
	"github.com/CP-Payne/taskflow/user/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	api.UnimplementedUserServer
	userService *service.UserService
	logger      *zap.SugaredLogger
}

func NewUserHandler(userService *service.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *api.RegisterUserRequest) (*api.RegisterUserResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}
	user := &model.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
	}
	user.Password.Set(req.Password)
	err := h.userService.RegisterUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	token, err := h.userService.AuthenticateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrInvalidPassword):
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	return &api.RegisterUserResponse{Jwt: token}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *api.AuthenticateUserRequest) (*api.AuthenticateUserResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	user := &model.User{
		Email: req.Email,
	}
	user.Password.Set(req.Password)
	token, err := h.userService.AuthenticateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrInvalidPassword):
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}
	if token == "" {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &api.AuthenticateUserResponse{Jwt: token}, nil
}
