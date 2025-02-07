package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/TonyGLL/caching-proxy-server/internal"
)

func main() {
	//* Definir flags
	port := flag.Int("port", 0, "Port to run the caching proxy on")
	origin := flag.String("origin", "", "Origin server to proxy requests to")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")

	//* Parsear los flags
	flag.Parse()

	//* Inicializar Redis
	internal.InitRedis("localhost:6379")

	if *clearCache {
		err := internal.ClearCache()
		if err != nil {
			log.Fatalf("Error clearing cache: %v", err)
		}
		fmt.Println("Cache cleared successfully")
		return
	}

	//* Iniciar el servidor con los parámetros proporcionados
	if *port == 0 {
		log.Fatal("Port server must be specified")
	}

	if *origin == "" {
		log.Fatal("Origin server must be specified")
	}

	internal.StartServer(*port, *origin)
}
