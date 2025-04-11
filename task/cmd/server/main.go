package main

import (
	"context"
	"log"
	"os"

	"github.com/CP-Payne/taskflow/task/internal/publisher"
	"github.com/CP-Payne/taskflow/task/internal/repository/memory"
	"github.com/CP-Payne/taskflow/task/internal/server"
	"github.com/CP-Payne/taskflow/task/internal/service"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load global config
	if err := godotenv.Load("./config/.env"); err != nil {
		logger.Panic("Failed to load global config")
	}

	redisAddr := os.Getenv("REDIS_NOTIFIER_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	redisPublisher := publisher.NewRedisPublisher(rdb, logger)

	// TODO: Create config to pass to layers
	// TODO: Define Handler in main, instead of StartGRPCServer
	repo := memory.NewInMemory()
	srv := service.New(repo, logger, redisPublisher)
	server.StartGRPCServer("localhost:3034", srv, logger)
}
