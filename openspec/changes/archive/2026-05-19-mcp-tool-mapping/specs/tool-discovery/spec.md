## ADDED Requirements

### Requirement: Discover tools from MCP server

The system SHALL connect to the configured MCP server and retrieve the complete list of available tools using the MCP tool listing protocol.

#### Scenario: Successful tool discovery at startup
- **WHEN** the bot starts with valid MCP server configuration
- **THEN** the system connects to the MCP server and retrieves all available tools with their schemas

#### Scenario: Handle empty tool list
- **WHEN** the MCP server returns an empty tool list
- **THEN** the system logs a warning and continues without registering any slash commands

#### Scenario: Connection failure during discovery
- **WHEN** the MCP server is unreachable or returns an error during tool listing
- **THEN** the system logs the error and fails bot startup with a descriptive message

### Requirement: Re-discover tools on reconnection

The system SHALL re-discover tools when the MCP server connection is lost and re-established.

#### Scenario: MCP server restarts during bot runtime
- **WHEN** the MCP server connection drops and reconnects
- **THEN** the system re-discovers tools and updates Discord slash command registrations

#### Scenario: Tool list changes after reconnection
- **WHEN** tools are added or removed in the MCP server between connections
- **THEN** the system registers new commands and removes old commands from Discord

### Requirement: Parse tool schemas

The system SHALL parse each discovered tool's input schema to extract parameter names, types, descriptions, and required/optional status.

#### Scenario: Tool with multiple parameter types
- **WHEN** a tool has string, number, boolean, array, and object parameters
- **THEN** the system correctly identifies each parameter's type and constraints

#### Scenario: Tool with required and optional parameters
- **WHEN** a tool schema specifies some parameters as required
- **THEN** the system distinguishes required from optional parameters for slash command registration

#### Scenario: Tool with missing or invalid schema
- **WHEN** a tool has no input schema or malformed schema
- **THEN** the system logs a warning and skips that tool, continuing with others

### Requirement: Validate tool metadata

The system SHALL validate that each tool has a name and description before attempting registration.

#### Scenario: Tool missing name
- **WHEN** a tool is returned without a name field
- **THEN** the system logs an error and skips that tool

#### Scenario: Tool with excessively long description
- **WHEN** a tool description exceeds 100 characters
- **THEN** the system truncates the description to 97 characters and appends "..."
