package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CP-Payne/taskflow/notifier/internal/gateway/user"
	"github.com/CP-Payne/taskflow/notifier/internal/notification"
	"github.com/CP-Payne/taskflow/notifier/internal/service"
	"github.com/CP-Payne/taskflow/notifier/internal/subscriber"
	"github.com/CP-Payne/taskflow/pkg/discovery"
	"github.com/CP-Payne/taskflow/pkg/discovery/consul"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	serviceName         = "notifier"
	shutdownTimeout     = 15 * time.Second
	healthCheckInterval = 5 * time.Second
	healthCheckTimeout  = 2 * time.Second
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	var port int
	flag.IntVar(&port, "port", 9003, "Service identification port (not listening)")
	flag.Parse()
	logger.Infof("Starting the notifier service on port %d", port)

	ctx, cancel := context.WithCancel(context.Background())

	if err := godotenv.Load("./config/.env"); err != nil {
		logger.Fatalw("Failed to load global config", "error", err)
	}

	if err := godotenv.Load("./notifier/internal/config/.env"); err != nil {
		logger.Fatalw("Failed to load notifer service config", "error", err)
	}

	// TODO: Add registry addr to global config
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		logger.Fatalw("Failed to create service discovery registry", "error", err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	serviceAddr := fmt.Sprintf("localhost:%d", port)

	// TODO: Add to .env and load from config
	redisAddr := os.Getenv("REDIS_NOTIFIER_ADDR")
	if redisAddr == "" {
		logger.Fatalw("Redis address not configured (REDIS_NOTIFIER_ADDR)")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)

	_, err = rdb.Ping(pingCtx).Result()
	// Release context resources
	pingCancel()
	if err != nil {
		logger.Fatalw("Failed to connect to Redis", "error", err, "address", redisAddr)
	}
	defer rdb.Close()
	logger.Infow("Connected to Redis", "address", redisAddr)

	// Gmail sender dependencies
	gmailSource := os.Getenv("GMAIL_SOURCE")
	gmailAppPass := os.Getenv("GMAIL_APP_PASSWORD")
	if gmailSource == "" || gmailAppPass == "" {
		logger.Warnw("Gmail credentials not fully configured")
	}

	userGtw := user.NewGateway(logger)
	notificationSender := notification.NewEmailSender(gmailSource, gmailAppPass, logger)
	notificationSrv := service.NewNotificationService(userGtw, notificationSender)
	redisSubscriber := subscriber.NewRedisSubscriber(rdb, notificationSrv, logger)

	// Register to consul
	registerCtx, registerCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer registerCancel()
	if err := registry.Register(registerCtx, instanceID, serviceName, serviceAddr); err != nil {
		logger.Fatalw("failed registering service to consul", "error", err)
	}
	logger.Infow("Service registered successfully", "id", instanceID, "name", serviceName, "address", serviceAddr)

	go func(ctx context.Context) {
		ticker := time.NewTicker(healthCheckInterval)
		defer ticker.Stop()
		logger.Info("Starting health reporter...")
		for {
			select {
			case <-ticker.C:
				reportCtx, reportCancel := context.WithTimeout(context.Background(), healthCheckTimeout)
				if err := registry.ReportHealthState(reportCtx, instanceID, serviceName); err != nil {
					logger.Warnw("failed to report healthy state", "error", err, "instanceID", instanceID)
				}
				reportCancel()
			case <-ctx.Done():
				logger.Info("Stopping health reporter due to context cancellation.")
				return
			}
		}
	}(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Starting Redis subscriber...")
		if err = redisSubscriber.SubscribeAndProcess(ctx); err != nil && !errors.Is(err, context.Canceled) {
			// Log unexpected errors. context.Canceled is expected on graceful shutdown.
			logger.Errorw("Redis subscriber stopped unexpectedly", "error", err)
		} else {
			logger.Info("Redis subscriber processing stopped.")
		}
	}()

	// --- Signal Handling ---
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// --- Wait for shutdown signal ---
	sig := <-shutdownChan
	logger.Infow("Received shutdown signal", "signal", sig)

	// --- Initiate Graceful Shutdown ---
	logger.Info("Initiating graceful shutdown...")

	// 1. Create a context with timeout for shutdown cleanup
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// 2. Deregister from Consul *first*
	logger.Info("Deregistering service from Consul...")
	if err := registry.Deregister(shutdownCtx, instanceID, serviceName); err != nil {
		logger.Errorw("failed to deregister service during shutdown", "error", err, "instanceID", instanceID)
		// Log error but continue shutdown
	} else {
		logger.Info("Service deregistered successfully.")
	}

	// 3. Signal background tasks (subscriber, health check) to stop
	logger.Info("Signaling background tasks to stop...")
	cancel() // Cancel the main context

	// 4. Wait for the subscriber goroutine (and potentially others) to finish
	logger.Info("Waiting for background tasks to finish...")
	wg.Wait() // Block until wg.Done() is called for all Add(1) calls
	logger.Info("Background tasks finished.")

	// 5. Close Redis connection explicitly (defer also works, but this is cleaner timing)
	logger.Info("Closing Redis connection...")
	if err := rdb.Close(); err != nil {
		logger.Errorw("Error closing Redis connection", "error", err)
	}

	logger.Info("Shutdown complete.")
}
