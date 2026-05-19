## ADDED Requirements

### Requirement: Translate string arguments

The system SHALL pass string arguments from Discord slash commands directly to MCP tool calls without modification.

#### Scenario: Simple string argument
- **WHEN** a user provides a string value for a parameter
- **THEN** the system passes the string unchanged to the MCP tool

#### Scenario: Empty string argument
- **WHEN** a user provides an empty string for an optional parameter
- **THEN** the system includes the empty string in the tool call

### Requirement: Translate number arguments

The system SHALL convert Discord number options to appropriate numeric types for MCP tool calls.

#### Scenario: Integer argument
- **WHEN** a user provides a number without decimal places
- **THEN** the system converts it to an integer

#### Scenario: Float argument
- **WHEN** a user provides a number with decimal places
- **THEN** the system preserves the decimal value as a float

### Requirement: Translate boolean arguments

The system SHALL pass Discord boolean options as boolean values to MCP tool calls.

#### Scenario: True boolean
- **WHEN** a user selects true for a boolean option
- **THEN** the system passes true to the MCP tool

#### Scenario: False boolean
- **WHEN** a user selects false for a boolean option
- **THEN** the system passes false to the MCP tool

### Requirement: Translate array arguments from CSV

The system SHALL parse comma-separated values into arrays for MCP tool calls.

#### Scenario: Simple CSV list
- **WHEN** a user provides "item1,item2,item3" for an array parameter
- **THEN** the system converts it to ["item1", "item2", "item3"]

#### Scenario: CSV with whitespace
- **WHEN** a user provides "item1, item2 , item3" with spaces
- **THEN** the system trims whitespace and converts to ["item1", "item2", "item3"]

#### Scenario: Single item array
- **WHEN** a user provides "singleitem" for an array parameter
- **THEN** the system converts it to ["singleitem"]

#### Scenario: Empty array
- **WHEN** a user provides an empty string for an optional array parameter
- **THEN** the system passes an empty array []

#### Scenario: CSV with commas in values
- **WHEN** values themselves contain commas
- **THEN** the system splits on all commas (limitation: no escaping support in v1)

### Requirement: Translate object arguments from JSON

The system SHALL parse JSON strings into objects for MCP tool calls.

#### Scenario: Valid JSON object
- **WHEN** a user provides '{"key": "value", "number": 42}' for an object parameter
- **THEN** the system parses it to a JSON object with those properties

#### Scenario: Nested JSON object
- **WHEN** a user provides '{"outer": {"inner": "value"}}' 
- **THEN** the system correctly parses the nested structure

#### Scenario: JSON array
- **WHEN** a user provides '[1, 2, 3]' for an array of objects parameter
- **THEN** the system parses it as a JSON array

#### Scenario: Invalid JSON
- **WHEN** a user provides malformed JSON like '{key: value}' (missing quotes)
- **THEN** the system returns an error message to the user via ephemeral message

#### Scenario: JSON parsing error details
- **WHEN** JSON parsing fails
- **THEN** the system provides a user-friendly error message with an example format

### Requirement: Handle missing optional parameters

The system SHALL omit optional parameters that are not provided by the user from the MCP tool call.

#### Scenario: Optional parameter not provided
- **WHEN** a user invokes a command without providing an optional parameter
- **THEN** the system does not include that parameter in the MCP tool call arguments

#### Scenario: All optional parameters omitted
- **WHEN** all parameters are optional and none are provided
- **THEN** the system calls the MCP tool with an empty arguments object

### Requirement: Validate required parameters

The system MUST ensure all required parameters are present before calling the MCP tool.

#### Scenario: Required parameter missing
- **WHEN** a required parameter is not provided (Discord should prevent this)
- **THEN** the system logs an error and returns a friendly message to the user

#### Scenario: All required parameters present
- **WHEN** all required parameters are provided
- **THEN** the system proceeds with the MCP tool call

### Requirement: Type validation

The system SHALL validate that parsed values match expected types before calling MCP tools.

#### Scenario: Number parameter receives non-numeric input
- **WHEN** Discord sends a non-numeric value for a number parameter (edge case)
- **THEN** the system returns a validation error to the user

#### Scenario: Array parameter receives invalid CSV
- **WHEN** CSV parsing produces unexpected results
- **THEN** the system logs details and returns a helpful error message
