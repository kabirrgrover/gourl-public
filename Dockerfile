# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o gourl ./cmd/server

# Final stage
FROM alpine:latest

# Install SQLite
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gourl .

# Copy database initialization (if needed)
# COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/data/gourl.db
ENV ENV=production

# Create data directory
RUN mkdir -p /data

# Run the application
CMD ["./gourl"]

