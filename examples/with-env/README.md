# Environment Variable Configuration Example

This directory contains an example configuration that demonstrates using environment variables to keep secrets out of configuration files.

## Why Use Environment Variables?

- **Security**: Keep sensitive tokens and API keys out of version control
- **Flexibility**: Different values for development, staging, and production
- **Deployment**: Works well with container orchestration (Docker, Kubernetes)
- **CI/CD**: Integrate with secret management systems

## Configuration File

The `mcp-config.json` uses the `${VARIABLE_NAME}` syntax for environment variable interpolation:

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
      "WEATHER_API_KEY": "${WEATHER_API_KEY}",
      "LOG_LEVEL": "${LOG_LEVEL}"
    },
    "transport": "stdio"
  }
}
```

## Setup

1. **Copy the example environment file:**
   ```bash
   cp examples/with-env/.env.example .env
   ```

2. **Edit `.env` with your actual values:**
   - Replace `your_discord_bot_token_here` with your Discord bot token
   - Add your API keys and other sensitive configuration
   - Customize log levels and other settings

3. **Load environment variables:**
   
   **Option A: Using direnv (recommended)**
   ```bash
   # Install direnv: https://direnv.net/
   echo "dotenv" > .envrc
   direnv allow
   ```

   **Option B: Manual loading**
   ```bash
   export $(cat .env | xargs)
   ```

   **Option C: Docker**
   ```bash
   docker run --env-file .env mcpdiscord
   ```

4. **Run the bot:**
   ```bash
   mcpdiscord --config examples/with-env/mcp-config.json
   ```

## Environment Variable Reference

### Required Variables

- **DISCORD_BOT_TOKEN**: Your Discord bot token from the Developer Portal

### Optional Variables

- **DISCORD_GUILD_ID**: Specific guild for command registration (faster updates during development)
- **WEATHER_API_KEY**: API key for weather services (if your MCP server requires it)
- **LOG_LEVEL**: Logging verbosity (`debug`, `info`, `warn`, `error`)

### MCP Server Variables

Variables in the `mcp.env` section are passed to the MCP server process. The exact variables needed depend on which MCP server you're using. Check your MCP server's documentation for requirements.

## Security Best Practices

1. **Never commit `.env` files to version control**
   - Add `.env` to your `.gitignore`
   - Only commit `.env.example` with placeholder values

2. **Use different `.env` files for each environment:**
   ```bash
   .env.development
   .env.staging
   .env.production
   ```

3. **In production, use secret management:**
   - Docker Secrets
   - Kubernetes Secrets
   - AWS Secrets Manager / Azure Key Vault
   - HashiCorp Vault

4. **Rotate tokens and keys regularly**

## Troubleshooting

### "unresolved environment variable" error

If you see validation errors about unresolved environment variables, make sure:
- The environment variable is set: `echo $DISCORD_BOT_TOKEN`
- The variable name matches exactly (case-sensitive)
- You've sourced the `.env` file before running the bot

### Variables not interpolating

The bot interpolates variables when loading the config file. If interpolation isn't working:
- Check the syntax: must be `${VAR_NAME}` (not `$VAR_NAME` or `%VAR_NAME%`)
- Ensure the variable is set before starting the bot
- Try running with explicit export: `export DISCORD_BOT_TOKEN=...`
