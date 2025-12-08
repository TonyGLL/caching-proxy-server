package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port         int
	OriginURL    string
	RedisAddr    string
	CacheExpires time.Duration
	LogLevel     string
}

// Load loads the configuration from environment variables
func Load() *Config {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "3000"))
	if err != nil {
		log.Fatalf("Invalid port format: %v", err)
	}

	cfg := &Config{
		Port:         port,
		OriginURL:    getEnv("ORIGIN_URL", "https://dummyjson.com"),
		RedisAddr:    getEnv("REDIS_ADDR", ""), // Default seguro
		CacheExpires: getDurationEnv("CACHE_EXPIRES", 10*time.Minute),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
	}

	// Build REDIS_URL if not provided: support user+password or only password
	if cfg.RedisAddr == "" {
		rHost := getEnv("REDIS_HOST", "")
		rPort := getEnv("REDIS_PORT", "")
		rUser := getEnv("REDIS_USER", "") // soporte para REDIS_USERNAME
		rPass := getEnv("REDIS_PASS", "")
		rDB := getEnv("REDIS_DB", "0")

		if rHost == "" || rPort == "" {
			return nil
		}

		// url-encode user and password
		encodedUser := url.QueryEscape(rUser)
		encodedPass := url.QueryEscape(rPass)

		var redisURL string
		if rUser != "" {
			// usuario presente
			redisURL = fmt.Sprintf("redis://%s:%s@%s:%s/%s", encodedUser, encodedPass, rHost, rPort, rDB)
		} else if rPass != "" {
			// solo password
			redisURL = fmt.Sprintf("redis://:%s@%s:%s/%s", encodedPass, rHost, rPort, rDB)
		} else {
			// sin auth
			redisURL = fmt.Sprintf("redis://%s:%s/%s", rHost, rPort, rDB)
		}

		cfg.RedisAddr = redisURL
	}

	return cfg
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getDurationEnv retrieves a duration environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration format for %s, using default: %v", key, err)
		return defaultValue
	}

	return value
}
