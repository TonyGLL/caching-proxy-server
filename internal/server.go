package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/TonyGLL/caching-proxy-server/pkg/config"
	"github.com/go-redis/redis/v8"
)

// Server holds the dependencies for the HTTP server
type Server struct {
	httpServer *http.Server
	redis      *redis.Client
	cfg        *config.Config
}

// NewServer creates and configures a new server
func NewServer(cfg *config.Config, redisClient *redis.Client) *Server {
	proxy, err := NewProxy(cfg, redisClient)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
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
	}
}

// Start runs the HTTP server
func (s *Server) Start() error {
	log.Printf("Server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}