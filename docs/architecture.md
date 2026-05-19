# Architecture Overview

This document describes the high-level architecture of the MCP-Discord bot, including component relationships, data flow, and key design decisions.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                       MCP-Discord Bot                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐      ┌──────────────┐      ┌──────────────┐   │
│  │   Discord   │◄────►│  Translator  │◄────►│  MCP Client  │   │
│  │   Bot       │      │              │      │              │   │
│  └─────────────┘      └──────────────┘      └──────────────┘   │
│         ▲                    ▲                       ▲          │
│         │                    │                       │          │
│         └────────────────────┴───────────────────────┘          │
│                           Config                                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
           ▲                                            ▼
           │                                            │
    ┌──────┴──────┐                            ┌───────┴────────┐
    │   Discord   │                            │   MCP Server   │
    │   Gateway   │                            │   (External)   │
    └─────────────┘                            └────────────────┘
```

## Core Components

### 1. Main Application (`cmd/mcpdiscord`)

**Responsibilities:**

- Application entry point
- Configuration loading and validation
- Component initialization and lifecycle management
- Signal handling and graceful shutdown
- Structured logging setup

**Key Features:**

- Multi-source configuration (flag > env > default)
- Environment variable interpolation
- Version information
- Signal-based shutdown (SIGINT, SIGTERM)

**Code Location:** `cmd/mcpdiscord/main.go`

---

### 2. Configuration System (`internal/config`)

**Responsibilities:**

- Parse JSON configuration files
- Interpolate environment variables
- Validate configuration
- Multi-source loading with precedence

**Components:**

- `types.go`: Configuration struct definitions
- `loader.go`: File loading and env interpolation
- `validator.go`: Validation logic with descriptive errors

**Configuration Flow:**

```yaml
1. Determine config path (flag > env > default)
2. Read JSON file
3. Interpolate ${VAR_NAME} with environment values
4. Unmarshal into typed structs
5. Validate required fields and constraints
```

**Code Location:** `internal/config/`

---

### 3. Discord Bot (`internal/bot`)

**Responsibilities:**

- Connect to Discord Gateway
- Register slash commands from MCP tools
- Handle slash command interactions
- Format and send responses

**Data Flow:**

```
Discord User Types /weather location:Seattle
         ↓
Discord sends InteractionCreate event
         ↓
Bot receives interaction
         ↓
Extracts command name and options
         ↓
Passes to Translator
```

**Key Operations:**

- Command registration (global or guild-specific)
- Interaction handling
- Response formatting (text, embeds, ephemeral)
- Error handling and user feedback

**Code Location:** `internal/bot/` (stub created, full implementation in future changes)

---

### 4. MCP Client (`internal/mcp`)

**Responsibilities:**

- Connect to MCP server via stdio/SSE/WebSocket
- Call `tools/list` to discover available tools
- Call `tools/call` to execute tools
- Handle MCP protocol messages

**Connection Types:**

- **stdio**: Subprocess communication (most common)
- **SSE**: Server-sent events over HTTP
- **WebSocket**: Bidirectional websocket connection

**Protocol Flow:**

```md
1. Launch MCP server as subprocess (if stdio)
2. Send tools/list request
3. Receive tool definitions with JSON Schema
4. When tool is invoked:
   a. Send tools/call with arguments
   b. Receive response (content blocks)
   c. Return to translator
```

**Code Location:** `internal/mcp/` (stub created, full implementation in future changes)

---

### 5. Translator (`internal/translator`)

**Responsibilities:**

- Convert MCP tool definitions → Discord command definitions
- Map JSON Schema types → Discord option types
- Convert Discord options → MCP tool arguments
- Format MCP responses → Discord messages

**Schema Translation:**

| MCP (JSON Schema) | Discord Command Option |
| ------------------- | ------------------------ |
| `string` | `STRING` |
| `number` | `NUMBER` |
| `integer` | `INTEGER` |
| `boolean` | `BOOLEAN` |
| `enum` (≤25 values) | `STRING` + Choices |
| `object` (nested) | JSON string fallback |
| `array` | JSON string fallback |

**Translation Phases:**

#### Phase 1: Tool Discovery

```json
MCP Tool Definition:
{
  "name": "get_weather",
  "description": "Get current weather",
  "inputSchema": {
    "type": "object",
    "properties": {
      "location": { "type": "string" },
      "units": { "enum": ["C", "F"] }
    }
  }
}

         ↓ Translate ↓

Discord Command:
{
  "name": "get_weather",
  "description": "Get current weather",
  "options": [
    { "name": "location", "type": STRING, "required": true },
    { "name": "units", "type": STRING, "choices": ["C", "F"] }
  ]
}
```

#### Phase 2: Execution

```json
Discord Interaction:
{ "name": "get_weather", "options": { "location": "Seattle", "units": "F" } }

         ↓ Translate ↓

