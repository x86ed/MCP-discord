# MCP-Discord

A Discord bot that automatically wraps Model Context Protocol (MCP) server tools as Discord slash commands, providing a 1-to-1 bridge between MCP servers and Discord.

[![Test](https://github.com/yourusername/MCP-discord/workflows/Test/badge.svg)](https://github.com/yourusername/MCP-discord/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/MCP-discord)](https://goreportcard.com/report/github.com/yourusername/MCP-discord)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

MCP-Discord automatically discovers tools from any MCP server and registers them as Discord slash commands. When users invoke these commands in Discord, the bot translates the interaction into an MCP tool call and returns the response.

```
Discord Slash Command  →  MCP Tool Call  →  MCP Server
      ↓                                          ↓
Discord Response      ←  Bot Translation  ←  MCP Response
```

## Features

- 🔄 **Automatic Discovery**: Connects to an MCP server and auto-registers all available tools as slash commands
- 🔐 **Secure Configuration**: Environment variable interpolation keeps secrets out of config files
- 🧪 **Well Tested**: Comprehensive test suite with 95%+ code coverage requirement
- 📝 **Schema Translation**: Intelligent mapping between JSON Schema (MCP) and Discord command options
- ⚡ **Fast Development**: Guild-specific command registration for instant testing
- 🛠️ **Standard Go Layout**: Clean, maintainable codebase following Go best practices

## Quick Start

### Prerequisites

- Go 1.21 or higher
- A Discord bot token ([Create one here](https://discord.com/developers/applications))
- An MCP server (e.g., `@modelcontextprotocol/server-weather`)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/MCP-discord.git
cd MCP-discord

# Build the bot
go build -o mcpdiscord ./cmd/mcpdiscord

# Or install directly
go install ./cmd/mcpdiscord
```

### Basic Configuration

Create a `mcp-config.json` file:

```json
{
  "discord": {
    "token": "YOUR_DISCORD_BOT_TOKEN",
    "guildId": ""
  },
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"],
    "transport": "stdio"
  }
}
```

See [examples/basic/](examples/basic/) for a complete working example.

### Running the Bot

```bash
# Using the default config path (./mcp-config.json)
./mcpdiscord

# Or specify a config file
./mcpdiscord --config /path/to/config.json

# Or use environment variable
export MCP_CONFIG_PATH=/path/to/config.json
./mcpdiscord

# Show version
./mcpdiscord --version
```

## Using Slash Commands

Once the bot is running and connected to Discord, it will automatically register all MCP tools as slash commands with a **1-to-1 mapping** (e.g., MCP tool "list" → Discord command "/list").

### Command Discovery

The bot discovers available tools from your MCP server at startup and registers them as Discord slash commands. Each tool's name, description, and parameters are automatically translated.

### Parameter Descriptions

The bot reads parameter descriptions from your MCP server's tool schema and displays them as hints in Discord's slash command UI. To ensure users see helpful parameter descriptions:

**MCP Server Tool Schema Example:**
```json
{
  "name": "get_weather",
  "description": "Get current weather for a location",
  "inputSchema": {
    "type": "object",
    "properties": {
      "location": {
        "type": "string",
        "description": "City name or ZIP code"
      },
      "units": {
        "type": "string",
        "description": "Temperature units (celsius or fahrenheit)"
      }
    },
    "required": ["location"]
  }
}
```

**Best Practices:**
- Always include `description` fields for all parameters in your MCP tool schemas
- Keep descriptions concise (under 100 characters) to fit Discord's limits
- Be specific about expected formats (e.g., "Comma-separated list of IDs" for arrays)
- The bot will log warnings if parameters are missing descriptions to help you identify issues

### Parameter Formats

MCP tools use JSON Schema for parameter definitions. The bot translates these to Discord's option types:

| MCP Type | Discord Type | Format |
|----------|-------------|--------|
| `string` | String | Direct text input |
| `number`/`integer` | Number | Numeric input |
| `boolean` | Boolean | True/false toggle |
| `array` | String | **Comma-separated values** |
| `object` | String | **JSON object string** |

### Array Parameters (CSV Format)

For array parameters, provide comma-separated values. The bot will trim whitespace around items.

**Examples:**
```
/get_weather locations: Seattle, Portland, San Francisco
/list_users ids: user1, user2, user3
/process_items items: apple,banana,orange
```

**Notes:**
- Whitespace around commas is trimmed automatically
- Empty items are filtered out
- Single items work without commas: `Seattle`
- **Limitation**: Cannot include commas within items (no escaping support)

### Object Parameters (JSON Format)

For object parameters, provide a JSON string. The bot validates JSON syntax and provides helpful error messages.

**Examples:**
```
/configure settings: {"theme": "dark", "notifications": true}
/create_user data: {"name": "Alice", "age": 30, "email": "alice@example.com"}
/query filters: {"status": "active", "role": "admin"}
```

**Nested objects:**
```
/create_profile data: {"name": "Bob", "address": {"city": "Seattle", "state": "WA"}}
```

**Notes:**
- Must be valid JSON (use double quotes for strings)
- Can include nested objects and arrays
- Bot will show error with example if JSON is invalid

### Discord Limitations

The bot respects Discord's platform constraints:

| Limit | Value |
|-------|-------|
| Max commands per guild | 100 |
| Max options per command | 25 |
| Command description length | 100 characters |
| Option description length | 100 characters |
| Embed description length | 4096 characters |

Descriptions exceeding limits are automatically truncated with "..." suffix.

### Response Format

Responses are displayed as rich embeds:

**Success (Green):**
```
✓ tool_name
Result content here
```

**Error (Red):**
```
✗ tool_name
Error message here
```

JSON/structured responses are automatically formatted with code blocks.

## Configuration

### Configuration File Format

The bot uses JSON configuration files with support for environment variable interpolation:

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": "${DISCORD_GUILD_ID}"
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

### Configuration Sources Priority

1. **CLI Flag** (highest): `--config /path/to/config.json`
2. **Environment Variable**: `MCP_CONFIG_PATH=/path/to/config.json`
3. **Default Path** (lowest): `./mcp-config.json`

### Environment Variable Interpolation

Use `${VAR_NAME}` syntax in your config file to reference environment variables:

```bash
export DISCORD_BOT_TOKEN="your-token-here"
export WEATHER_API_KEY="your-api-key"
./mcpdiscord --config config.json
```

See [examples/with-env/](examples/with-env/) for a complete example.

## Usage Examples

### Connecting to Different MCP Servers

**Weather Server:**

```json
{
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"]
  }
}
```

**Custom Server:**

```json
{
  "mcp": {
    "command": "/path/to/your/mcp-server",
    "args": ["--port", "8080"]
  }
}
```

**Node.js Server:**

```json
{
  "mcp": {
    "command": "node",
    "args": ["./my-mcp-server.js"]
  }
}
```

### Development vs Production

**Development** (fast command updates):

```json
{
  "discord": {
    "guildId": "your-test-server-id"
  }
}
```

**Production** (global commands):

```json
{
  "discord": {
    "guildId": ""
  }
}
```

## Deployment

### Quick Comparison

| Method | Complexity | Best For | Auto-Restart | Scaling |
|--------|-----------|----------|--------------|---------|
| **Docker** | Low | Local/single-server | Yes | Manual |
| **AWS ECS** | Medium | Production cloud | Yes | Auto |
| **Systemd** | Low | Linux VPS | Yes | N/A |

### Docker Deployment (Recommended)

The easiest way to deploy MCP-Discord in production. The included Dockerfile builds a self-contained image with:

- **Embedded MCP server binary** (~50MB at `/app/mcp-server`)
- **SQLite database** (~20MB at `/app/data/bot.db`)
- **Multi-process management** (bot + MCP server)
- **Health checks** with automatic restart

**Quick Start:**

```bash
# Build the image
docker build -t mcpdiscord .

# Run with environment config
docker run -d \
  --name mcpdiscord \
  -e MCP_CONFIG_JSON='{"discord":{"token":"YOUR_TOKEN"},"mcp":{"command":"/app/mcp-server"}}' \
  mcpdiscord

# View logs
docker logs -f mcpdiscord
```

**Building with Your Go MCP Server:**

If your MCP server is also a Go program:

```bash
# Add your MCP server code to the repo
mkdir -p cmd/mcpserver
cp -r /path/to/your/mcp-server/* cmd/mcpserver/

# Build image (automatically detects and builds cmd/mcpserver)
docker build -t mcpdiscord .
```

The Dockerfile automatically builds any Go program in `cmd/mcpserver/` and embeds it at `/app/mcp-server`.

**Using docker-compose:**

```yaml
version: '3.8'
services:
  mcpdiscord:
    build: .
    environment:
      - MCP_CONFIG_JSON=${MCP_CONFIG_JSON}
      # Or use individual vars and mount config:
      - DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN}
    volumes:
      - ./mcp-config.json:/app/config.json:ro
    restart: unless-stopped
```

See [examples/deployment/](examples/deployment/) for complete Docker and docker-compose configurations.

### AWS ECS Deployment with Terraform

Fully automated cloud deployment with infrastructure-as-code and CI/CD:

**Infrastructure** (Terraform):
- ECR container registry with lifecycle policies
- ECS Fargate cluster (no server management)
- VPC with public subnets and internet gateway
- IAM roles for task execution
- CloudWatch logging (14-day retention)
- Task definition with environment variable injection

**CI/CD** (GitHub Actions):
- Automatic build and push to ECR on commits to `main`
- Rolling deployment to ECS with service stability checks
- Zero-downtime updates

**Prerequisites:**
- AWS account with appropriate IAM permissions
- Terraform >= 1.5 installed
- GitHub repository

**Step 1: Provision Infrastructure**

```bash
cd infra

# Initialize Terraform
terraform init

# Review planned changes
terraform plan \
  -var="discord_bot_token=YOUR_DISCORD_TOKEN" \
  -var="mcp_config_json=$(cat ../mcp-config.json)"

# Apply infrastructure
terraform apply \
  -var="discord_bot_token=YOUR_DISCORD_TOKEN" \
  -var="mcp_config_json=$(cat ../mcp-config.json)"

# Note the outputs
terraform output ecr_repository_url
terraform output ecs_cluster_arn
terraform output ecs_service_arn
```

**Step 2: Configure GitHub Secrets**

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

| Secret | Description | Example |
|--------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | AWS IAM access key | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS IAM secret key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `AWS_REGION` | AWS region | `us-east-1` |
| `ECR_REPOSITORY` | ECR repository name | `mcpdiscord` |
| `ECS_CLUSTER` | ECS cluster name | `mcpdiscord-cluster` |
| `ECS_SERVICE` | ECS service name | `mcpdiscord-service` |

**Step 3: Deploy**

```bash
# Push to main branch triggers automatic deployment
git push origin main

# Monitor deployment in GitHub Actions tab
# Or watch ECS service in AWS console
```

**Step 4: Verify**

```bash
# Check ECS service status
aws ecs describe-services \
  --cluster mcpdiscord-cluster \
  --services mcpdiscord-service \
  --query 'services[0].{status:status,running:runningCount,desired:desiredCount}'

# Tail CloudWatch logs
aws logs tail /ecs/mcpdiscord --follow
```

**Cost Estimate** (us-east-1):
- ECS Fargate (0.25 vCPU, 0.5GB): ~$5/month
- ECR storage: <$1/month
- Data transfer (minimal): <$1/month
- **Total**: ~$7/month

**Terraform Configuration:**

Key variables in `infra/variables.tf`:

```hcl
variable "aws_region" {
  default = "us-east-1"
}

variable "task_cpu" {
  default = "256"  # 0.25 vCPU
}

variable "task_memory" {
  default = "512"  # 512MB
}

variable "ephemeral_storage_gib" {
  default = 21  # 20GB for embedded assets + buffer
}

variable "discord_bot_token" {
  sensitive = true
}

variable "mcp_config_json" {
  description = "Complete MCP config as JSON string"
  default     = ""
}
```

See [infra/README.md](infra/README.md) for detailed Terraform documentation.

### Linux Systemd Deployment

Best for VPS or dedicated servers running Linux:

```bash
# Install binary
sudo cp mcpdiscord /usr/local/bin/
sudo chmod +x /usr/local/bin/mcpdiscord

# Create config directory
sudo mkdir -p /etc/mcpdiscord
sudo cp mcp-config.json /etc/mcpdiscord/

# Create systemd service
sudo cp examples/deployment/mcpdiscord.service /etc/systemd/system/

# Start and enable
sudo systemctl daemon-reload
sudo systemctl enable --now mcpdiscord

# Check status
sudo systemctl status mcpdiscord
sudo journalctl -u mcpdiscord -f
```

### Other Deployment Options

- **AWS EC2**: Install on EC2 instance, use Systemd for process management
- **Google Cloud Compute Engine**: Similar to EC2 deployment
- **DigitalOcean Droplet**: Use Systemd on Ubuntu/Debian droplet
- **Heroku**: Use `Procfile` with worker dyno
- **Railway**: Connect GitHub repo, set environment variables

### Configuration in Production

**Option 1: JSON Environment Variable** (Docker/Cloud-native)

```bash
export MCP_CONFIG_JSON='{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": ""
  },
  "mcp": {
    "command": "/app/mcp-server",
    "transport": "stdio"
  }
}'
```

**Option 2: Mounted Configuration File** (Traditional)

```bash
# Mount config as read-only volume
docker run -v /path/to/mcp-config.json:/app/config.json:ro mcpdiscord
```

**Option 3: Environment File** (Systemd)

```bash
# /etc/mcpdiscord/environment
DISCORD_BOT_TOKEN=your_token
WEATHER_API_KEY=your_key
```

### Security Best Practices

1. **Never commit secrets** to version control
2. **Use environment variables** for sensitive values
3. **Restrict IAM permissions** to minimum required (AWS)
4. **Enable CloudWatch logs** for audit trail
5. **Use readonly root filesystem** when possible (Docker)
6. **Run as non-root user** (Dockerfile uses uid 1000)
7. **Rotate tokens regularly** (Discord bot tokens)

### Monitoring and Health Checks

**Docker Health Check:**
```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD pgrep mcpdiscord || exit 1
```

**ECS Health Check:**
- Configured automatically in task definition
- Monitors process existence
- Restarts container on failure

**Logs:**
```bash
# Docker
docker logs -f mcpdiscord

# AWS CloudWatch
aws logs tail /ecs/mcpdiscord --follow

# Systemd
journalctl -u mcpdiscord -f
```

### Troubleshooting Deployments

**Container won't start:**
```bash
# Check logs for error messages
docker logs mcpdiscord

# Common issues:
# - Missing DISCORD_BOT_TOKEN
# - Invalid MCP_CONFIG_JSON syntax
# - MCP server binary not found or not executable
```

**ECS task failing:**
```bash
# Check task stopped reason
aws ecs describe-tasks --cluster mcpdiscord-cluster --tasks TASK_ID

# View CloudWatch logs
aws logs tail /ecs/mcpdiscord --since 10m

# Common issues:
# - Insufficient CPU/memory (increase task_cpu/task_memory)
# - ECR pull permissions (check IAM role)
# - Environment variables not set
```

**Bot connects but commands don't appear:**
- Wait up to 1 hour for global commands to propagate
- Use `guildId` in config for instant updates (development)
- Check bot has `applications.commands` scope in OAuth2 URL

For detailed deployment instructions and advanced configurations, see the [Deployment Guide](docs/deployment.md).

## Documentation

- [Architecture Overview](docs/architecture.md) - System design and component relationships
- [Configuration Guide](docs/configuration.md) - Detailed configuration reference
- [Deployment Guide](docs/deployment.md) - Production deployment instructions
- [Development Guide](docs/development.md) - Contributing and development setup

## Project Structure

```md
cmd/mcpdiscord/         # Application entry point
internal/
  bot/                  # Discord bot logic
  mcp/                  # MCP client implementation
  translator/           # Schema translation (MCP ↔ Discord)
  config/               # Configuration management
  testutil/             # Testing utilities
examples/               # Example configurations
docs/                   # Documentation
```

## Development

### Running Tests

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building

```bash
# Build for current platform
go build -o mcpdiscord ./cmd/mcpdiscord

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o mcpdiscord-linux-amd64 ./cmd/mcpdiscord
GOOS=darwin GOARCH=arm64 go build -o mcpdiscord-darwin-arm64 ./cmd/mcpdiscord
GOOS=windows GOARCH=amd64 go build -o mcpdiscord-windows-amd64.exe ./cmd/mcpdiscord
```

## Troubleshooting

### Missing Parameter Descriptions in Discord

**Symptom:** Slash command parameters show "No description" or generic hints like "(JSON object)" without context.

**Cause:** Your MCP server's tool schema doesn't include `description` fields for parameters.

**Solution:**
1. Check your MCP server logs for warnings:
   ```
   WARN Parameter missing description tool=your_tool parameter=param_name
   ```
2. Update your MCP server to include descriptions in the tool schema (see "Parameter Descriptions" section above)
3. Restart the bot to refresh command registrations

**Example Fix:**
```json
// Before (missing description)
{
  "properties": {
    "location": {
      "type": "string"
    }
  }
}

// After (with description)
{
  "properties": {
    "location": {
      "type": "string",
      "description": "City name or ZIP code"
    }
  }
}
```

### Commands Not Updating

**Symptom:** Changes to MCP tools don't appear in Discord.

**Cause:** Discord caches slash commands.

**Solution:**
- **Guild commands** (development): Updates are instant
- **Global commands** (production): Can take up to 1 hour to propagate
- Restart the bot to force re-registration

### Connection Errors

**Symptom:** Bot fails to connect to MCP server.

**Solution:**
1. Verify the MCP server command/path is correct in your config
2. Check MCP server logs for errors
3. Ensure the MCP server supports stdio transport
4. Test the MCP server independently:
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | your-mcp-command
   ```

## Contributing

Contributions are welcome! Please see [docs/development.md](docs/development.md) for guidelines.

### Code Quality Standards

- **Test Coverage**: Minimum 95% required (enforced in CI)
- **Linting**: All code must pass `golangci-lint`
- **Documentation**: All exported functions and types must be documented

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) for the MCP specification
- [discordgo](https://github.com/bwmarrin/discordgo) for Discord API interactions

## Support

- 🐛 [Report a Bug](https://github.com/yourusername/MCP-discord/issues/new?labels=bug)
- 💡 [Request a Feature](https://github.com/yourusername/MCP-discord/issues/new?labels=enhancement)
- 💬 [Discussions](https://github.com/yourusername/MCP-discord/discussions)
