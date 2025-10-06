package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TonyGLL/caching-proxy-server/internal"
	"github.com/TonyGLL/caching-proxy-server/pkg/config"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Initialize Redis client
	redisClient, err := internal.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Create a new server
	server := internal.NewServer(cfg, redisClient)

	// Set up graceful shutdown
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	fmt.Printf("Server started on port %d, proxying to %s\n", cfg.Port, cfg.OriginURL)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server exited gracefully")
}