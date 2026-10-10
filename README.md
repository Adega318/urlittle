# URLittle

![CI](https://github.com/Adega318/urlittle/actions/workflows/ci.yml/badge.svg)

A minimal URL shortener service written in Go.

## Features

- **Shorten URLs**: Create short URLs via HTTP POST
- **Redirect**: Resolve short URLs with HTTP GET
- **Health check**: Built-in `/health` endpoint
- **URL normalization**: Enforces HTTPS, strips credentials, removes fragments and default ports
- **TTL expiration**: Links expire after 10 minutes (configurable in code)
- **Docker support**: Multi-stage build, security-hardened runtime

## API

### Create short URL

```bash
POST /
Content-Type: text/plain

https://example.com/very/long/url
```

Response:

```
HTTP/1.1 201 Created
Location: http://localhost:8080/abc123
Content-Type: text/plain; charset=utf-8

http://localhost:8080/abc123
```

- Request body: raw URL string (max 8 KB)
- Returns the shortened URL in `Location` header and response body

### Resolve short URL

```bash
GET /abc123
```

Response:

```
HTTP/1.1 307 Temporary Redirect
Location: https://example.com/very/long/url
```

- Returns 404 if the ID is not found or expired

### Health check

```bash
GET /health
```

Response:

```
HTTP/1.1 204 No Content
```

## Configuration

All configuration is via environment variables:

| Variable       | Required | Description                                 |
| -------------- | -------- | ------------------------------------------- |
| `LOG_LEVEL`    | Yes      | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `PORT`         | Yes      | HTTP server port (e.g., `8080`)             |
| `TTL_MINUTES`  | Yes      | Time to live in minutes                     |
| `DATABASE_URL` | Yes      | PostgreSQL connection string                |
| `CACHE_SIZE`   | Yes      | LRU cache capacity (number of entries)      |

Example `.env`:

```
LOG_LEVEL=INFO
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/urlittle
CACHE_SIZE=1000
```

## Running locally

### Prerequisites

- Go 1.27.1+
- PostgreSQL 18+

### With Docker Compose (recommended)

```bash
# Copy example env and adjust values
cp .env.example .env

# Start services
docker compose up -d
```

The service will be available at `http://localhost:8080`.

### Without Docker

```bash
# Start PostgreSQL and run migrations
# (see migrate/001_urls_table.sql)

# Set environment variables
export LOG_LEVEL=INFO
export PORT=8080
export TTL_MINUTES=5
export DATABASE_URL=postgres://user:pass@localhost:5432/urlittle
export CACHE_SIZE=1000

# Run
go run .
```

## Development

### Build

```bash
go build -o urlittle .
```

### Run tests

```bash
go test ./...
```

### Lint

```bash
golangci-lint run
```

## Docker

### Production image

```dockerfile
# Multi-stage build
# Builder: golang:1.27.1-alpine
# Runtime: alpine:3.24 (non-root, read-only, dropped capabilities)
```

### Environment variables for docker-compose

| Variable            | Description                          |
| ------------------- | ------------------------------------ |
| `POSTGRES_DB`       | Database name                        |
| `POSTGRES_USER`     | Database user                        |
| `POSTGRES_PASSWORD` | Database password                    |
| `LOG_LEVEL`         | Log level (default: INFO)            |
| `PORT`              | Host port mapping (default: 8080)    |
| `TTL_MINUTES`       | Time to live in minutes (default: 5) |
| `CACHE_SIZE`        | LRU cache size (default: 1000)       |
| `VERSION`           | Build version tag (default: dev)     |

## Security

- Runs as non-root user (UID 65532)
- Read-only root filesystem
- All capabilities dropped
- No new privileges
- Health check via `wget`

## Roadmap

- [ ] **Rate limiting**: the API is currently unauthenticated and unthrottled, so a public deployment could be abused.
- [x] Configurable TTL (currently hardcoded to 10 minutes)
- [ ] Metrics endpoint (Prometheus)
- [ ] Kubernetes configuration

## License

MIT
