# ==============================================================================
# Makefile for the Caching Proxy Server
# ==============================================================================

# Go parameters
BINARY_NAME=caching-proxy
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Environment variables (with default values)
export PORT ?= 8080
export ORIGIN_URL ?= https://dummyjson.com
export REDIS_ADDR ?= localhost:6379
export CACHE_EXPIRES ?= 10m
export LOG_LEVEL ?= INFO

# Docker parameters
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: all build run clean test tidy start-redis stop-redis check-redis

all: build

# ------------------------------------------------------------------------------
# Build and Run
# ------------------------------------------------------------------------------

build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd/main.go

run: build
	@echo "Starting server on port $(PORT)"
	@echo "Proxying to: $(ORIGIN_URL)"
	@./$(BINARY_NAME)

# ------------------------------------------------------------------------------
# Testing and Linting
# ------------------------------------------------------------------------------

test:
	@echo "Running tests..."
	@go test -v ./...

tidy:
	@echo "Tidying go modules..."
	@go mod tidy

# ------------------------------------------------------------------------------
# Docker and Redis
# ------------------------------------------------------------------------------

start-redis:
	@echo "Starting Redis container..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop-redis:
	@echo "Stopping Redis container..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

check-redis:
	@echo "Pinging Redis..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec redis redis-cli ping

# ------------------------------------------------------------------------------
# Cleanup
# ------------------------------------------------------------------------------

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@go clean -testcache