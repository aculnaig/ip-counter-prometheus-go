# Log Tracker Service

Production-grade HTTP service that tracks unique IP addresses from log entries and exposes metrics in Prometheus format.

## Features

- Thread-safe IP tracking with concurrent request handling
- Dual HTTP servers (log ingestion and metrics)
- Graceful shutdown with configurable timeout
- Structured JSON logging with configurable levels
- Environment-based configuration
- Health check endpoints
- Request validation and sanitization
- Middleware for logging and panic recovery
- Prometheus-compatible metrics endpoint
- Docker and Docker Compose support

## Project Structure


```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   └── log.go               # Data models and validation
│   ├── server/
│   │   ├── log_server.go        # Log ingestion server
│   │   └── metrics_server.go   # Metrics server
│   └── tracker/
│       └── tracker.go           # IP tracking logic
├── pkg/
│   ├── logger/
│   │   └── logger.go            # Structured logger
│   └── middleware/
│       └── middleware.go        # HTTP middlewares
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Configuration

Set environment variables:

```bash
LOG_LEVEL=info                       # debug, info, warn, error
LOG_SERVER_PORT=5000
LOG_SERVER_READ_TIMEOUT=10s
LOG_SERVER_WRITE_TIMEOUT=10s
METRICS_SERVER_PORT=9102
```

## Running

### Local Development

```bash
# Build
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Run locally
make run
```

### Docker

```bash
# Build image
make docker-build

# Run with Docker Compose
make docker-run
```

## API Endpoints

### Log Server (Port 5000)

**POST /logs**
```bash
curl -X POST http://localhost:5000/logs \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2025-10-24T10:00:00Z","ip":"192.168.1.1","url":"/api/users"}'
```

### Metrics Server (Port 9102)

**GET /metrics**
```bash
curl http://localhost:9102/metrics
```

Returns Prometheus format:
```
# HELP unique_ip_addresses Total number of unique IP addresses seen
# TYPE unique_ip_addresses gauge
unique_ip_addresses 42
```

## Testing

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Production Considerations

1. **Monitoring**: Integrate with Prometheus for metrics collection
2. **Load Balancing**: Deploy behind a load balancer for horizontal scaling
3. **Rate Limiting**: Add rate limiting middleware for production traffic
4. **TLS**: Enable HTTPS with proper certificates
5. **Logging**: Ship logs to centralized logging system (ELK, Loki, etc.)
6. **Persistence**: Consider persistent storage for IP tracking if needed across restarts
7. **Memory Management**: Monitor memory usage; implement TTL for old IPs if needed
