package server

import (
	"net"

	grpcApi "github.com/CP-Payne/taskflow/pkg/gen/task/v1"
	grpchandler "github.com/CP-Payne/taskflow/task/internal/handler/grpc"

	"github.com/CP-Payne/taskflow/task/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartGRPCServer(addr string, taskService *service.TaskService, logger *zap.SugaredLogger) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	taskHandler := grpchandler.NewTaskHandler(taskService, logger)

	grpcApi.RegisterTaskServiceServer(grpcServer, taskHandler)

	logger.Infof("gRPC server listening on %s", addr)
	return grpcServer.Serve(lis)
}
