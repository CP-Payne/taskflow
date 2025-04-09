package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func subscribeAndSendNotification() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	pubsub := rdb.Subscribe(ctx, "task_notification")
	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println("Sending Notification:", msg.Payload)
		// Call Email/SMS API
	}
}

func main() {
}
