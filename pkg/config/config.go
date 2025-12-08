package config

import (
	"log"
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

	return &Config{
		Port:         port,
		OriginURL:    getEnv("ORIGIN_URL", "https://dummyjson.com"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"), // Default seguro
		CacheExpires: getDurationEnv("CACHE_EXPIRES", 10*time.Minute),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
	}
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
