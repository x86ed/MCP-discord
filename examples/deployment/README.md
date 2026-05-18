# Deployment Examples

This directory contains example deployment configurations for production environments.

## Files

- [`mcpdiscord.service`](mcpdiscord.service) - systemd service file for Linux systems
- [`Dockerfile`](Dockerfile) - Multi-stage Docker build
- [`docker-compose.yml`](docker-compose.yml) - Docker Compose configuration
- [`.env.example`](.env.example) - Example environment variables for Docker

## Systemd Deployment

### Quick Setup

```bash
# 1. Build and install binary
go build -o mcpdiscord ./cmd/mcpdiscord
sudo cp mcpdiscord /usr/local/bin/
sudo chmod +x /usr/local/bin/mcpdiscord

# 2. Create config directory
sudo mkdir -p /etc/mcpdiscord
sudo cp mcp-config.json /etc/mcpdiscord/config.json

# 3. Create environment file with secrets
sudo tee /etc/mcpdiscord/environment << EOF
DISCORD_BOT_TOKEN=your_token_here
EOF
sudo chmod 600 /etc/mcpdiscord/environment

# 4. Create service user
sudo useradd -r -s /bin/false mcpdiscord

# 5. Install service
sudo cp examples/deployment/mcpdiscord.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable mcpdiscord
sudo systemctl start mcpdiscord

# 6. Check status
sudo systemctl status mcpdiscord
sudo journalctl -u mcpdiscord -f
```

See [docs/deployment.md](../../docs/deployment.md) for detailed instructions.

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
# 1. Create .env file
cat > .env << EOF
DISCORD_BOT_TOKEN=your_token_here
DISCORD_GUILD_ID=
LOG_LEVEL=info
EOF

# 2. (Optional) Create custom config
cp ../../examples/basic/mcp-config.json ./mcp-config.json
# Edit mcp-config.json as needed

# 3. Start the bot
docker-compose up -d

# 4. View logs
docker-compose logs -f mcpdiscord

# 5. Stop the bot
docker-compose down
```

### Using Docker Directly

```bash
# 1. Build image
docker build -t mcpdiscord:latest -f examples/deployment/Dockerfile .

# 2. Run container
docker run -d \
  --name mcpdiscord \
  --restart unless-stopped \
  -e DISCORD_BOT_TOKEN="your_token_here" \
  -v $(pwd)/mcp-config.json:/app/mcp-config.json:ro \
  mcpdiscord:latest

# 3. View logs
docker logs -f mcpdiscord

# 4. Stop container
docker stop mcpdiscord
docker rm mcpdiscord
```

## Configuration

All deployment methods use the same configuration file format. See [docs/configuration.md](../../docs/configuration.md) for details.

### Example Production Config

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": ""
  },
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"],
    "env": {
      "API_KEY": "${WEATHER_API_KEY}"
    },
    "transport": "stdio"
  }
}
```

## Security Notes

- **Never commit secrets** - use environment variables
- **Restrict permissions** - config files should be mode 600
- **Use dedicated user** - don't run as root
- **Enable security features** - systemd hardening, Docker user namespaces
- **Monitor logs** - check for suspicious activity

## Resource Requirements

| Deployment | CPU | RAM | Disk |
|-----------|-----|-----|------|
| Minimal | 1 core | 256MB | 100MB |
| Recommended | 2 cores | 512MB | 500MB |
| High Volume | 4 cores | 1GB | 1GB |

## Troubleshooting

### Systemd Issues

```bash
# Check service status
sudo systemctl status mcpdiscord

# View logs
sudo journalctl -u mcpdiscord -n 100

# Restart service
sudo systemctl restart mcpdiscord

# Check configuration
sudo mcpdiscord --config /etc/mcpdiscord/config.json
```

### Docker Issues

```bash
# Check container logs
docker-compose logs mcpdiscord

# Check container status
docker-compose ps

# Restart container
docker-compose restart mcpdiscord

# Rebuild image
docker-compose build --no-cache
docker-compose up -d
```

## Additional Resources

- [Full Deployment Guide](../../docs/deployment.md)
- [Configuration Reference](../../docs/configuration.md)
- [Architecture Overview](../../docs/architecture.md)
