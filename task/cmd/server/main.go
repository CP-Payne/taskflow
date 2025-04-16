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

	"github.com/CP-Payne/taskflow/pkg/discovery"
	"github.com/CP-Payne/taskflow/pkg/discovery/consul"
	grpcApi "github.com/CP-Payne/taskflow/pkg/gen/task/v1"
	grpchandler "github.com/CP-Payne/taskflow/task/internal/handler/grpc"
	"github.com/CP-Payne/taskflow/task/internal/publisher"
	"github.com/CP-Payne/taskflow/task/internal/repository/memory"
	"github.com/CP-Payne/taskflow/task/internal/service"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	serviceName     = "task"
	shutdownTimeout = 15 * time.Second
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	var port int
	flag.IntVar(&port, "port", 9002, "API handler port")
	flag.Parse()
	logger.Infof("Starting the task service on port %d", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load global config
	if err := godotenv.Load("./config/.env"); err != nil {
		logger.Fatalw("failed to load global config", "error", err)
	}

	// TODO: Add registry addr to global config
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		logger.Fatalw("failed to create service discovery registry", "error", err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	serviceAddr := fmt.Sprintf("localhost:%d", port)

	redisAddr := os.Getenv("REDIS_NOTIFIER_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatalw("Failed to connect to Redis", "error", err)
	}
	defer rdb.Close()
	logger.Info("Connected to Redis")

	redisPublisher := publisher.NewRedisPublisher(rdb, logger)

	// TODO: Create config to pass to layers
	// TODO: Define Handler in main, instead of StartGRPCServer
	repo := memory.NewInMemory()
	srv := service.New(repo, logger, redisPublisher)

	grpcServer := grpc.NewServer()
	taskHandler := grpchandler.NewTaskHandler(srv, logger)
	grpcApi.RegisterTaskServiceServer(grpcServer, taskHandler)

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
