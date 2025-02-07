# Setting Up and Running Caching Proxy Server with Go

https://roadmap.sh/projects/caching-server

## Overview

This project is a **Caching Proxy Server** built in Go. It acts as an intermediary between clients and an origin server, caching responses to improve performance and reduce load on the origin server. The server uses **Redis** for caching and supports dynamic configuration of ports and origin servers.

---

## Prerequisites

Ensure you have the following tools installed on your system:

- **Golang** ([Download](https://go.dev/dl/))
- **Make** (Available by default on most Linux/macOS systems, installable on Windows via [Chocolatey](https://chocolatey.org/) or [Scoop](https://scoop.sh/))
- **Docker** ([Install Docker](https://docs.docker.com/get-docker/))
- **Docker Compose** ([Install Docker Compose](https://docs.docker.com/compose/install/))

---

## Steps to Set Up the Application

### 1. Start Redis with Docker

Run the following command to start a Redis container:

```sh
make start-redis
```

This command will:

- Pull the necessary Redis Docker image (if not already available)
- Start a Redis container using `docker-compose.yml`
- Expose Redis on port `6379`

---

### 2. Build the Application

To build the Go application, run:

```sh
make build
```

This will:

- Compile the Go code and generate a binary named `main`.

---

### 3. Run the Application

To start the caching proxy server, use:

```sh
make run PORT=<port> ORIGIN=<origin_url>
```

#### Example:

```sh
make run PORT=3000 ORIGIN=http://dummyjson.com
```

This will:

- Start the caching proxy server on the specified port.
- Proxy requests to the provided origin URL.
- Cache responses in Redis.

---

### 4. Clear the Cache

To clear the Redis cache, run:

```sh
make clear-cache
```

This will:

- Remove all cached data from Redis.

---

### 5. Stop Redis

To stop the Redis container, use:

```sh
make stop-redis
```

This will:

- Stop and remove the Redis container.

---

### 6. Clean Up

To remove the compiled binary and other generated files, run:

```sh
make clean
```

This will:

- Delete the `main` binary.

---

## Additional Commands

- **Check Redis Connection**

  ```sh
  make check-redis
  ```

  Verifies if the Redis container is running and accessible.

- **Run Tests**

  ```sh
  make test
  ```

  Runs all unit tests in the project (if available).

---

## Configuration

### Environment Variables

You can configure the application using the following environment variables:

| Variable       | Description                          | Default Value           |
|----------------|--------------------------------------|-------------------------|
| `PORT`         | Port to run the caching proxy on     | `3000`                  |
| `ORIGIN`       | Origin server to proxy requests to   | `http://dummyjson.com`  |
| `REDIS_ADDR`   | Redis server address                 | `localhost:6379`        |

---

## Example Workflow

1. Start Redis:

   ```sh
   make start-redis
   ```

2. Build and run the application:

   ```sh
   make run PORT=4000 ORIGIN=http://api.example.com
   ```

3. Test the caching proxy by making requests:

   ```sh
   curl http://localhost:4000/products
   ```

4. Clear the cache:

   ```sh
   make clear-cache
   ```

5. Stop Redis:

   ```sh
   make stop-redis
   ```

6. Clean up:

   ```sh
   make clean
   ```

---

## Troubleshooting

### 1. Redis Container Fails to Start

- Ensure Docker is running.
- Check for port conflicts (default Redis runs on `6379`).
- Run `docker ps` to see if the container is already running.

### 2. Application Fails to Connect to Redis

- Verify that Redis is running using `docker ps`.
- Ensure the `REDIS_ADDR` environment variable matches the Redis container's address.

### 3. Cache Not Working

- Check if Redis is properly configured and running.
- Verify that the application is storing and retrieving data from Redis.

---

## Project Structure

```
.
├── cmd
│   └── main.go            # Entry point of the application
├── internal
│   └── server.go          # Core logic for the caching proxy server
├── Makefile               # Automation commands
├── docker-compose.yml     # Docker Compose configuration for Redis
└── README.md              # Project documentation
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Now your **Caching Proxy Server** is ready to use! 🚀