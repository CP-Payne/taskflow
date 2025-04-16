package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CP-Payne/taskflow/pkg/authkeys"
	"github.com/CP-Payne/taskflow/pkg/discovery"
	"github.com/CP-Payne/taskflow/pkg/discovery/consul"
	grpcApi "github.com/CP-Payne/taskflow/pkg/gen/user/v1"
	"github.com/CP-Payne/taskflow/pkg/secrets"
	"github.com/CP-Payne/taskflow/user/config"
	"github.com/CP-Payne/taskflow/user/internal/auth"
	grpchandler "github.com/CP-Payne/taskflow/user/internal/handler/grpc"
	"github.com/CP-Payne/taskflow/user/internal/repository/memory"
	"github.com/CP-Payne/taskflow/user/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	serviceName     = "user"
	shutdownTimeout = 15 * time.Second
)

func main() {
	_ = config.New("user/.env")

	logger := zap.Must(zap.NewDevelopment()).Sugar()
	defer logger.Sync()

	var port int
	flag.IntVar(&port, "port", 9001, "API handler port")
	flag.Parse()
	logger.Infof("Starting the user service on port %d", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register service on Consul
	// TODO: Add registry addr to global config
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		logger.Fatalw("failed to create service discovery registry", "error", err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	serviceAddr := fmt.Sprintf("localhost:%d", port)

	vaultAddr := os.Getenv("VAULT_ADDR")
	roleID := os.Getenv("APPROLE_ROLE_ID")
	secretID := os.Getenv("APPROLE_SECRET_ID")

	secretsManager, err := secrets.NewVaultSecretManager(vaultAddr, roleID, secretID)
	if err != nil {
		logger.Fatalw("failed to create secrets manager", "error", err)
	}

	keyPath := os.Getenv("VAULT_KEY_PATH")
	keyName := os.Getenv("VAULT_KEY_NAME")

	authKeys := authkeys.NewAuthKeys()
	err = authKeys.LoadFromSecretsManager(secretsManager, keyPath, keyName, authkeys.Private)
	if err != nil {
		logger.Fatalw("failed to fetch auth key", "error", err)
	}

	authenticator := auth.NewJWTAuthenticator(authKeys)

	repo := memory.NewInMemory()
	srv := service.New(repo, authenticator, logger)

	grpcServer := grpc.NewServer()
	userHandler := grpchandler.NewUserHandler(srv, logger)
	grpcApi.RegisterUserServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalw("failed to listen", "port", port, "error", err)
	}
	logger.Infof("gRPC server listening on %s", lis.Addr().String())

	// Register service to consul
	registerCtx, registerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer registerCancel()
	if err := registry.Register(registerCtx, instanceID, serviceName, serviceAddr); err != nil {
		logger.Fatalw("failed registering service to consul", "error", err)
	}
	logger.Infow("Service registered successfully", "id", instanceID, "name", serviceName, "address", serviceAddr)

	// HealthCheck

	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		logger.Info("Starting health reporter")
		for {
			select {
			case <-ticker.C:
				reportCtx, reportCancel := context.WithTimeout(context.Background(), 2*time.Second)
				if err := registry.ReportHealthState(reportCtx, instanceID, serviceName); err != nil {
					logger.Warnw("failed to report healthy state", "error", err, "instanceID", instanceID)
				}
				reportCancel()
			case <-ctx.Done():
				logger.Info("Stopping health reported due to context cancellation")
				return
			}
		}
	}(ctx)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting gRPC server...")
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			logger.Errorf("gRPC server failed: %v", err)
			errChan <- err
		} else {
			logger.Info("gRPC server stopped gracefully")
			close(errChan)
		}
	}()

	select {
	case sig := <-shutdownChan:
		logger.Infow("Received shutdown signal", "signal", sig)
	case err := <-errChan:
		logger.Errorw("gRPC server failed to start or run", "error", err)
	}

	logger.Info("Initiating graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	logger.Info("Deregistering service from Consul...")
	if err := registry.Deregister(shutdownCtx, instanceID, serviceName); err != nil {
		logger.Errorw("failed to deregister service during shutdown", "error", err, "instanceID", instanceID)
	} else {
		logger.Info("Service deregistered successfully")
	}

	logger.Info("Stopping background tasks (like health check)...")
	cancel()

	logger.Info("Stopping gRPC server gracefully...")
	grpcServer.GracefulStop()
	logger.Info("gRPC server has been stopped.")

	logger.Info("Flushing logs...")
	_ = logger.Sync()

	logger.Info("Shutdown complete.")
}
