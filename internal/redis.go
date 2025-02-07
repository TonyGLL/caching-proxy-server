package internal

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// InitRedis inicializa la conexión con Redis
func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr, // Dirección de Redis, por ejemplo: "localhost:6379"
		Password: "",   // Contraseña (si es necesaria)
		DB:       0,    // Número de base de datos
	})

	// Verificar la conexión con Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to Redis: %v", err))
	}
	fmt.Println("Connected to Redis")
}
