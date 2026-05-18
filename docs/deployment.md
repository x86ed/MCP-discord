# Deployment Guide

Production deployment guide for the MCP-Discord bot.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Deployment Options](#deployment-options)
- [Systemd Deployment (Linux)](#systemd-deployment-linux)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployments](#cloud-deployments)
- [Configuration Management](#configuration-management)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)

## Prerequisites

### System Requirements

- **CPU**: 1 core minimum, 2+ cores recommended
- **RAM**: 256MB minimum, 512MB+ recommended
- **Disk**: 100MB for binary, additional space for MCP servers
- **Network**: Outbound HTTPS to Discord API (discord.com)

### Software Requirements

- Go 1.21+ (for building from source)
- Git (for cloning repository)
- MCP server dependencies (Node.js, Python, etc. depending on your MCP server)

### Discord Setup

1. Create application at [Discord Developer Portal](https://discord.com/developers/applications)
2. Create bot user under "Bot" section
3. Enable required Privileged Gateway Intents if needed:
   - `GUILD_MEMBERS` (if using member info)
   - `MESSAGE_CONTENT` (if reading message content)
4. Generate OAuth2 URL with scopes:
   - `bot`
   - `applications.commands`
5. Set permissions (minimum):
   - `Send Messages`
   - `Use Slash Commands`

## Deployment Options

### Quick Comparison

| Method | Complexity | Isolation | Auto-Restart | Resource Usage |
|--------|-----------|-----------|--------------|----------------|
| **Systemd** | Low | Process | Yes | Low |
| **Docker** | Medium | Container | Yes | Medium |
| **Kubernetes** | High | Pod | Yes | Medium-High |
| **Cloud Services** | Low-Medium | Varies | Yes | Varies |

## Systemd Deployment (Linux)

Best for: VPS, dedicated servers, self-hosted environments

### 1. Build the Binary

```bash
# Clone repository
git clone https://github.com/yourusername/MCP-discord.git
cd MCP-discord

# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o mcpdiscord \
  ./cmd/mcpdiscord

# Verify build
./mcpdiscord --version
```

### 2. Install Binary

```bash
# Copy binary to system location
sudo cp mcpdiscord /usr/local/bin/
sudo chmod +x /usr/local/bin/mcpdiscord
```

### 3. Create Configuration

```bash
# Create configuration directory
sudo mkdir -p /etc/mcpdiscord
sudo chmod 755 /etc/mcpdiscord

# Create config file (edit as needed)
sudo nano /etc/mcpdiscord/config.json
```

**Example `/etc/mcpdiscord/config.json`:**
```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": ""
  },
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"],
    "env": {},
    "transport": "stdio"
  }
}
```

### 4. Create Environment File

```bash
# Create environment file for secrets
sudo nano /etc/mcpdiscord/environment
```

**Example `/etc/mcpdiscord/environment`:**
```bash
DISCORD_BOT_TOKEN=your_actual_bot_token_here
# Add other environment variables as needed
```

**Secure the file:**
```bash
sudo chmod 600 /etc/mcpdiscord/environment
sudo chown root:root /etc/mcpdiscord/environment
```

### 5. Create Systemd Service

See [examples/deployment/mcpdiscord.service](../examples/deployment/mcpdiscord.service) for the complete service file.

```bash
# Copy service file
sudo cp examples/deployment/mcpdiscord.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload
```

### 6. Start and Enable Service

```bash
# Start the service
sudo systemctl start mcpdiscord

# Check status
sudo systemctl status mcpdiscord

# Enable auto-start on boot
sudo systemctl enable mcpdiscord

# View logs
sudo journalctl -u mcpdiscord -f
```

### 7. Verify Deployment

```bash
# Check service is running
sudo systemctl is-active mcpdiscord

# Check bot is connected to Discord
# Commands should appear in your Discord server

# Test a command
# Type /help or one of your MCP server's commands
```

## Docker Deployment

Best for: Containerized environments, consistent deployments, easy scaling

### 1. Create Dockerfile

See [examples/deployment/Dockerfile](../examples/deployment/Dockerfile) for the complete Dockerfile.

### 2. Build Image

```bash
# Build image
docker build -t mcpdiscord:latest .

# Tag for registry (optional)
docker tag mcpdiscord:latest registry.example.com/mcpdiscord:latest
```

### 3. Create docker-compose.yml

See [examples/deployment/docker-compose.yml](../examples/deployment/docker-compose.yml) for the complete compose file.

### 4. Deploy with Docker Compose

```bash
# Create .env file with secrets
cat > .env << EOF
DISCORD_BOT_TOKEN=your_token_here
EOF

# Start container
docker-compose up -d

# View logs
docker-compose logs -f mcpdiscord

# Stop container
docker-compose down
```

### 5. Health Check

```bash
# Check container status
docker-compose ps

# Check logs for errors
docker-compose logs mcpdiscord | grep -i error

# Restart if needed
docker-compose restart mcpdiscord
```

## Cloud Deployments

### AWS EC2

1. Launch EC2 instance (t3.micro or larger)
2. Install dependencies (Go, Node.js, etc.)
3. Follow [Systemd Deployment](#systemd-deployment-linux) steps
4. Configure security group:
   - Allow outbound HTTPS to Discord (443)
   - No inbound ports needed

### Google Cloud Platform (Compute Engine)

Similar to AWS EC2:
1. Create Compute Engine instance
2. Follow Systemd deployment steps
3. Configure firewall rules for outbound HTTPS

### DigitalOcean Droplet

1. Create droplet (1GB RAM recommended)
2. SSH into droplet
3. Follow [Systemd Deployment](#systemd-deployment-linux) steps
4. (Optional) Set up monitoring via DigitalOcean dashboard

### Heroku

Create `Procfile`:
```
worker: ./mcpdiscord --config config.json
```

Deploy:
```bash
# Login
heroku login

# Create app
heroku create my-mcp-discord-bot

# Set environment variables
heroku config:set DISCORD_BOT_TOKEN=your_token

# Deploy
git push heroku main

# Check logs
heroku logs --tail
```

### Railway

1. Connect GitHub repository to Railway
2. Set environment variables in Railway dashboard
3. Deploy automatically on push
4. Monitor via Railway dashboard

## Configuration Management

### Environment Variables Best Practices

**DO:**
- ✅ Use environment variables for all secrets
- ✅ Use different configs for dev/staging/prod
- ✅ Document required environment variables
- ✅ Use secret management tools (Vault, AWS Secrets Manager)

**DON'T:**
- ❌ Commit tokens to version control
- ❌ Share production configs via chat/email
- ❌ Use same tokens for dev and prod
- ❌ Log sensitive values

### Secret Management

**AWS Secrets Manager:**
```bash
# Store secret
aws secretsmanager create-secret \
  --name mcpdiscord/token \
  --secret-string "your-token"

# Retrieve in deployment script
export DISCORD_BOT_TOKEN=$(aws secretsmanager get-secret-value \
  --secret-id mcpdiscord/token \
  --query SecretString \
  --output text)
```

**HashiCorp Vault:**
```bash
# Store secret
vault kv put secret/mcpdiscord token="your-token"

# Retrieve
export DISCORD_BOT_TOKEN=$(vault kv get -field=token secret/mcpdiscord)
```

## Monitoring

### Logs

**Systemd (journalctl):**
```bash
# Follow logs in real-time
sudo journalctl -u mcpdiscord -f

# Show last 100 lines
sudo journalctl -u mcpdiscord -n 100

# Show errors only
sudo journalctl -u mcpdiscord -p err
```

**Docker:**
```bash
# Follow logs
docker-compose logs -f mcpdiscord

# Last 100 lines
docker-compose logs --tail=100 mcpdiscord
```

### Health Checks

Create a simple health check script:

```bash
#!/bin/bash
# /usr/local/bin/mcpdiscord-health-check

if systemctl is-active --quiet mcpdiscord; then
    echo "OK: Service is running"
    exit 0
else
    echo "ERROR: Service is not running"
    exit 1
fi
```

**Add to cron for monitoring:**
```bash
# Check every 5 minutes
*/5 * * * * /usr/local/bin/mcpdiscord-health-check || systemctl restart mcpdiscord
```

### Metrics (Future Enhancement)

Once Prometheus metrics are added:
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'mcpdiscord'
    static_configs:
      - targets: ['localhost:2112']
```

## Troubleshooting

### Service Won't Start

**Check logs:**
```bash
sudo journalctl -u mcpdiscord -n 50
```

**Common issues:**
- Missing environment variables
- Invalid configuration
- Port already in use
- Insufficient permissions

### Bot Appears Offline

**Check Discord status:**
- Verify token is correct
- Check Discord API status (status.discord.com)
- Ensure bot has been invited to server

**Verify network:**
```bash
# Test Discord API connectivity
curl -I https://discord.com/api/v10/gateway
```

### Commands Not Appearing

**Solutions:**
1. Wait up to 1 hour for global commands (or use guildId for instant)
2. Verify bot has `applications.commands` scope
3. Check logs for registration errors
4. Try reinviting bot with updated permissions

### High Memory Usage

**Investigate:**
```bash
# Check memory usage
systemctl status mcpdiscord

# View detailed process info
ps aux | grep mcpdiscord

# Monitor over time
top -p $(pgrep mcpdiscord)
```

**Solutions:**
- Restart service periodically (cron job)
- Upgrade to larger instance
- Check for memory leaks in MCP server

## Security Best Practices

### Access Control

```bash
# Run as dedicated user (not root)
sudo useradd -r -s /bin/false mcpdiscord

# Set ownership
sudo chown -R mcpdiscord:mcpdiscord /etc/mcpdiscord

# Restrict permissions
sudo chmod 700 /etc/mcpdiscord
sudo chmod 600 /etc/mcpdiscord/environment
```

### Firewall

```bash
# UFW (Ubuntu)
sudo ufw allow out 443/tcp  # Discord API
sudo ufw enable

# iptables
sudo iptables -A OUTPUT -p tcp --dport 443 -j ACCEPT
```

### Updates

```bash
# Update bot
cd /path/to/MCP-discord
git pull
go build -o mcpdiscord ./cmd/mcpdiscord
sudo cp mcpdiscord /usr/local/bin/
sudo systemctl restart mcpdiscord

# Update dependencies
go get -u ./...
go mod tidy
```

### Backup

```bash
# Backup configuration
sudo tar -czf mcpdiscord-backup-$(date +%Y%m%d).tar.gz \
  /etc/mcpdiscord/ \
  /usr/local/bin/mcpdiscord

# Backup to S3
aws s3 cp mcpdiscord-backup-*.tar.gz s3://my-backups/mcpdiscord/
```

## Scaling Considerations

### Single Instance
- Sufficient for most use cases
- One bot token = one instance
- Handles thousands of commands per hour

### High Availability
- Run multiple instances with different tokens
- Use load balancer for MCP servers
- Consider bot sharding for 2500+ guilds

### Kubernetes Deployment
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mcpdiscord
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mcpdiscord
  template:
    metadata:
      labels:
        app: mcpdiscord
    spec:
      containers:
      - name: mcpdiscord
        image: mcpdiscord:latest
        env:
        - name: DISCORD_BOT_TOKEN
          valueFrom:
            secretKeyRef:
              name: discord-secrets
              key: bot-token
```

---

For more deployment examples, see the [examples/deployment/](../examples/deployment/) directory.
