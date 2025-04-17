package user

import (
	"context"
	"math/rand"

	"github.com/CP-Payne/taskflow/notifier/internal/model"
	"github.com/CP-Payne/taskflow/pkg/discovery"
	gen "github.com/CP-Payne/taskflow/pkg/gen/user/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	logger   *zap.SugaredLogger
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry, logger *zap.SugaredLogger) *Gateway {
	return &Gateway{
		logger:   logger,
		registry: registry,
	}
}

func (g *Gateway) GetUserDetails(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	// Get service address from registry
	addrs, err := g.registry.ServiceAddresses(ctx, "user")
	if err != nil {
		return nil, err
	}

	// g.logger.Infof("ADDRESSES RETRIEVED: %v", addrs)

	target := addrs[rand.Intn(len(addrs))]

	g.logger.Infow("Calling user service", "address", target)

	// conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.logger.Errorw("Failed to open connection to user service", "address", target)
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserClient(conn)
	res, err := client.GetByID(ctx, &gen.GetByIDRequest{
		UserId: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		UserID:   userID,
		Username: res.GetUsername(),
		Email:    res.GetEmail(),
	}, nil
}
