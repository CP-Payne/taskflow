package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CP-Payne/taskflow/notifier/internal/gateway/user"
	"github.com/CP-Payne/taskflow/notifier/internal/notification"
	"github.com/CP-Payne/taskflow/notifier/internal/service"
	"github.com/CP-Payne/taskflow/notifier/internal/subscriber"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// TODO: Implement Gateway to retrieve user email by ID from User Service
// TODO: Implement Email sender

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load("./config/.env"); err != nil {
		logger.Panic("Failed to load global config")
	}

	if err := godotenv.Load("./notifier/internal/config/.env"); err != nil {
		logger.Panic("Failed to load notifer service config")
	}

	// TODO: Add to .env and load from config
	redisAddr := os.Getenv("REDIS_NOTIFIER_ADDR")
	// redisAddr := "localhost:6379"
	//
	gmailSource := os.Getenv("GMAIL_SOURCE")
	gmailAppPass := os.Getenv("GMAIL_APP_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Errorw("Failed to connect to Redis", "error", err)
	}
	defer rdb.Close()
	logger.Info("Connected to Redis")

	userGtw := user.NewGateway(logger)
	notificationSender := notification.NewEmailSender(gmailSource, gmailAppPass, logger)
	notificationSrv := service.NewNotificationService(userGtw, notificationSender)
	redisSubscriber := subscriber.NewRedisSubscriber(rdb, notificationSrv, logger)

	go redisSubscriber.SubscribeAndProcess(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down notifier service...")
	cancel()

	logger.Info("Notifier service stopped.")
}
