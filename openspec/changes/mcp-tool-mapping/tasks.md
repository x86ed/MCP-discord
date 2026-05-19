## 1. MCP Client Foundation

- [ ] 1.1 Define `Client` interface in `internal/mcp/client.go` with Connect, ListTools, CallTool, Close methods
- [ ] 1.2 Define `Tool` struct to represent MCP tool metadata (name, description, input schema)
- [ ] 1.3 Define `ToolResult` struct to represent MCP tool execution results
- [ ] 1.4 Implement stdio transport client in `internal/mcp/stdio.go` with subprocess management
- [ ] 1.5 Add JSON-RPC message encoding/decoding utilities in `internal/mcp/jsonrpc.go`
- [ ] 1.6 Implement `ListTools` method to call MCP "tools/list" endpoint
- [ ] 1.7 Implement `CallTool` method to call MCP "tools/call" endpoint with arguments
- [ ] 1.8 Add connection lifecycle management (Connect, Close, reconnection logic)

## 2. Tool Discovery & Schema Parsing

- [ ] 2.1 Implement tool discovery on bot startup in `internal/mcp/discovery.go`
- [ ] 2.2 Add JSON schema parsing for tool input parameters (type, required, description)
- [ ] 2.3 Add validation for tool metadata (name and description presence)
- [ ] 2.4 Implement tool description truncation to 100 characters with "..." suffix
- [ ] 2.5 Add detection and reporting of tool name collisions after sanitization
- [ ] 2.6 Implement reconnection trigger for tool re-discovery
- [ ] 2.7 Add error handling for empty tool lists and missing schemas
- [ ] 2.8 Add timeout handling for MCP tool listing calls (30 second timeout)

## 3. Translation Layer - Command Mapping

- [ ] 3.1 Define `Translator` interface in `internal/translator/translator.go`
- [ ] 3.2 Implement 1-to-1 tool name mapping: MCP tool name → Discord slash command name (e.g., "list" → "/list")
- [ ] 3.3 Implement tool name sanitization (lowercase, special chars to `-`) for names with spaces/special chars
- [ ] 3.4 Implement `ToolToSlashCommand` method to convert MCP tool to Discord command structure
- [ ] 3.5 Add parameter type mapping (string→String, number→Number, boolean→Boolean, array→String, object→String)
- [ ] 3.6 Add parameter description hints for arrays ("Comma-separated list") and objects ("JSON object")
- [ ] 3.7 Implement required/optional parameter mapping from MCP schema to Discord options
- [ ] 3.8 Add validation for Discord limits (100 commands, 25 options per command)
- [ ] 3.9 Add parameter description truncation to fit Discord's limits

## 4. Translation Layer - Argument Parsing

- [ ] 4.1 Implement CSV parsing for array parameters in `internal/translator/csv.go`
- [ ] 4.2 Add whitespace trimming for CSV items
- [ ] 4.3 Implement JSON parsing for object parameters in `internal/translator/json.go`
- [ ] 4.4 Add JSON validation with user-friendly error messages
- [ ] 4.5 Implement `TranslateArguments` method to convert Discord options to MCP tool arguments
- [ ] 4.6 Add type validation for parsed values (ensure numbers are numeric, etc.)
- [ ] 4.7 Add handling for missing optional parameters (omit from MCP call)
- [ ] 4.8 Add input sanitization to prevent injection attacks

## 5. Discord Bot - Session & Registration

- [ ] 5.1 Define `Bot` interface in `internal/bot/bot.go` with Start, Stop, RegisterCommands methods
- [ ] 5.2 Implement Discord session initialization using discordgo
- [ ] 5.3 Add OAuth2 token authentication from configuration
- [ ] 5.4 Implement `RegisterCommands` method to bulk register slash commands
- [ ] 5.5 Add command deregistration for removed tools
- [ ] 5.6 Add command update logic when tool schemas change
- [ ] 5.7 Implement ready event handler to trigger tool discovery and registration
- [ ] 5.8 Add graceful shutdown handling for Discord session

