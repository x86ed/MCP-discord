## Context

The project has a scaffold in place with configuration loading, but no runtime functionality. We need to implement the core feature: discovering MCP tools and exposing them as Discord slash commands. The MCP specification defines a JSON-RPC protocol for tool discovery and execution. Discord's slash command API requires pre-registration with defined parameters. We need to bridge these two systems dynamically.

Current state:
- Configuration system loads MCP server connection details (command, args, transport)
- Package stubs exist for `internal/bot`, `internal/mcp`, and `internal/translator`
- No MCP client or Discord bot implementation yet

## Goals / Non-Goals

**Goals:**
- Dynamically discover all tools from the configured MCP server
- Register each tool as a Discord slash command with appropriate parameter types
- Translate Discord interaction arguments to MCP tool input format
- Execute MCP tools and return results to Discord users
- Support all MCP argument types (strings, numbers, booleans, arrays, objects)
- Provide clear error messages when tools fail or arguments are invalid

**Non-Goals:**
- Custom command filtering or whitelisting (can be added later if needed)
- Command permission management beyond Discord's default guild-level permissions
- Multi-server MCP aggregation (single server only as per scaffold spec)
- Persistent command history or audit logging
- Rate limiting beyond what Discord provides

## Decisions

### 1. Tool Discovery Timing: Startup + Reconnect

**Decision**: Discover tools at bot startup and re-discover on MCP connection loss/reconnect.

**Rationale**: 
- MCP servers are typically long-lived processes with stable tool sets
- Startup discovery keeps runtime simple and predictable
- Re-discovery on reconnect handles server restarts gracefully
- Alternative (polling) adds complexity and unnecessary API calls

**Trade-off**: Tool changes require bot restart unless connection is lost. This is acceptable for v1.

### 2. Slash Command Mapping Strategy

**Decision**: Map MCP tool names directly to slash command names, with schema properties to command options.

**Mapping rules**:
- Tool name → command name (sanitized: lowercase, replace spaces/special chars with `-`)
- Tool description → command description (truncated to Discord's 100 char limit)
- Each input schema property → slash command option
- Property types map as: `string`→String, `number`→Number, `boolean`→Boolean, `array`→String (CSV), `object`→String (JSON)

**Rationale**:
- 1:1 mapping is simple and predictable for users
- Discord's 100-command-per-guild limit is generous for typical MCP servers
- Type mapping balances Discord's constraints with MCP's flexibility

**Alternative considered**: Command prefixing (e.g., `/mcp-*`). Rejected because it adds noise without clear benefit.

### 3. Argument Translation: CSV for Arrays, JSON for Objects

**Decision**: 
- Array parameters accept comma-separated values (CSV): `item1,item2,item3`
- Object parameters accept JSON strings: `{"key": "value", "nested": {"x": 1}}`
- Whitespace trimming for CSV items

**Rationale**:
- Discord slash command options are single strings - we need a serialization format
- CSV is intuitive for simple lists (most common case)
- JSON is standard for structured data, users familiar with it
- Alternative (query-string format `key=val&key2=val2`) less familiar for objects

**Trade-off**: CSV can't represent arrays of objects. In that case, users must use JSON array notation.

### 4. MCP Client Architecture: Interface + stdio/sse/websocket Implementations

**Decision**: Define `Client` interface in `internal/mcp`, with concrete implementations for each transport type.

```go
type Client interface {
    Connect() error
    ListTools() ([]Tool, error)
    CallTool(name string, arguments map[string]interface{}) (*ToolResult, error)
    Close() error
}
```

**Rationale**:
- Interface allows testing with mocks
- Transport abstraction matches MCP spec's flexibility
- stdio implementation is simplest and most common (priority)
- sse and websocket can be added incrementally

**Implementation priority**: stdio first, others in future changes.

### 5. Error Handling Strategy

**Decision**: Surface MCP errors to Discord users with context, log detailed errors server-side.

- MCP connection errors → Discord message: "Unable to connect to MCP server. Try again later."
- Tool execution errors → Discord message: Tool error message from MCP + hint
- Argument parsing errors → Discord ephemeral message: "Invalid format for X. Expected: Y"

**Rationale**:
- Users need actionable feedback
- Security: Don't leak internal paths or stack traces
- Ephemeral messages for validation errors reduce channel noise

### 6. Response Formatting

**Decision**: Format tool results as Discord embeds with structured fields.

- Success: Green embed with tool name as title, result as description (truncated to 4096 chars)
- Error: Red embed with error message
- Large responses: Truncate with "...see logs for full output" message

**Rationale**:
- Embeds are more readable than plain text
- Color coding provides instant feedback
- Discord message limits (2000 chars content, 4096 in embed description) require truncation

## Risks / Trade-offs

### Risk: Command Name Collisions
**Scenario**: Multiple tools with names that sanitize to the same slash command name (e.g., "Get Data" and "get_data" both become "get-data").

**Mitigation**: 
- Detect collisions during registration
- Fail fast with clear error message listing conflicting tools
- Log warning with suggestion to rename tools in MCP server
- Future: Add config option for command prefix/suffix

### Risk: Discord Slash Command Limit (100 per guild)
**Scenario**: MCP server exposes >100 tools.

**Mitigation**:
- Document the limit in README and error messages
- Future: Add tool filtering config option
- For v1: Fail registration with clear message listing tool count

### Risk: Large Tool Responses Exceed Discord Limits
**Scenario**: Tool returns JSON or text >4096 characters.

**Mitigation**:
- Truncate with "..." suffix
- Log full response server-side
- Future: Pagination or file upload for large responses

### Risk: Invalid JSON in Object Parameters
**Scenario**: User provides malformed JSON for object parameter.

**Mitigation**:
- Parse JSON during argument translation
- Return ephemeral error message with example
- Consider: Add JSON validation hints in command option descriptions

### Risk: MCP Server Crashes or Hangs
**Scenario**: stdio process dies or stops responding.

**Mitigation**:
- Set timeouts on MCP calls (e.g., 30 seconds)
- Detect broken pipe / EOF and attempt reconnect
- Gracefully handle and report to user: "Server connection lost"
- Future: Health check monitoring

## Migration Plan

Not applicable - this is the initial implementation of core functionality. No existing state to migrate.

## Open Questions

1. **Command description length**: If MCP tool description >100 chars, truncate or abbreviate intelligently?
   - **Lean**: Truncate at 97 chars + "..."
   
2. **Required vs optional parameters**: Should all slash command options be required, or respect MCP schema's required/optional?
   - **Lean**: Respect MCP schema - mark Discord options as required only if MCP requires them
   
3. **Array element types**: How to communicate expected element types in CSV? (e.g., `"1,2,3"` vs `"a,b,c"`)
   - **Lean**: Include hint in option description: "Comma-separated list of X"
   - Let MCP server validate and return errors for type mismatches
