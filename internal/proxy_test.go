package internal

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TonyGLL/caching-proxy-server/pkg/config"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxy_CacheHit(t *testing.T) {
	// Setup
	db, mock := redismock.NewClientMock()
	cfg := &config.Config{
		OriginURL:    "http://dummy-origin.com",
		CacheExpires: 10 * time.Minute,
	}
	proxy, err := NewProxy(cfg, db)
	require.NoError(t, err)

	// Mock Redis GET
	cacheKey := "/test"
	expectedBody := []byte("cached response")
	mock.ExpectGet(cacheKey).SetVal(string(expectedBody))

	// Create request and response recorder
	req := httptest.NewRequest("GET", cacheKey, nil)
	rr := httptest.NewRecorder()

	// Execute
	proxy.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "HIT", rr.Header().Get("X-Cache"))
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, expectedBody, body)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProxy_CacheMiss(t *testing.T) {
	// Setup origin server
	originServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("origin response"))
	}))
	defer originServer.Close()

	// Setup proxy
	db, mock := redismock.NewClientMock()
	cfg := &config.Config{
		OriginURL:    originServer.URL,
		CacheExpires: 10 * time.Minute,
	}
	proxy, err := NewProxy(cfg, db)
	require.NoError(t, err)

	// Mock Redis GET (miss) and SET
	cacheKey := "/test"
	expectedBody := []byte("origin response")
	mock.ExpectGet(cacheKey).RedisNil()
	mock.ExpectSet(cacheKey, expectedBody, cfg.CacheExpires).SetVal("OK")

	// Create request and response recorder
	req := httptest.NewRequest("GET", cacheKey, nil)
	rr := httptest.NewRecorder()

	// Execute
	proxy.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "MISS", rr.Header().Get("X-Cache"))
	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, expectedBody, body)
	assert.NoError(t, mock.ExpectationsWereMet())
}