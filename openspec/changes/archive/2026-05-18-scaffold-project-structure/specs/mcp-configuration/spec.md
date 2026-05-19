## ADDED Requirements

### Requirement: JSON configuration file format
The bot SHALL use a JSON configuration file to specify MCP server connection details.

#### Scenario: Configuration file location
- **WHEN** bot starts up
- **THEN** it SHALL look for a configuration file at a path specified via environment variable or flag
- **THEN** default path SHALL be `./mcp-config.json`

#### Scenario: Required configuration fields
- **WHEN** parsing the configuration file
- **THEN** it SHALL include a `discord` section with bot token
- **THEN** it SHALL include an `mcp` section with server connection details

### Requirement: MCP server connection specification
The configuration SHALL specify how to launch and connect to an MCP server.

#### Scenario: Command-based server launch
- **WHEN** MCP server needs to be launched as subprocess
- **THEN** configuration SHALL include `command` field with executable path
- **THEN** configuration SHALL include `args` array with command arguments
- **THEN** configuration SHALL support optional `env` object for environment variables

#### Scenario: Connection type specification
- **WHEN** configuring MCP connection
- **THEN** configuration SHALL specify connection type (stdio, SSE, or WebSocket)
- **THEN** default connection type SHALL be stdio

### Requirement: Discord bot configuration
The configuration SHALL include Discord-specific settings.

#### Scenario: Bot token configuration
- **WHEN** connecting to Discord
- **THEN** configuration SHALL include `discord.token` field
- **THEN** token SHALL be loadable from environment variable reference

#### Scenario: Optional guild-specific commands
- **WHEN** faster command updates are needed during development
- **THEN** configuration SHALL support optional `discord.guildId` field
- **THEN** when guildId is present, commands SHALL be registered to that guild only

### Requirement: Configuration validation
The bot SHALL validate configuration on startup and provide clear error messages.

#### Scenario: Missing required fields
- **WHEN** configuration is missing required fields
- **THEN** bot SHALL fail to start with descriptive error message
- **THEN** error message SHALL indicate which fields are missing

#### Scenario: Invalid command path
- **WHEN** MCP server command path does not exist
- **THEN** bot SHALL fail to start with clear error about invalid path
- **THEN** error message SHALL show the attempted path

### Requirement: Example configuration provided
The project SHALL include example configuration files for common use cases.

#### Scenario: Weather server example
- **WHEN** user wants to connect to weather MCP server
- **THEN** example configuration SHALL be available
- **THEN** example SHALL include comments explaining each field

#### Scenario: Environment variable usage
- **WHEN** user wants to keep secrets out of config file
- **THEN** example SHALL demonstrate environment variable references
- **THEN** format SHALL support `${ENV_VAR_NAME}` syntax
