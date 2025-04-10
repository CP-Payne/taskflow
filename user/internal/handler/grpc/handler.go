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
		h.logger.Warnw("RegisterUser validation failed: invalid arguments",
			"email", req.GetEmail(),
			"username", req.GetUsername(),
		)
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
			h.logger.Warnw("User registration conflict: user already exists",
				"email", req.Email,
				"username", req.Username,
				zap.Error(err),
			)
			return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
		default:
			h.logger.Errorw("Internal error during user registration",
				"email", req.Email,
				"username", req.Username,
				"assigned_userID", user.ID,
				zap.Error(err),
			)
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	token, err := h.userService.AuthenticateUser(ctx, user)
	if err != nil {
		h.logger.Errorw("Internal error during post-registration authentication",
			"userID", user.ID.String(),
			"email", user.Email,
			zap.Error(err),
		)
		switch {
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrInvalidPassword):
			return nil, status.Errorf(codes.Unauthenticated, "failed to authenticate newly registered user: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	if token == "" {
		h.logger.Errorw("Post-registration authentication returned empty token without error",
			"userID", user.ID.String(),
			"email", user.Email,
		)
		return nil, status.Errorf(codes.Internal, "internal server error: empty token received")
	}

	h.logger.Infow("User registered and authenticated successfully",
		"userID", user.ID.String(),
		"email", user.Email,
	)
	return &api.RegisterUserResponse{Jwt: token}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *api.AuthenticateUserRequest) (*api.AuthenticateUserResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		h.logger.Warnw("AuthenticateUser validation failed: invalid arguments",
			"email", req.GetEmail(), // Log email for context
		)
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
			h.logger.Warnw("User authentication failed: invalid credentials",
				"email", req.Email,
				zap.Error(err), // Include specific reason (NotFound vs InvalidPassword)
			)
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		default:
			h.logger.Errorw("Internal error during user authentication",
				"email", req.Email,
				zap.Error(err),
			)
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}
	if token == "" {
		h.logger.Errorw("Authentication service returned empty token without error",
			"email", req.Email,
		)
		return nil, status.Errorf(codes.Internal, "internal server error: empty token received")
	}

	h.logger.Infow("User authenticated successfully",
		"email", req.Email,
	)
	return &api.AuthenticateUserResponse{Jwt: token}, nil
}

func (h *UserHandler) GetByID(ctx context.Context, req *api.GetByIDRequest) (*api.GetByIDResponse, error) {
	if req == nil || req.GetUserId() == "" {
		h.logger.Warnw("GetByID validation failed: invalid arguments",
			"userID", req.GetUserId(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	// Parse UUID
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		h.logger.Warnw("GetByID invalid userID",
			"userID", req.GetUserId(),
			"error", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid userID")
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		h.logger.Errorw("GetByID internal error",
			"userID", userID.String(),
			"error", err,
		)
		return nil, status.Errorf(codes.Internal, "internal server error")

	}

	h.logger.Infow("User retrieved successfully", "userID", userID)
	return &api.GetByIDResponse{
		UserId:   user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}, nil
}
