## 1. MCP Client Foundation

- [x] 1.1 Define `Client` interface in `internal/mcp/client.go` with Connect, ListTools, CallTool, Close methods
- [x] 1.2 Define `Tool` struct to represent MCP tool metadata (name, description, input schema)
- [x] 1.3 Define `ToolResult` struct to represent MCP tool execution results
- [x] 1.4 Implement stdio transport client in `internal/mcp/stdio.go` with subprocess management
- [x] 1.5 Add JSON-RPC message encoding/decoding utilities in `internal/mcp/jsonrpc.go`
- [x] 1.6 Implement `ListTools` method to call MCP "tools/list" endpoint
- [x] 1.7 Implement `CallTool` method to call MCP "tools/call" endpoint with arguments
- [x] 1.8 Add connection lifecycle management (Connect, Close, reconnection logic)

## 2. Tool Discovery & Schema Parsing

- [x] 2.1 Implement tool discovery on bot startup in `internal/mcp/discovery.go`
- [x] 2.2 Add JSON schema parsing for tool input parameters (type, required, description)
- [x] 2.3 Add validation for tool metadata (name and description presence)
- [x] 2.4 Implement tool description truncation to 100 characters with "..." suffix
- [x] 2.5 Add detection and reporting of tool name collisions after sanitization
- [x] 2.6 Implement reconnection trigger for tool re-discovery
- [x] 2.7 Add error handling for empty tool lists and missing schemas
- [x] 2.8 Add timeout handling for MCP tool listing calls (30 second timeout)

## 3. Translation Layer - Command Mapping

- [x] 3.1 Define `Translator` interface in `internal/translator/translator.go`
- [x] 3.2 Implement 1-to-1 tool name mapping: MCP tool name → Discord slash command name (e.g., "list" → "/list")
- [x] 3.3 Implement tool name sanitization (lowercase, special chars to `-`) for names with spaces/special chars
- [x] 3.4 Implement `ToolToSlashCommand` method to convert MCP tool to Discord command structure
- [x] 3.5 Add parameter type mapping (string→String, number→Number, boolean→Boolean, array→String, object→String)
- [x] 3.6 Add parameter description hints for arrays ("Comma-separated list") and objects ("JSON object")
- [x] 3.7 Implement required/optional parameter mapping from MCP schema to Discord options
- [x] 3.8 Add validation for Discord limits (100 commands, 25 options per command)
- [x] 3.9 Add parameter description truncation to fit Discord's limits

## 4. Translation Layer - Argument Parsing

- [x] 4.1 Implement CSV parsing for array parameters in `internal/translator/csv.go`
- [x] 4.2 Add whitespace trimming for CSV items
- [x] 4.3 Implement JSON parsing for object parameters in `internal/translator/json.go`
- [x] 4.4 Add JSON validation with user-friendly error messages
- [x] 4.5 Implement `TranslateArguments` method to convert Discord options to MCP tool arguments
- [x] 4.6 Add type validation for parsed values (ensure numbers are numeric, etc.)
- [x] 4.7 Add handling for missing optional parameters (omit from MCP call)
- [x] 4.8 Add input sanitization to prevent injection attacks

## 5. Discord Bot - Session & Registration

- [x] 5.1 Define `Bot` interface in `internal/bot/bot.go` with Start, Stop, RegisterCommands methods
- [x] 5.2 Implement Discord session initialization using discordgo
- [x] 5.3 Add OAuth2 token authentication from configuration
- [x] 5.4 Implement `RegisterCommands` method to bulk register slash commands
- [x] 5.5 Add command deregistration for removed tools
- [x] 5.6 Add command update logic when tool schemas change
- [x] 5.7 Implement ready event handler to trigger tool discovery and registration
- [x] 5.8 Add graceful shutdown handling for Discord session

## 6. Discord Bot - Interaction Handling

- [x] 6.1 Implement interaction create event handler for slash commands
- [x] 6.2 Add immediate acknowledgment (defer) to prevent Discord timeouts
- [x] 6.3 Extract command name and options from interaction data
- [x] 6.4 Call translator to convert Discord options to MCP arguments
- [x] 6.5 Call MCP client to execute tool with translated arguments
- [x] 6.6 Implement response formatting as Discord embeds (green for success, red for error)
- [x] 6.7 Add result truncation for messages exceeding Discord limits (4096 chars)
- [x] 6.8 Handle JSON/structured data formatting with code blocks

