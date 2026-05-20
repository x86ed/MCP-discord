## MODIFIED Requirements

### Requirement: Map tool parameters to command options

The system SHALL create a Discord slash command option for each parameter in the tool's input schema, using the parameter's description from the MCP tool schema as the Discord option description.

#### Scenario: Map string parameter with description
- **WHEN** a tool parameter has type "string" and a description in the MCP schema
- **THEN** the system creates a String option with the parameter name and the description from the MCP schema

#### Scenario: Map string parameter without description
- **WHEN** a tool parameter has type "string" but no description in the MCP schema
- **THEN** the system creates a String option with a default description "No description"

#### Scenario: Map number parameter with description
- **WHEN** a tool parameter has type "number" and a description
- **THEN** the system creates a Number option with the description from the MCP schema

#### Scenario: Map boolean parameter with description
- **WHEN** a tool parameter has type "boolean" and a description
- **THEN** the system creates a Boolean option with the description from the MCP schema

#### Scenario: Map array parameter with type hint
- **WHEN** a tool parameter has type "array" and a description
- **THEN** the system creates a String option with the description followed by " (Comma-separated list)"

#### Scenario: Map array parameter without description
- **WHEN** a tool parameter has type "array" but no description
- **THEN** the system creates a String option with description "(Comma-separated list)"

#### Scenario: Map object parameter with type hint
- **WHEN** a tool parameter has type "object" and a description
- **THEN** the system creates a String option with the description followed by " (JSON object)"

#### Scenario: Map object parameter without description
- **WHEN** a tool parameter has type "object" but no description
- **THEN** the system creates a String option with description "(JSON object)"

#### Scenario: Empty string description
- **WHEN** a tool parameter has an empty string as its description
- **THEN** the system treats it as missing and uses the default "No description"

## ADDED Requirements

### Requirement: Log missing parameter descriptions

The system SHALL log a warning when registering commands with parameters that lack descriptions to help diagnose MCP server configuration issues.

#### Scenario: Parameter missing description during registration
- **WHEN** a parameter in the tool schema has no description or an empty description
- **THEN** the system logs a warning with the tool name and parameter name

#### Scenario: Multiple parameters missing descriptions
- **WHEN** multiple parameters lack descriptions in a single tool
- **THEN** the system logs a warning for each missing description
