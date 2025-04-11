package user

import (
	"context"

	"github.com/CP-Payne/taskflow/notifier/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	logger *zap.SugaredLogger
}

func NewGateway(logger *zap.SugaredLogger) *Gateway {
	return &Gateway{
		logger: logger,
	}
}

func (g *Gateway) GetUserDetails(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	conn, err := grpc.NewClient("localhost:3032", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.logger.Error("Failed to open connection to user service")
		return nil, err
	}

	defer conn.Close()
}