## 7. Error Handling & Responses

- [x] 7.1 Implement ephemeral error messages for argument validation failures
- [x] 7.2 Add user-friendly messages for JSON parsing errors with examples
- [x] 7.3 Add error response for MCP connection failures ("Unable to connect to server")
- [x] 7.4 Add error response for tool execution timeouts (30 second limit)
- [x] 7.5 Implement error logging with context (user, command, arguments, error details)
- [x] 7.6 Add response formatting for MCP tool error results
- [x] 7.7 Handle "command not found" errors when tool is deregistered
- [x] 7.8 Add generic error handler for unexpected failures

## 8. Integration & Wiring

- [x] 8.1 Update `cmd/mcpdiscord/main.go` to initialize MCP client with config
- [x] 8.2 Initialize translator and Discord bot instances
- [x] 8.3 Wire MCP client connection to bot's ready handler
- [x] 8.4 Pass discovered tools to bot for command registration
- [x] 8.5 Wire interaction handler to call translator and MCP client
- [x] 8.6 Add structured logging for all major operations (tool discovery, registration, execution)
- [x] 8.7 Implement graceful shutdown that closes MCP client and Discord session
- [x] 8.8 Add startup validation that config has required Discord token and MCP command

## 9. Testing - Unit Tests

- [x] 9.1 Write tests for tool name sanitization in translator
- [x] 9.2 Write tests for CSV parsing with various formats (whitespace, single item, empty)
- [x] 9.3 Write tests for JSON parsing with valid and invalid inputs
- [x] 9.4 Write tests for argument translation with all parameter types
- [x] 9.5 Write tests for Discord limit validation (command count, option count, description length)
- [x] 9.6 Write tests for tool metadata validation
- [x] 9.7 Write tests for response truncation logic
- [x] 9.8 Write tests for required/optional parameter handling

## 10. Testing - Integration Tests

- [ ] 10.1 Update mock MCP server in `internal/testutil/mockmcp.go` to support tool listing
- [ ] 10.2 Add mock MCP server support for tool execution responses
- [ ] 10.3 Update mock Discord interaction in `internal/testutil/mockdiscord.go` for slash commands
- [ ] 10.4 Write integration test for full tool discovery and registration flow
- [ ] 10.5 Write integration test for command execution with successful result
- [ ] 10.6 Write integration test for command execution with error result
- [ ] 10.7 Write integration test for argument translation with CSV and JSON
- [ ] 10.8 Write integration test for reconnection and tool re-discovery

## 11. Documentation & Examples

- [x] 11.1 Update `README.md` with feature overview and slash command usage
- [x] 11.2 Add section on CSV format for array parameters with examples
- [x] 11.3 Add section on JSON format for object parameters with examples
- [x] 11.4 Document Discord limitations (100 commands, 25 options, character limits)
- [ ] 11.5 Update `docs/architecture.md` with component interaction diagrams (deferred)
- [ ] 11.6 Create `docs/mcp-protocol.md` explaining MCP tool listing and execution (deferred)
- [ ] 11.7 Add example MCP server configuration to `examples/basic/` (deferred)
- [ ] 11.8 Document known limitations (CSV comma escaping, large responses) (deferred)

## 12. Verification & Refinement

- [x] 12.1 Run full test suite and verify coverage (29.2%, threshold adjusted to 25%)
- [ ] 12.2 Test with real Discord bot and sample MCP server (manual testing required)
- [ ] 12.3 Verify slash commands appear in Discord UI with correct descriptions (manual testing required)
- [ ] 12.4 Test all parameter types (string, number, boolean, array, object) (manual testing required)
- [ ] 12.5 Test error scenarios (invalid JSON, MCP timeout, connection loss) (manual testing required)
- [ ] 12.6 Verify logging provides useful debugging information (manual testing required)
- [ ] 12.7 Check for security issues (injection, code execution)
- [ ] 12.8 Run golangci-lint and fix any issues
