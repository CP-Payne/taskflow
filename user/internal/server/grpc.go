package server

import (
	"net"

	grpcApi "github.com/CP-Payne/taskflow/pkg/gen/user/v1"
	grpchandler "github.com/CP-Payne/taskflow/user/internal/handler/grpc"
	"github.com/CP-Payne/taskflow/user/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartGRPCServer(addr string, userService *service.UserService, logger *zap.SugaredLogger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	userHandler := grpchandler.NewUserHandler(userService, logger)

	grpcApi.RegisterUserServer(grpcServer, userHandler)

	logger.Infof("gRPC server listening on %s", addr)
	return grpcServer.Serve(lis)
}
