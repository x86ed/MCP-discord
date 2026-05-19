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

## Documentation

- [Architecture Overview](docs/architecture.md) - System design and component relationships
- [Configuration Guide](docs/configuration.md) - Detailed configuration reference
- [Deployment Guide](docs/deployment.md) - Production deployment instructions
- [Development Guide](docs/development.md) - Contributing and development setup

## Project Structure

```
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
