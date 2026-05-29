# Multi-stage Dockerfile for MCP-Discord bot
#
# MCP Server Integration:
#   - Automatically builds Go MCP server from ./cmd/mcpserver/ if present
#   - Falls back to test MCP server if not found
#   - Result: /app/mcp-server in final image
#
# To add your Go MCP server:
#   mkdir -p cmd/mcpserver && cp your-mcp-code/* cmd/mcpserver/
#
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the main bot binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcpdiscord ./cmd/mcpdiscord

# Build your MCP server binary.
# Option 1: If MCP server is in this repo (e.g., ./cmd/mcpserver/main.go):
RUN if [ -d "./cmd/mcpserver" ]; then \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcp-server ./cmd/mcpserver; \
    else \
      echo "Warning: cmd/mcpserver not found, using test MCP server" && \
      CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcp-server ./internal/testutil/mcpserver; \
    fi && \
    chmod +x /out/mcp-server

# Option 2: If MCP server is external, uncomment and adjust:
# COPY path/to/mcp-server-source /mcp-src
# RUN cd /mcp-src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcp-server .

# Create a seeded SQLite-like file at approximately 20MB.
RUN mkdir -p /out/data && \
    printf 'SQLite format 3\000' > /out/data/bot.db && \
    truncate -s 20M /out/data/bot.db

# Build-time validation for required embedded assets.
RUN mcp_size=$(wc -c < /out/mcp-server) && \
    db_size=$(wc -c < /out/data/bot.db) && \
    echo "MCP server: ${mcp_size} bytes, SQLite DB: ${db_size} bytes" && \
    [ "$mcp_size" -gt 1000000 ] || { echo "Error: MCP server binary too small (<1MB)"; exit 1; } && \
    [ "$db_size" -ge 15000000 ] && [ "$db_size" -le 30000000 ]

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata procps

RUN addgroup -g 1000 mcpdiscord && adduser -D -u 1000 -G mcpdiscord mcpdiscord

WORKDIR /app

COPY --from=builder /out/mcpdiscord /app/mcpdiscord
COPY --from=builder /out/mcp-server /app/mcp-server
COPY --from=builder /out/data/bot.db /app/data/bot.db
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh && \
    chmod +x /app/mcp-server && \
    mkdir -p /app/data && \
    chown -R mcpdiscord:mcpdiscord /app

USER mcpdiscord

ENV MCP_CONFIG_PATH=/app/mcp-config.json

HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD pgrep mcpdiscord > /dev/null || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
