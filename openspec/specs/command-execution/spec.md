## ADDED Requirements

### Requirement: Execute MCP tools on slash command invocation

The system SHALL call the corresponding MCP tool when a Discord slash command is invoked.

#### Scenario: Successful tool execution
- **WHEN** a user invokes a registered slash command with valid arguments
- **THEN** the system calls the MCP tool with translated arguments and returns the result

#### Scenario: Tool execution with no parameters
- **WHEN** a user invokes a command for a tool that requires no parameters
- **THEN** the system calls the MCP tool with an empty arguments object

#### Scenario: Tool execution with all parameter types
- **WHEN** a user provides string, number, boolean, array (CSV), and object (JSON) parameters
- **THEN** the system correctly translates and passes all parameters to the MCP tool

### Requirement: Handle MCP tool execution errors

The system SHALL gracefully handle errors returned by MCP tools and communicate them to users.

#### Scenario: Tool returns error response
- **WHEN** an MCP tool call fails with an error message
- **THEN** the system sends a Discord message with the error in a red embed

#### Scenario: Tool execution timeout
- **WHEN** an MCP tool call does not respond within 30 seconds
- **THEN** the system cancels the request and notifies the user of the timeout

#### Scenario: MCP server connection lost during execution
- **WHEN** the MCP server disconnects while executing a tool
- **THEN** the system returns an error message and attempts to reconnect

### Requirement: Format tool results for Discord

The system SHALL format MCP tool results as Discord embeds for readability.

#### Scenario: Successful result
- **WHEN** a tool returns a successful result
- **THEN** the system creates a green embed with the tool name as title and result as description

#### Scenario: Result exceeds Discord limit
- **WHEN** a tool result is longer than 4096 characters
- **THEN** the system truncates the result to 4093 characters and appends "..."

#### Scenario: Empty result
- **WHEN** a tool returns an empty or null result
- **THEN** the system displays a message indicating successful execution with no output

#### Scenario: Structured result data
- **WHEN** a tool returns JSON object or array data
- **THEN** the system formats it as code block markdown for readability

### Requirement: Provide immediate feedback

The system SHALL acknowledge slash command invocations immediately to prevent Discord timeouts.

#### Scenario: Command acknowledged within 3 seconds
- **WHEN** a user invokes a slash command
- **THEN** the system sends a "thinking" state or initial response within 3 seconds

#### Scenario: Long-running tool execution
- **WHEN** a tool takes longer than 3 seconds to execute
- **THEN** the system shows a "processing" message and updates it when the result is ready

### Requirement: Handle concurrent command invocations

The system SHALL support multiple users invoking commands simultaneously.

#### Scenario: Multiple commands from different users
- **WHEN** multiple users invoke different commands at the same time
- **THEN** the system executes all commands concurrently without interference

#### Scenario: Same tool invoked multiple times
- **WHEN** multiple users invoke the same tool with different arguments
- **THEN** each invocation executes independently with correct argument isolation

### Requirement: Log command execution

The system SHALL log all command invocations and results for debugging and monitoring.

#### Scenario: Log successful execution
- **WHEN** a tool executes successfully
- **THEN** the system logs the command name, user, arguments, and execution time

#### Scenario: Log failed execution
- **WHEN** a tool execution fails
- **THEN** the system logs the error details, user, command, and arguments

#### Scenario: Log argument translation errors
- **WHEN** argument parsing fails
- **THEN** the system logs the raw input and parsing error details

### Requirement: Error messages are user-friendly

The system SHALL provide clear, actionable error messages to users when execution fails.

#### Scenario: Invalid JSON argument
- **WHEN** a user provides malformed JSON for an object parameter
- **THEN** the system sends an ephemeral message: "Invalid JSON format. Example: {"key": "value"}"

#### Scenario: Tool not found
- **WHEN** a slash command is invoked but the tool is no longer available
- **THEN** the system responds: "This command is no longer available. The bot may need to restart."

#### Scenario: Permission denied by MCP server
- **WHEN** the MCP server rejects a tool call due to permissions
- **THEN** the system displays the MCP error message without exposing internal details

### Requirement: Security - no injection vulnerabilities

The system MUST sanitize all user input to prevent command injection or code execution vulnerabilities.

#### Scenario: Attempt to inject shell commands
- **WHEN** a user provides input like "; rm -rf /" in a parameter
- **THEN** the system passes it as a literal string to the MCP tool without interpretation

#### Scenario: Attempt SQL injection patterns
- **WHEN** a user provides input like "'; DROP TABLE users; --"
- **THEN** the system passes it as a string value to the MCP tool, with no special handling

#### Scenario: JSON with malicious content
- **WHEN** a user provides JSON with executable code patterns
- **THEN** the system parses it as data only, with no code evaluation
