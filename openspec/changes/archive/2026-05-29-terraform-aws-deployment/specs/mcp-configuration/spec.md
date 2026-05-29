## MODIFIED Requirements

### Requirement: JSON configuration file format
The bot SHALL use a JSON configuration file OR a `MCP_CONFIG_JSON` environment variable to specify MCP server connection details.

#### Scenario: Configuration file location
- **WHEN** bot starts up
- **THEN** it SHALL look for a configuration file at a path specified via environment variable or flag
- **THEN** default path SHALL be `./mcp-config.json`

#### Scenario: Required configuration fields
- **WHEN** parsing the configuration
- **THEN** it SHALL include a `discord` section with bot token
- **THEN** it SHALL include an `mcp` section with server connection details

#### Scenario: Environment variable config injection
- **WHEN** the `MCP_CONFIG_JSON` environment variable is set to a non-empty string
- **THEN** the bot SHALL parse that string as JSON configuration instead of reading a file
- **THEN** `MCP_CONFIG_JSON` SHALL take precedence over any config file path
- **THEN** the bot SHALL log at info level that config was loaded from environment variable
