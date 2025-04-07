package main

import (
	"github.com/CP-Payne/taskflow/user/config"
	"github.com/CP-Payne/taskflow/user/config/authkeys"
	"github.com/CP-Payne/taskflow/user/internal/auth"
	"github.com/CP-Payne/taskflow/user/internal/repository/memory"
	"github.com/CP-Payne/taskflow/user/internal/server"
	"github.com/CP-Payne/taskflow/user/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New("user/config/.env")

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	authKeys := authkeys.NewAuthKeys()
	authKeys.Load(cfg.KeyPath)

	authenticator := auth.NewJWTAuthenticator(authKeys)

	repo := memory.NewInMemory()
	srv := service.New(repo, authenticator, logger)
	server.StartGRPCServer("localhost:3033", srv, logger)
}
