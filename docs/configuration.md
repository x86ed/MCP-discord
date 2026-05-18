# Configuration Guide

Complete reference for configuring the MCP-Discord bot.

## Table of Contents

- [Configuration File Format](#configuration-file-format)
- [Configuration Sources](#configuration-sources)
- [Discord Configuration](#discord-configuration)
- [MCP Configuration](#mcp-configuration)
- [Environment Variable Interpolation](#environment-variable-interpolation)
- [Validation Rules](#validation-rules)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## Configuration File Format

The bot uses JSON configuration files. The basic structure:

```json
{
  "discord": {
    "token": "string",
    "guildId": "string (optional)"
  },
  "mcp": {
    "command": "string",
    "args": ["string", "..."],
    "env": {"KEY": "value"},
    "transport": "string (optional)"
  }
}
```

## Configuration Sources

The bot loads configuration from multiple sources with the following priority order:

1. **CLI Flag** (highest priority)

   ```bash
   mcpdiscord --config /path/to/config.json
   ```

2. **Environment Variable**

   ```bash
   export MCP_CONFIG_PATH=/path/to/config.json
   mcpdiscord
   ```

3. **Default Path** (lowest priority)

   ```bash
   # Looks for ./mcp-config.json in current directory
   mcpdiscord
   ```

## Discord Configuration

### `discord.token` (required)

Your Discord bot token from the [Discord Developer Portal](https://discord.com/developers/applications).

**Type:** `string`  
**Required:** Yes  
**Environment Variable:** Recommended  

**Example:**

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}"
  }
}
```

**How to Get:**

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create or select your application
3. Go to "Bot" section
4. Click "Reset Token" or "Copy" to get your token

**Security:** Never commit tokens to version control. Use environment variables.

### `discord.guildId` (optional)

A specific Discord guild (server) ID to register commands to. When set, commands are registered only to that guild, which enables instant command updates during development.

**Type:** `string`  
**Required:** No  
**Default:** `""` (global commands)  

**Global vs Guild Commands:**

| Aspect | Global Commands | Guild Commands |
|--------|----------------|----------------|
| Registration Time | Up to 1 hour | Instant |
| Scope | All servers with bot | Single server only |
| Use Case | Production | Development/Testing |

**Example:**

```json
{
  "discord": {
    "guildId": "123456789012345678"
  }
}
```

**How to Get Guild ID:**

1. Enable Developer Mode in Discord (User Settings > Advanced > Developer Mode)
2. Right-click on your server name
3. Click "Copy Server ID"

## MCP Configuration

### `mcp.command` (required)

The executable command to launch the MCP server.

**Type:** `string`  
**Required:** Yes  

**Examples:**

**NPM package (recommended for packages):**

```json
{
  "mcp": {
    "command": "npx"
  }
}
```

**Node.js script:**

```json
{
  "mcp": {
    "command": "node"
  }
}
```

**Absolute path:**

```json
{
  "mcp": {
    "command": "/usr/local/bin/my-mcp-server"
  }
}
```

**Python script:**

```json
{
  "mcp": {
    "command": "python3"
  }
}
```

### `mcp.args` (optional)

Command-line arguments passed to the MCP server command.

**Type:** `array of strings`  
**Required:** No  
**Default:** `[]`  

**Examples:**

**NPM package with version:**

```json
{
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather@1.0.0"]
  }
}
```

**Script with arguments:**

```json
{
  "mcp": {
    "command": "node",
    "args": ["./server.js", "--port", "8080", "--verbose"]
  }
}
```

**Python module:**

```json
{
  "mcp": {
    "command": "python3",
    "args": ["-m", "my_mcp_server", "--config", "server-config.json"]
  }
}
```

### `mcp.env` (optional)

Environment variables to set for the MCP server process. These are passed to the subprocess.

**Type:** `object (key-value pairs)`  
**Required:** No  
**Default:** `{}`  

**Example:**

```json
{
  "mcp": {
    "env": {
      "API_KEY": "${WEATHER_API_KEY}",
      "LOG_LEVEL": "debug",
      "CACHE_DIR": "/tmp/mcp-cache"
    }
  }
}
```

**Use Cases:**

- API keys and secrets
- Configuration paths
- Debug settings
- Feature flags

**Note:** The MCP server also inherits environment variables from the bot process.

### `mcp.transport` (optional)

The communication protocol to use with the MCP server.

**Type:** `string`  
**Required:** No  
**Default:** `"stdio"`  
**Valid Values:** `"stdio"`, `"sse"`, `"websocket"`  

**Transport Types:**

#### stdio (Standard Input/Output)

```json
{
  "mcp": {
    "transport": "stdio"
  }
}
```

**Characteristics:**

- Default and most common
- Subprocess communication via stdin/stdout
- Works with any MCP server that supports stdio
- Simple and reliable

**Use When:**

- Running local MCP servers
- Server is a subprocess
- No special networking requirements

#### sse (Server-Sent Events)

```json
{
  "mcp": {
    "transport": "sse",
    "command": "http://localhost:8080/mcp"
  }
}
```

**Characteristics:**

- HTTP-based streaming
- Server must be running separately
- One-way server → client events

**Use When:**

- MCP server is a separate HTTP service
- Server is remote or containerized
- Want to decouple bot and server lifecycles

#### websocket

```json
{
  "mcp": {
    "transport": "websocket",
    "command": "ws://localhost:8080/mcp"
  }
}
```

**Characteristics:**

- Bidirectional communication
- Server must be running separately
- Real-time messaging

**Use When:**

- Need bidirectional communication
- Server requires persistent connection
- Server is remote or scalable

## Environment Variable Interpolation

The bot supports `${VAR_NAME}` syntax for environment variable interpolation anywhere in the configuration file.

### Syntax

```json
{
  "field": "${ENVIRONMENT_VARIABLE_NAME}"
}
```

### Examples

**Single variable:**

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}"
  }
}
```

**Multiple variables in one string:**

```json
{
  "mcp": {
    "command": "http://${MCP_HOST}:${MCP_PORT}/api"
  }
}
```

**Nested in objects:**

```json
{
  "mcp": {
    "env": {
      "DATABASE_URL": "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
    }
  }
}
```

### Behavior

- **Variable exists**: Replaced with the environment variable value
- **Variable doesn't exist**: Kept as literal `${VAR_NAME}` (triggers validation error)
- **Case sensitive**: `${TOKEN}` ≠ `${token}`
- **No defaults**: Must be set before starting the bot

### Setting Environment Variables

**Linux/macOS:**

```bash
export DISCORD_BOT_TOKEN="your-token"
export WEATHER_API_KEY="your-key"
```

**Windows (PowerShell):**

```powershell
$env:DISCORD_BOT_TOKEN="your-token"
$env:WEATHER_API_KEY="your-key"
```

**Using .env file (with direnv):**

```bash
# .env
DISCORD_BOT_TOKEN=your-token
WEATHER_API_KEY=your-key

# Load with direnv
echo "dotenv" > .envrc
direnv allow
```

## Validation Rules

The bot validates configuration on startup. Common validation errors:

### Required Fields

❌ **Error:** `discord.token is required`

```json
{
  "discord": {
    "token": ""  // Empty token
  }
}
```

✅ **Fix:**

```json
{
  "discord": {
    "token": "your-actual-token"
  }
}
```

### Unresolved Environment Variables

❌ **Error:** `discord.token contains unresolved environment variable`

```json
{
  "discord": {
    "token": "${DISCORD_TOKEN}"  // Variable not set
  }
}
```

✅ **Fix:**

```bash
export DISCORD_TOKEN="your-token"
```

### Invalid Transport Type

❌ **Error:** `mcp.transport must be one of: stdio, sse, websocket`

```json
{
  "mcp": {
    "transport": "http"  // Invalid
  }
}
```

✅ **Fix:**

```json
{
  "mcp": {
    "transport": "stdio"
  }
}
```

## Examples

### Minimal Configuration

```json
{
  "discord": {
    "token": "YOUR_DISCORD_BOT_TOKEN_HERE"
  },
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"]
  }
}
```

### Development Configuration

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": "123456789012345678"
  },
  "mcp": {
    "command": "node",
    "args": ["./my-dev-server.js", "--verbose"],
    "env": {
      "LOG_LEVEL": "debug",
      "DEV_MODE": "true"
    },
    "transport": "stdio"
  }
}
```

### Production Configuration

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": ""
  },
  "mcp": {
    "command": "/opt/mcp/server",
    "args": ["--config", "/etc/mcp/prod-config.json"],
    "env": {
      "API_KEY": "${MCP_API_KEY}",
      "LOG_LEVEL": "info",
      "ENVIRONMENT": "production"
    },
    "transport": "stdio"
  }
}
```

### Remote MCP Server (SSE)

```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}"
  },
  "mcp": {
    "command": "https://mcp-server.example.com/stream",
    "transport": "sse"
  }
}
```

## Troubleshooting

### Bot Won't Start

**Symptom:** Bot exits immediately with error  
**Common Causes:**

- Missing or invalid configuration file
- Unresolved environment variables
- Invalid Discord token
- Missing required fields

**Solutions:**

1. Check file exists: `ls -l mcp-config.json`
2. Validate JSON syntax: `jq . mcp-config.json`
3. Verify environment variables: `echo $DISCORD_BOT_TOKEN`
4. Check logs for specific error messages

### Commands Not Appearing

**Symptom:** Slash commands don't show up in Discord  
**Causes:**

- Global commands take up to 1 hour to register
- Bot lacks proper permissions
- Bot not invited with `applications.commands` scope

**Solutions:**

1. Use `guildId` for instant registration during development
2. Re-invite bot with correct OAuth2 scopes: `bot`, `applications.commands`
3. Check bot has required permissions in Discord server

### Environment Variables Not Working

**Symptom:** Error about unresolved variables  
**Causes:**

- Variable not exported
- Typo in variable name (case sensitive)
- Variable set after bot started

**Solutions:**

1. Export before running: `export VAR_NAME=value`
2. Check exact name: `echo $VAR_NAME`
3. Restart bot after setting variables

### MCP Server Connection Fails

**Symptom:** Bot starts but tools don't register  
**Causes:**

- MCP server command not found
- Wrong arguments
- Server crashes on startup
- Transport mismatch

**Solutions:**

1. Test command manually: `npx -y @modelcontextprotocol/server-weather`
2. Check logs for MCP server errors
3. Verify server supports specified transport type
4. Check network connectivity (for SSE/WebSocket)

### Configuration Validation Fails

**Symptom:** Specific validation error on startup  
**Solution:** Read the error message carefully - it tells you exactly what's wrong and which field needs attention.

---

For more examples, see the [examples/](../examples/) directory.
