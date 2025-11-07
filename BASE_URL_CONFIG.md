# Base URL Configuration Guide

## Problem
When hosting your URL shortener, short URLs were showing `http://localhost:8080/hello` instead of your actual domain like `https://yoursite.com/hello`.

## Solution
The application now supports a `BASE_URL` environment variable to set your production domain.

## How to Use

### For Local Development
No configuration needed! The app will automatically detect the URL from the request:
- `http://localhost:8080/hello` ✅

### For Production (Vercel, Railway, Render, etc.)

Set the `BASE_URL` environment variable to your domain:

```bash
export BASE_URL=https://yoursite.com
```

**Important:** 
- Don't include a trailing slash
- Use `https://` for production
- Use your actual domain (e.g., `https://gourl.vercel.app`)

## Examples

### Vercel
In your Vercel dashboard:
1. Go to Project Settings → Environment Variables
2. Add: `BASE_URL` = `https://your-app.vercel.app`

### Railway
In Railway dashboard:
1. Go to Variables tab
2. Add: `BASE_URL` = `https://your-app.railway.app`

### Render
In Render dashboard:
1. Go to Environment tab
2. Add: `BASE_URL` = `https://your-app.onrender.com`

### Docker/Docker Compose
In your `.env` file or `docker-compose.yml`:
```yaml
environment:
  - BASE_URL=https://yoursite.com
```

### Local Testing
To test with a custom domain locally:
```bash
export BASE_URL=https://yoursite.com
go run cmd/server/main.go
```

## How It Works

1. **If `BASE_URL` is set**: Uses that value for all short URLs
2. **If `BASE_URL` is not set**: Automatically detects from the HTTP request
   - Checks `X-Forwarded-Proto` header (for proxies)
   - Uses request hostname
   - Handles HTTPS automatically

## Updated Files

- `internal/config/config.go` - Added `BaseURL` config field
- `internal/handlers/utils.go` - New `getBaseURL()` function
- `internal/handlers/url_handler.go` - Uses `getBaseURL()`
- `internal/handlers/qr_handler.go` - Uses `getBaseURL()`
- `internal/handlers/bulk_handler.go` - Uses `getBaseURL()`
- `cmd/server/main.go` - Stores config in context

## Testing

After setting `BASE_URL`, create a new short URL and verify it uses your domain:

```bash
curl -X POST https://yoursite.com/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

Response should show:
```json
{
  "short_url": "https://yoursite.com/abc123",
  ...
}
```

Not `http://localhost:8080/abc123` ✅

