# GOURL - Go-Based URL Shortener

A lightweight, self-hostable URL shortener API built with Go, similar to Bitly.

## Features

- ✅ Create short URLs via REST API
- ✅ Redirect to original URLs
- ✅ **Analytics tracking** - Track clicks, unique visitors, time-based stats
- ✅ **Enhanced Analytics** - Daily clicks, top referrers, user agent breakdown
- ✅ **JWT Authentication** - Secure API access with user accounts
- ✅ **Rate limiting** - Protect API from abuse
- ✅ **CORS support** - Cross-origin requests
- ✅ **Docker support** - Easy deployment with Docker
- ✅ SQLite database for storage
- ✅ Base62 code generation
- ✅ Production-ready middleware and error handling

## Prerequisites

- Go 1.21 or higher
- SQLite3 (usually comes with Go SQLite driver)

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd GOURL
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

## Usage

### Create a Short URL

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

Response:
```json
{
  "short_url": "localhost:8080/abc12345",
  "original_url": "https://example.com",
  "code": "abc12345",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### Redirect to Original URL

Simply visit: `http://localhost:8080/{code}`

Or use curl:
```bash
curl -L http://localhost:8080/abc12345
```

**Note:** Each redirect is automatically tracked for analytics.

### Get Analytics/Stats

**Basic Stats:**
```bash
curl http://localhost:8080/api/stats/{code}
```

Response:
```json
{
  "code": "abc12345",
  "original_url": "https://example.com",
  "created_at": "2024-01-01T12:00:00Z",
  "total_clicks": 42,
  "unique_ips": 15
}
```

**Enhanced Stats (with time-based analytics):**
```bash
curl http://localhost:8080/api/stats/{code}/enhanced
```

Response includes:
- Total clicks and unique IPs
- Clicks by day (last 30 days)
- Top 10 referrers
- User agent breakdown

### Authentication

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Use token for authenticated requests:**
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"url": "https://example.com"}'
```

Note: Authentication is optional. URLs can be created without authentication, but authenticated users can track their URLs.

### Health Check

```bash
curl http://localhost:8080/health
```

## Configuration

Environment variables:
- `PORT` - Server port (default: 8080)
- `DB_PATH` - SQLite database path (default: gourl.db)
- `ENV` - Environment mode: `development` or `production` (default: development)
- `RATE_LIMIT_RPS` - Rate limit requests per second per IP (default: 10)
- `RATE_LIMIT_BURST` - Rate limit burst size (default: 20)
- `CORS_ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins (default: `*`)
- `JWT_SECRET` - Secret key for JWT tokens (default: development key, **change in production!**)

### Example Configuration

```bash
# Production setup
export PORT=8080
export DB_PATH=/var/lib/gourl/gourl.db
export ENV=production
export RATE_LIMIT_RPS=5
export RATE_LIMIT_BURST=10
export CORS_ALLOWED_ORIGINS=https://example.com,https://app.example.com
export JWT_SECRET=$(openssl rand -hex 32)  # Generate secure secret

go run cmd/server/main.go
```

### Rate Limiting

The API endpoints (`/api/shorten` and `/api/stats/:code`) are protected by rate limiting:
- Default: 10 requests per second per IP
- Burst: 20 requests
- Redirects (`/{code}`) are not rate limited

When rate limit is exceeded, the API returns:
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```
Status code: `429 Too Many Requests`

## Project Structure

```
GOURL/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers (URL, auth, analytics)
│   ├── models/              # Data models
│   ├── database/            # Database setup
│   ├── middleware/          # CORS, rate limiting, error handling
│   ├── config/              # Configuration management
│   ├── auth/                # JWT authentication utilities
│   └── utils/               # Utility functions
├── Dockerfile               # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
├── Makefile                # Common tasks
├── go.mod
└── README.md
```

## Deployment

### Quick Deploy Options

**Recommended: Railway** (Easiest, best for Go apps)
- See [DEPLOYMENT.md](./DEPLOYMENT.md) for Railway setup
- Free tier: $5/month credit
- Persistent storage included
- Auto-deploys from GitHub

**Alternative: Render** (Free tier available)
- See [DEPLOYMENT_RENDER.md](./DEPLOYMENT_RENDER.md)
- Free tier with auto-spin-down
- Good for demos/testing

**Alternative: Fly.io** (Production-ready)
- See [DEPLOYMENT_FLY.md](./DEPLOYMENT_FLY.md)
- Generous free tier
- Global edge network

### Docker Deployment

### Using Docker Compose (Recommended)

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Using Docker directly

```bash
# Build image
docker build -t gourl:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -v gourl_data:/data \
  -e JWT_SECRET=your-secret-key \
  -e ENV=production \
  gourl:latest
```

## Makefile Commands

The project includes a Makefile for common tasks:

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
make docker-stop   # Stop Docker containers
make clean         # Clean build artifacts
make dev           # Run in development mode
make prod          # Run in production mode
```

## Development

Run tests (when added):
```bash
go test ./...
```

Build binary:
```bash
go build -o gourl cmd/server/main.go
```

Or use Makefile:
```bash
make build
```

## License

MIT