MCP Tool Call:
{ "name": "get_weather", "arguments": { "location": "Seattle", "units": "F" } }

         ↓ Execute ↓

MCP Response:
{ "content": [{ "type": "text", "text": "Seattle: 72°F, Sunny" }] }

         ↓ Format ↓

Discord Response:
"Seattle: 72°F, Sunny"
```

**Code Location:** `internal/translator/` (stub created, full implementation in future changes)

---

### 6. Test Utilities (`internal/testutil`)

**Responsibilities:**

- Provide test helpers and mocks
- Simplify test setup
- Mock external dependencies

**Components:**

- `config.go`: Configuration test helpers
- `mockmcp.go`: Mock MCP server for testing
- `mockdiscord.go`: Mock Discord interactions

**Usage:**

```go
func TestMyFeature(t *testing.T) {
    cfg := testutil.TestConfig(t)
    mockMCP := testutil.NewMockMCPServer(t)
    mockMCP.AddSimpleTool("test_tool", "A test tool")
    // ... test logic
}
```

**Code Location:** `internal/testutil/`

---

## Data Flow

### Startup Flow

```md
1. Parse CLI flags (--config, --version)
2. Initialize logger
3. Load configuration
   a. Determine source (flag > env > default)
   b. Read and parse JSON
   c. Interpolate environment variables
4. Validate configuration
5. Initialize MCP client
6. Connect to MCP server
7. Discover tools (tools/list)
8. Initialize Discord bot
9. Register commands
10. Start listening for interactions
11. Wait for shutdown signal
```

### Command Execution Flow

```md
User types /weather location:Seattle in Discord
         ↓
1. Discord Gateway sends InteractionCreate
         ↓
2. Bot receives interaction
         ↓
3. Extract command name and options
         ↓
4. Translator converts to MCP format
         ↓
5. MCP Client calls tools/call
         ↓
6. MCP Server executes and returns result
         ↓
7. Translator formats response for Discord
         ↓
8. Bot sends InteractionResponse
         ↓
User sees result in Discord
```

### Shutdown Flow

```md
1. Receive SIGINT or SIGTERM
2. Cancel context
3. Stop accepting new interactions
4. Wait for in-flight requests to complete
5. Disconnect from Discord
6. Shutdown MCP client
7. Close MCP server process
8. Flush logs
9. Exit gracefully
```

---

## Package Dependencies

```file
cmd/mcpdiscord
  └─> internal/config

internal/bot
  ├─> internal/config
  ├─> internal/translator
  └─> internal/mcp

internal/translator
  └─> internal/mcp

internal/mcp
  └─> internal/config

internal/testutil
  ├─> internal/config
  ├─> internal/mcp
  └─> internal/bot
```

**Dependency Rules:**

- `internal/` packages cannot be imported by external projects
- No circular dependencies
- Each package has a single, well-defined responsibility
- Test utilities can import from any package

---

## Concurrency Model

### Goroutines

- **Main goroutine**: Configuration, initialization, shutdown coordination
- **Discord event handler**: Processes incoming Discord interactions
- **MCP client**: Manages subprocess stdio streams
- **Signal handler**: Listens for OS signals

### Synchronization

- **Context**: Used for cancellation and shutdown coordination
- **Channels**: Signal handling, shutdown notifications
- **Mutexes**: Protect shared state (if any)

### Error Handling

- Errors propagate up via return values
- Structured logging for operational errors
- Graceful degradation where possible
- Clear error messages to users

---

## Security Considerations

### Configuration

- Environment variables for secrets
- Never log sensitive values
- Validate all inputs

### Discord

- Verify interaction signatures (handled by discordgo)
- Rate limit handling
- Permission checks

### MCP

- Subprocess isolation
- Environment variable injection security
- Input sanitization before calling tools

---

## Performance Characteristics

### Bottlenecks

1. **Discord API rate limits**: ~50 command registrations per 10 seconds
2. **MCP tool execution**: Depends on server implementation
3. **Network latency**: Discord Gateway + MCP server

### Optimization Strategies

- Guild-specific commands for development (instant vs 1 hour)
- Connection pooling for SSE/WebSocket transports
- Caching of tool definitions
- Parallel tool discovery

---

## Future Architecture Extensions

### Multi-Server Support

```file
One bot → Multiple MCP servers
├─ Server A: Weather tools (prefix: weather_)
├─ Server B: Search tools (prefix: search_)
└─ Server C: Calc tools (prefix: calc_)
```

### Command Caching

- Cache translated commands
- Update only when MCP server restarts
- Reduce registration overhead

### Advanced Translation

- Handle complex nested objects
- Support for attachments/images
- Modal forms for multi-step workflows

### Monitoring

- Prometheus metrics
- Distributed tracing
- Health check endpoints
