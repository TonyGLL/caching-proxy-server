# Go Caching Proxy Server

[![Go Report Card](https://goreportcard.com/badge/github.com/TonyGLL/caching-proxy-server)](https://goreportcard.com/report/github.com/TonyGLL/caching-proxy-server)

## Overview

This project is a high-performance **Caching Proxy Server** built in Go. It acts as an intermediary between clients and an origin server, caching responses in **Redis** to dramatically improve response times and reduce the load on the origin.

This refactored version includes:
- **Content-Agnostic Caching**: Caches any type of content, not just JSON.
- **Configuration via Environment Variables**: Easy to configure and deploy.
- **Graceful Shutdown**: Ensures no in-flight requests are dropped on shutdown.
- **Modular Codebase**: A clean, maintainable, and testable structure.
- **Unit Tests**: Robust tests to ensure reliability.

---

## Prerequisites

- **Go** (v1.18+)
- **Make**
- **Docker** & **Docker Compose**

---

## 🚀 Getting Started

### 1. Start Redis
To run the proxy, you first need a Redis instance. A `docker-compose.yml` file is provided for convenience.

```sh
make start-redis
```
This will start a Redis container in the background and expose it on port `6379`.

### 2. Run the Application
Once Redis is running, you can build and start the proxy server with a single command:

```sh
make run
```
By default, the server will start on port `8080` and proxy requests to `https://dummyjson.com`.

### 3. Test the Proxy
You can now send requests to the proxy. The `X-Cache` header in the response will tell you if you got a `HIT` or a `MISS`.

```sh
# First request (MISS)
curl -v http://localhost:8080/products/1

# Second request (HIT)
curl -v http://localhost:8080/products/1
```

---

## ⚙️ Configuration

The application is configured using environment variables. You can override the default values by setting them in your shell or including them in the `make run` command.

| Variable        | Description                                  | Default Value            |
|-----------------|----------------------------------------------|--------------------------|
| `PORT`          | Port for the proxy server to listen on.      | `8080`                   |
| `ORIGIN_URL`    | The target origin server to proxy requests to. | `https://dummyjson.com`  |
| `REDIS_ADDR`    | The address of the Redis server.             | `localhost:6379`         |
| `CACHE_EXPIRES` | Cache expiration time (e.g., `5m`, `1h`).    | `10m`                    |

#### Example with Custom Configuration:
```sh
export PORT=4000
export ORIGIN_URL=http://my-api.com
make run
```

---

## ✨ Available `make` Commands

- `make build`: Compiles the Go binary.
- `make run`: Builds and runs the server.
- `make test`: Runs all unit tests.
- `make tidy`: Tidies up Go module dependencies.
- `make clean`: Removes the compiled binary and test cache.
- `make start-redis`: Starts the Redis container.
- `make stop-redis`: Stops the Redis container.
- `make check-redis`: Pings the Redis container to check connectivity.

---

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point with graceful shutdown
├── internal/
│   ├── proxy.go             # Core proxy and caching logic
│   ├── proxy_test.go        # Unit tests for the proxy
│   ├── redis.go             # Redis client initialization
│   └── server.go            # Server setup and lifecycle management
├── pkg/
│   └── config/
│       └── config.go        # Configuration loader from environment variables
├── Makefile                 # Automation for build, run, test, etc.
├── go.mod                   # Go module dependencies
├── go.sum
└── docker-compose.yml       # Docker Compose for Redis service
```

---

## License

This project is licensed under the MIT License.