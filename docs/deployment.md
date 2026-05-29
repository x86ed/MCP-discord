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

### 1. Prepare Your MCP Server Binary

The Dockerfile embeds your MCP server binary into the container image. You have several options:

#### Option A: Build Go MCP Server from Source (Recommended for Go)

If your MCP server is a Go program, add it to the repository and the Dockerfile will build it automatically.

**Structure 1: MCP server in same repo**

```bash
# Add your MCP server code to the repo
mkdir -p cmd/mcpserver
# Copy your main.go and other files
cp /path/to/your/mcp-server/*.go cmd/mcpserver/

# Add dependencies to go.mod if needed
go get your-mcp-dependencies

# Test build locally
go build -o mcp-server ./cmd/mcpserver
```

The Dockerfile automatically detects and builds `cmd/mcpserver` → `/app/mcp-server`

**Example cmd/mcpserver/main.go:**

```go
package main

import (
    "fmt"
    "os"
    // your MCP server imports
)

func main() {
    // Your MCP server implementation
    fmt.Fprintln(os.Stderr, "MCP server starting...")
    // ... handle stdio MCP protocol ...
}
```

**Structure 2: MCP server in separate directory**

```dockerfile
# In Dockerfile, uncomment Option 2 and adjust:
COPY path/to/mcp-server-source /mcp-src
RUN cd /mcp-src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /out/mcp-server .
```

**Structure 3: MCP server as Go module dependency**

```dockerfile
# In Dockerfile, add:
RUN go install github.com/yourorg/mcp-server@latest && \
    cp $(go env GOPATH)/bin/mcp-server /out/mcp-server
```

#### Option B: Copy Pre-built Binary from Local Filesystem

Place your MCP server binary in the repository root:

```bash
# Copy your binary to the repo
cp /path/to/your/mcp-server ./mcp-server
chmod +x ./mcp-server

# Verify it's a Linux binary (if building on macOS/Windows)
file ./mcp-server
# Should show: ELF 64-bit LSB executable

# Optional: Add to .gitignore if you don't want to commit it
echo "mcp-server" >> .gitignore
```

The Dockerfile will automatically copy it: `COPY mcp-server /out/mcp-server`

#### Option B: Download During Build

Add to Dockerfile before the COPY step:

```dockerfile
# Download your MCP binary
RUN wget -O /out/mcp-server https://your-server.com/mcp-server && \
    chmod +x /out/mcp-server
```

#### Option C: Build from Source in Docker

If your MCP server is buildable, add a build stage:

```dockerfile
# Add before the Go builder stage
FROM node:20-alpine AS mcp-builder
WORKDIR /mcp
COPY mcp-server-source/ .
RUN npm install && npm run build
RUN pkg . -t node20-linux-x64 -o mcp-server

# Then in the Go builder stage:
COPY --from=mcp-builder /mcp/mcp-server /out/mcp-server
```

#### Option D: Use npx/Node.js MCP Server

If using a Node.js-based MCP server (e.g., `@modelcontextprotocol/server-weather`), install Node.js in the runtime image:

```dockerfile
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata procps nodejs npm
# ... rest of Dockerfile

# Update config to use npx
ENV MCP_CONFIG_JSON='{"mcp":{"command":"npx","args":["-y","@modelcontextprotocol/server-weather"]}}'
```

**Important**: Ensure your binary is compiled for Linux x86_64, even if building on macOS/Windows.

### 2. Create Dockerfile

See [examples/deployment/Dockerfile](../examples/deployment/Dockerfile) for the complete Dockerfile.

The Dockerfile is already configured in the repository root and includes:
- Multi-stage build with Go 1.25
- Embedded MCP server binary validation (must be >1MB)
- SQLite database (~20MB) for persistent storage
- Non-root user (mcpdiscord, uid 1000)
- Health check using `pgrep`
- Entrypoint script for process management

### 3. Build Image

### 3. Build Image

```bash
# Build image
docker build -t mcpdiscord:latest .

# Tag for registry (optional)
docker tag mcpdiscord:latest registry.example.com/mcpdiscord:latest
```

### 4. Create docker-compose.yml

See [examples/deployment/docker-compose.yml](../examples/deployment/docker-compose.yml) for the complete compose file.

### 5. Deploy with Docker Compose

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

### 6. Health Check

```bash
# Check container status
docker-compose ps

# Check logs for errors
docker-compose logs mcpdiscord | grep -i error

# Restart if needed
docker-compose restart mcpdiscord
```

## Cloud Deployments

### AWS ECS Fargate (Terraform + GitHub Actions)

This repository includes an opinionated AWS deployment path using:

- Terraform infrastructure in `infra/`
- Container image push + ECS rollout in `.github/workflows/deploy.yml`

#### 1. Build and Runtime Model

- The root `Dockerfile` builds a self-contained runtime image.
- Embedded assets include:
  - MCP server binary at `/app/mcp-server` (approximately 50MB)
  - SQLite DB file at `/app/data/bot.db` (approximately 20MB)
- SQLite data is stored in the container writable layer for this initial rollout.
  - Important: replacing tasks can reset writable-layer changes.

#### 2. Provision AWS Infrastructure with Terraform

```bash
cd infra
terraform init
terraform plan -var="discord_bot_token=YOUR_DISCORD_TOKEN"
terraform apply -var="discord_bot_token=YOUR_DISCORD_TOKEN"
```

Terraform creates:

- ECR repository
- ECS cluster, task definition, and service
- VPC + public subnets + internet gateway
- IAM roles for task execution and runtime
- CloudWatch log group

#### 3. Configure GitHub Secrets

Add the following repository secrets before enabling automated deployment:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_REGION`
- `ECR_REPOSITORY`
- `ECS_CLUSTER`
- `ECS_SERVICE`

`DISCORD_BOT_TOKEN` is provided to ECS at runtime (Terraform variable or task definition environment) and must not be baked into the image.

#### 4. CI/CD on Push to Main

On each push to `main`, `.github/workflows/deploy.yml`:

1. Configures AWS credentials
2. Logs in to ECR
3. Builds and pushes image tags (`<sha>`, `latest`)
4. Renders ECS task definition with the new image
5. Deploys and waits for ECS service stability

#### 5. Verify Deployment

```bash
# ECS service status
aws ecs describe-services \
  --cluster <your-cluster> \
  --services <your-service> \
  --query 'services[0].{status:status,running:runningCount,desired:desiredCount}'

# Tail application logs from CloudWatch
aws logs tail /ecs/mcpdiscord --follow
```

Verify in Discord that slash commands are available and tool calls succeed.

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
