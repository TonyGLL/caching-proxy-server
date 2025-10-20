package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TonyGLL/caching-proxy-server/internal"
	"github.com/TonyGLL/caching-proxy-server/pkg/config"
	"github.com/TonyGLL/caching-proxy-server/pkg/logger"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize Redis client
	redisClient, err := internal.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Create a new server
	server := internal.NewServer(cfg, redisClient, log)

	// Set up graceful shutdown
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("Could not start server", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started", "port", cfg.Port, "proxy_to", cfg.OriginURL)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}
