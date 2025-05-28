# URL Shortener Service

Scalable URL shortener service built with Go, Redis, and PostgreSQL.

## Features

- Shorten long URLs to short codes
- Redirect short codes to original URLs
- Redis caching for improved performance
- PostgreSQL for persistent storage
- Docker containerization

## Tech Stack

- Go 1.21
- Gin Web Framework
- Redis for caching
- PostgreSQL for storage
- Docker for containerization

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/url_shortner.git
   cd url_shortner
   ```

2. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. The service will be available at `http://localhost:8080`

## API Endpoints

### Create Short URL

```http
POST /shorten
Content-Type: application/json

{
    "url": "https://example.com/very/long/url"
}
```

Response:
```json
{
    "short_code": "abc123"
}
```

### Redirect to Original URL

```http
GET /{short_code}
```

Response: HTTP 302 redirect to the original URL

## Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the application locally:
   ```bash
   go run cmd/shortener/main.go
   ```

## Testing

Run the tests:
```bash
go test ./...
```

## Performance

The service is designed to handle high concurrency with response times under 20ms. It achieves this through:

- Redis caching with 24-hour TTL
- Optimized PostgreSQL queries with proper indexing
- Efficient URL code generation
