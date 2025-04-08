package main

import (
	"github.com/CP-Payne/taskflow/task/internal/repository/memory"
	"github.com/CP-Payne/taskflow/task/internal/server"
	"github.com/CP-Payne/taskflow/task/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// TODO: Create config to pass to layers
	// TODO: Define Handler in main, instead of StartGRPCServer
	repo := memory.NewInMemory()
	srv := service.New(repo, logger)
	server.StartGRPCServer("localhost:3034", srv, logger)
}
