package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/TonyGLL/caching-proxy-server/pkg/config"
	"github.com/go-redis/redis/v8"
)

// Server holds the dependencies for the HTTP server
type Server struct {
	httpServer *http.Server
	redis      *redis.Client
	cfg        *config.Config
	log        *slog.Logger
}

// NewServer creates and configures a new server
func NewServer(cfg *config.Config, redisClient *redis.Client, log *slog.Logger) *Server {
	proxy, err := NewProxy(cfg, redisClient, log)
	if err != nil {
		log.Error("Failed to create proxy", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxy)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: mux,
		},
		redis: redisClient,
		cfg:   cfg,
		log:   log,
	}
}

// Start runs the HTTP server
func (s *Server) Start() error {
	s.log.Info("Server listening", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
