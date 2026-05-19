## Why

The Discord bot currently has no mechanism to discover and expose MCP server tools as Discord slash commands. This prevents users from interacting with MCP functionality through Discord, making the bot non-functional for its core purpose: bridging MCP servers to Discord interfaces.

## What Changes

- Add automatic discovery of MCP tools from the configured server
- Implement 1-to-1 dynamic Discord slash command registration for each discovered MCP tool (e.g., MCP "list" → Discord "/list")
- Add argument type translation (CSV for list parameters, JSON strings for object parameters)
- Implement command execution pipeline that translates Discord interactions to MCP tool calls
- Add response formatting to display MCP tool results in Discord messages

## Capabilities

### New Capabilities

- `tool-discovery`: Discovering available tools from MCP servers via tool listing protocol
- `slash-command-registration`: Dynamically registering Discord slash commands based on discovered MCP tools
- `argument-translation`: Converting Discord slash command arguments to MCP tool input schema (CSV→arrays, JSON strings→objects)
- `command-execution`: Executing MCP tools when Discord slash commands are invoked and handling responses

### Modified Capabilities

None - this is new functionality building on the existing scaffold.

## Impact

- **New packages**: Core functionality in `internal/bot`, `internal/mcp`, and `internal/translator`
- **Discord API**: Uses discordgo for slash command registration and interaction handling
- **MCP Protocol**: Implements tool listing and tool execution from MCP specification
- **Testing**: Requires integration tests with mock MCP servers and Discord interactions
- **Configuration**: May need extension for tool filtering or customization options
