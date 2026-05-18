# Basic Configuration Example

This directory contains a basic configuration example for the MCP-Discord bot that connects to the weather MCP server.

## Configuration File

The `mcp-config.json` file demonstrates the minimal required configuration:

```json
{
  "discord": {
    "token": "YOUR_DISCORD_BOT_TOKEN_HERE",
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

### Fields

#### Discord Configuration

- **token** (required): Your Discord bot token from the [Discord Developer Portal](https://discord.com/developers/applications)
- **guildId** (optional): A specific guild/server ID for faster command registration during development. Leave empty for global commands.

#### MCP Configuration

- **command** (required): The executable to run the MCP server. Here we use `npx` to run an npm package without installing it.
- **args** (required): Command-line arguments passed to the MCP server command.
- **env** (optional): Environment variables to set for the MCP server process.
- **transport** (optional): Communication method with the MCP server. Defaults to `"stdio"`. Other options: `"sse"`, `"websocket"`.

## Usage

1. **Get a Discord bot token:**
   - Go to the [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application
   - Go to the "Bot" section and create a bot
   - Copy the token and replace `YOUR_DISCORD_BOT_TOKEN_HERE` in the config file

2. **Invite the bot to your server:**
   - In the Discord Developer Portal, go to OAuth2 > URL Generator
   - Select scopes: `bot`, `applications.commands`
   - Select permissions: `Send Messages`, `Use Slash Commands`
   - Use the generated URL to invite the bot

3. **Run the bot:**
   ```bash
   mcpdiscord --config examples/basic/mcp-config.json
   ```

4. **Test the weather command:**
   In Discord, type `/` and you should see commands from the weather server automatically registered as slash commands.

## Notes

- This example uses the public weather MCP server which doesn't require API keys
- The bot will automatically discover all tools from the MCP server and create Discord slash commands for them
- Commands are registered globally by default (can take up to 1 hour to propagate)
- For faster testing during development, set `guildId` to a specific server ID