## 6. Discord Bot - Interaction Handling

- [ ] 6.1 Implement interaction create event handler for slash commands
- [ ] 6.2 Add immediate acknowledgment (defer) to prevent Discord timeouts
- [ ] 6.3 Extract command name and options from interaction data
- [ ] 6.4 Call translator to convert Discord options to MCP arguments
- [ ] 6.5 Call MCP client to execute tool with translated arguments
- [ ] 6.6 Implement response formatting as Discord embeds (green for success, red for error)
- [ ] 6.7 Add result truncation for messages exceeding Discord limits (4096 chars)
- [ ] 6.8 Handle JSON/structured data formatting with code blocks

## 7. Error Handling & Responses

- [ ] 7.1 Implement ephemeral error messages for argument validation failures
- [ ] 7.2 Add user-friendly messages for JSON parsing errors with examples
- [ ] 7.3 Add error response for MCP connection failures ("Unable to connect to server")
- [ ] 7.4 Add error response for tool execution timeouts (30 second limit)
- [ ] 7.5 Implement error logging with context (user, command, arguments, error details)
- [ ] 7.6 Add response formatting for MCP tool error results
- [ ] 7.7 Handle "command not found" errors when tool is deregistered
- [ ] 7.8 Add generic error handler for unexpected failures

## 8. Integration & Wiring

- [ ] 8.1 Update `cmd/mcpdiscord/main.go` to initialize MCP client with config
- [ ] 8.2 Initialize translator and Discord bot instances
- [ ] 8.3 Wire MCP client connection to bot's ready handler
- [ ] 8.4 Pass discovered tools to bot for command registration
- [ ] 8.5 Wire interaction handler to call translator and MCP client
- [ ] 8.6 Add structured logging for all major operations (tool discovery, registration, execution)
- [ ] 8.7 Implement graceful shutdown that closes MCP client and Discord session
- [ ] 8.8 Add startup validation that config has required Discord token and MCP command

## 9. Testing - Unit Tests

- [ ] 9.1 Write tests for tool name sanitization in translator
- [ ] 9.2 Write tests for CSV parsing with various formats (whitespace, single item, empty)
- [ ] 9.3 Write tests for JSON parsing with valid and invalid inputs
- [ ] 9.4 Write tests for argument translation with all parameter types
- [ ] 9.5 Write tests for Discord limit validation (command count, option count, description length)
- [ ] 9.6 Write tests for tool metadata validation
- [ ] 9.7 Write tests for response truncation logic
- [ ] 9.8 Write tests for required/optional parameter handling

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

- [ ] 11.1 Update `README.md` with feature overview and slash command usage
- [ ] 11.2 Add section on CSV format for array parameters with examples
- [ ] 11.3 Add section on JSON format for object parameters with examples
- [ ] 11.4 Document Discord limitations (100 commands, 25 options, character limits)
- [ ] 11.5 Update `docs/architecture.md` with component interaction diagrams
- [ ] 11.6 Create `docs/mcp-protocol.md` explaining MCP tool listing and execution
- [ ] 11.7 Add example MCP server configuration to `examples/basic/`
- [ ] 11.8 Document known limitations (CSV comma escaping, large responses)

## 12. Verification & Refinement

- [ ] 12.1 Run full test suite and verify 95%+ coverage
- [ ] 12.2 Test with real Discord bot and sample MCP server
- [ ] 12.3 Verify slash commands appear in Discord UI with correct descriptions
- [ ] 12.4 Test all parameter types (string, number, boolean, array, object)
- [ ] 12.5 Test error scenarios (invalid JSON, MCP timeout, connection loss)
- [ ] 12.6 Verify logging provides useful debugging information
- [ ] 12.7 Check for security issues (injection, code execution)
- [ ] 12.8 Run golangci-lint and fix any issues
