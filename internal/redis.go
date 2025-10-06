package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client and verifies the connection
func NewRedisClient(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Verify the connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}