# Multi-stage build for optimized image size
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o mcpdiscord \
    ./cmd/mcpdiscord

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    nodejs \
    npm \
    tzdata

# Create non-root user
RUN addgroup -g 1000 mcpdiscord && \
    adduser -D -u 1000 -G mcpdiscord mcpdiscord

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mcpdiscord /app/mcpdiscord

# Copy example config (can be overridden via volume mount)
COPY examples/with-env/mcp-config.json /app/mcp-config.json

# Set ownership
RUN chown -R mcpdiscord:mcpdiscord /app

# Switch to non-root user
USER mcpdiscord

# Health check (checks if process is running)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep mcpdiscord || exit 1

# Run the bot
ENTRYPOINT ["/app/mcpdiscord"]
CMD ["--config", "/app/mcp-config.json"]
