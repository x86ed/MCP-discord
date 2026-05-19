## ADDED Requirements

### Requirement: Register slash commands for discovered tools

The system SHALL register a Discord slash command for each discovered MCP tool using the Discord API in a 1-to-1 mapping.

#### Scenario: Register command for simple tool
- **WHEN** an MCP tool with name "weather" and description "Get weather data" is discovered
- **THEN** the system registers Discord slash command "/weather" with the tool's description

#### Scenario: Direct mapping for list command
- **WHEN** an MCP tool named "list" is discovered
- **THEN** the system registers Discord slash command "/list"

#### Scenario: Sanitize tool names for Discord
- **WHEN** a tool name contains spaces or special characters (e.g., "Get User Data")
- **THEN** the system converts the name to lowercase kebab-case (e.g., "get-user-data")

#### Scenario: Handle command name collisions
- **WHEN** two tools sanitize to the same command name
- **THEN** the system logs an error with both tool names and fails registration

### Requirement: Map tool parameters to command options

The system SHALL create a Discord slash command option for each parameter in the tool's input schema.

#### Scenario: Map string parameter
- **WHEN** a tool parameter has type "string"
- **THEN** the system creates a String option with the parameter name and description

#### Scenario: Map number parameter
- **WHEN** a tool parameter has type "number"
- **THEN** the system creates a Number option

#### Scenario: Map boolean parameter
- **WHEN** a tool parameter has type "boolean"
- **THEN** the system creates a Boolean option

#### Scenario: Map array parameter
- **WHEN** a tool parameter has type "array"
- **THEN** the system creates a String option with description indicating "Comma-separated list"

#### Scenario: Map object parameter
- **WHEN** a tool parameter has type "object"
- **THEN** the system creates a String option with description indicating "JSON object"

### Requirement: Set option required status

The system SHALL mark Discord command options as required if the corresponding MCP parameter is required.

#### Scenario: Required parameter
- **WHEN** an MCP tool parameter is marked as required
- **THEN** the Discord slash command option is registered as required

#### Scenario: Optional parameter
- **WHEN** an MCP tool parameter is optional
- **THEN** the Discord slash command option is registered as optional

### Requirement: Enforce Discord slash command limits

The system SHALL validate that tool registration complies with Discord's slash command constraints.

#### Scenario: Exceed command limit per guild
- **WHEN** more than 100 tools are discovered
- **THEN** the system fails registration with an error message listing the tool count

#### Scenario: Parameter description too long
- **WHEN** a parameter description exceeds Discord's limit (100 characters)
- **THEN** the system truncates the description to fit

#### Scenario: Too many parameters
- **WHEN** a tool has more than 25 parameters (Discord's option limit)
- **THEN** the system logs a warning and skips that tool

### Requirement: Update commands on tool changes

The system SHALL update Discord slash command registrations when tools are added, removed, or modified.

#### Scenario: New tool discovered after reconnection
- **WHEN** the MCP server adds a new tool
- **THEN** the system registers the new slash command without affecting existing commands

#### Scenario: Tool removed after reconnection
- **WHEN** a tool is no longer available on the MCP server
- **THEN** the system deregisters the corresponding slash command from Discord

#### Scenario: Tool schema changes
- **WHEN** a tool's parameters or description change
- **THEN** the system updates the slash command registration with new metadata
