package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	userpb "github.com/CP-Payne/taskflow/pkg/gen/user/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gatewayPort     = ":8080"
	userServiceAddr = "localhost:9011"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	userConn, err := grpc.NewClient(userServiceAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial User service: %v", err)
	}
	defer userConn.Close()
	log.Printf("Connected to User service at %s", userServiceAddr)

	gwmux := runtime.NewServeMux()

	err = userpb.RegisterUserHandler(ctx, gwmux, userConn)
	if err != nil {
		log.Fatalf("Failed to register user gateway handler: %v", err)
	}
	log.Printf("Registered User service handler with gateway")

	mux := http.NewServeMux()

	mux.Handle("/", gwmux)

	log.Printf("Starting API Gateway on port %s", gatewayPort)
	httpServer := &http.Server{
		Addr:    gatewayPort,
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting.")
}
