# Variables
BINARY_NAME=main
REDIS_CONTAINER_NAME=redis_container
REDIS_PORT=6379
PORT?=3000  # Valor por defecto para el puerto
ORIGIN?=http://dummyjson.com  # Valor por defecto para el origen

# Comando para construir el proyecto
build:
	@echo "Building the project..."
	@go build -o $(BINARY_NAME) cmd/main.go
	@echo "Build complete."

# Comando para ejecutar el proyecto
run: build
	@echo "Starting the application on port $(PORT) with origin $(ORIGIN)..."
	@./$(BINARY_NAME) --port $(PORT) --origin $(ORIGIN)

# Comando para limpiar la caché de Redis
clear-cache:
	@echo "Clearing Redis cache..."
	@./$(BINARY_NAME) --clear-cache

# Comando para levantar Redis con Docker Compose
start-redis:
	@echo "Starting Redis container..."
	@docker-compose up -d
	@echo "Redis is running on port $(REDIS_PORT)."

# Comando para detener Redis
stop-redis:
	@echo "Stopping Redis container..."
	@docker-compose down
	@echo "Redis container stopped."

# Comando para ejecutar tests (si tienes tests)
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests complete."

# Comando para limpiar archivos generados
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@echo "Cleanup complete."

# Comando para verificar la conexión a Redis
check-redis:
	@echo "Checking Redis connection..."
	@docker exec $(REDIS_CONTAINER_NAME) redis-cli ping