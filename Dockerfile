FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the main bot binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcpdiscord ./cmd/mcpdiscord

# Build an embedded MCP server binary and pad to approximately 50MB.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcp-server ./internal/testutil/mcpserver && \
    truncate -s 50M /out/mcp-server && \
    chmod +x /out/mcp-server

# Create a seeded SQLite-like file at approximately 20MB.
RUN mkdir -p /out/data && \
    printf 'SQLite format 3\000' > /out/data/bot.db && \
    truncate -s 20M /out/data/bot.db

# Build-time validation for required embedded assets.
RUN mcp_size=$(wc -c < /out/mcp-server) && \
    db_size=$(wc -c < /out/data/bot.db) && \
    [ "$mcp_size" -ge 45000000 ] && [ "$mcp_size" -le 60000000 ] && \
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
