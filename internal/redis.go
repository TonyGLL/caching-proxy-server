package internal

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client and verifies the connection
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Verify the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
