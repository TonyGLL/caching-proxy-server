package internal

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/TonyGLL/caching-proxy-server/pkg/config"
	"github.com/go-redis/redis/v8"
)

// Proxy handles the proxying logic
type Proxy struct {
	target *url.URL
	client *http.Client
	cache  *redis.Client
	cfg    *config.Config
	log    *slog.Logger
}

// NewProxy creates a new Proxy
func NewProxy(cfg *config.Config, cache *redis.Client, log *slog.Logger) (*Proxy, error) {
	target, err := url.Parse(cfg.OriginURL)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		target: target,
		client: &http.Client{Timeout: 30 * time.Second},
		cache:  cache,
		cfg:    cfg,
		log:    log,
	}, nil
}

// ServeHTTP handles the incoming request
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.URL.String()
	ctx := r.Context()
	requestLog := p.log.With("method", r.Method, "url", r.URL.String())

	// Try to get the response from cache
	cachedResponse, err := p.cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		requestLog.Info("Cache HIT", "cache_status", "HIT")
		w.Header().Set("X-Cache", "HIT")
		w.Write(cachedResponse)
		return
	}

	if err != redis.Nil {
		requestLog.Warn("Redis error on GET", "error", err)
	}
	requestLog.Info("Cache MISS", "cache_status", "MISS")

	// Create a new request to the target
	proxyReq := r.Clone(ctx)
	proxyReq.Host = p.target.Host
	proxyReq.URL.Scheme = p.target.Scheme
	proxyReq.URL.Host = p.target.Host
	proxyReq.RequestURI = "" // RequestURI cannot be set in a client request

	// Forward the request to the origin server
	resp, err := p.client.Do(proxyReq)
	if err != nil {
		requestLog.Error("Error forwarding request", "error", err)
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		requestLog.Error("Error reading response body", "error", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Cache the response if it's a successful one
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		err := p.cache.Set(ctx, cacheKey, body, p.cfg.CacheExpires).Err()
		if err != nil {
			requestLog.Warn("Failed to cache response", "error", err)
		}
	}

	// Copy headers and status code to the response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
