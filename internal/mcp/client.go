// Package mcp implements the Model Context Protocol client for communicating
// with MCP servers.
//
// This package provides functionality to connect to MCP servers via various
// transports (stdio, SSE, WebSocket), discover available tools, and execute
// tool calls. It handles the low-level protocol details and message serialization.
//
// The MCP protocol allows servers to expose tools (functions) that can be
// discovered and invoked by clients. Each tool has a name, description, and
// JSON Schema defining its input parameters. Tool execution returns structured
// content blocks.
//
// Key responsibilities:
//   - Establish and maintain connection to MCP server
//   - Implement MCP protocol message handling
//   - Tool discovery (tools/list request)
//   - Tool execution (tools/call request)
//   - Transport abstraction (stdio, SSE, WebSocket)
//   - Error handling and retries
//
// Example usage:
//
//	cfg := config.MCPConfig{
//	    Command: "npx",
//	    Args: []string{"-y", "@modelcontextprotocol/server-weather"},
//	    Transport: "stdio",
//	}
//
//	client := mcp.NewClient(cfg)
//	if err := client.Connect(); err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	tools, err := client.ListTools()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	result, err := client.CallTool("get_weather", map[string]interface{}{
//	    "location": "Seattle",
//	})
package mcp

import (
	"context"
	"fmt"

	"mcpdiscord/internal/config"
)

// Client represents an MCP protocol client that communicates with an MCP server.
type Client interface {
	// Connect establishes a connection to the MCP server using the configured
	// transport. Returns an error if connection fails.
	Connect(ctx context.Context) error

	// Close gracefully shuts down the connection to the MCP server and
	// cleans up resources.
	Close() error

	// ListTools retrieves the list of available tools from the MCP server.
	// This calls the tools/list MCP method.
	ListTools(ctx context.Context) ([]Tool, error)

	// CallTool invokes a specific tool on the MCP server with the given arguments.
	// This calls the tools/call MCP method.
	CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*ToolResult, error)

	// IsConnected returns whether the client is currently connected to the server.
	IsConnected() bool
}

// Tool represents an MCP tool definition returned from tools/list.
type Tool struct {
	// Name is the unique identifier for the tool
	Name string `json:"name"`

	// Description explains what the tool does
	Description string `json:"description"`

	// InputSchema is the JSON Schema defining the tool's input parameters
	InputSchema Schema `json:"inputSchema"`
}

// Schema represents a JSON Schema definition for tool inputs.
type Schema struct {
	// Type is the JSON type (object, string, number, etc.)
	Type string `json:"type"`

	// Properties defines the schema for object properties (if Type is "object")
	Properties map[string]Property `json:"properties,omitempty"`

	// Required lists the names of required properties
	Required []string `json:"required,omitempty"`

	// Items defines the schema for array items (if Type is "array")
	Items *Schema `json:"items,omitempty"`

	// Enum lists allowed values (for enum types)
	Enum []interface{} `json:"enum,omitempty"`
}

// Property represents a JSON Schema property definition.
type Property struct {
	// Type is the JSON type of this property
	Type string `json:"type"`

	// Description explains what this property is for
	Description string `json:"description,omitempty"`

	// Enum lists allowed values (for enum types)
	Enum []interface{} `json:"enum,omitempty"`

	// Default is the default value if not provided
	Default interface{} `json:"default,omitempty"`

	// Properties for nested objects
	Properties map[string]Property `json:"properties,omitempty"`

	// Items for arrays
	Items *Schema `json:"items,omitempty"`
}

// ToolResult represents the result of a tool execution.
type ToolResult struct {
	// Content contains the response content blocks from the tool
	Content []ContentBlock `json:"content"`

	// IsError indicates whether the tool execution resulted in an error
	IsError bool `json:"isError,omitempty"`
}

// ContentBlock represents a content block in an MCP response.
type ContentBlock struct {
	// Type is the content type (text, image, resource, etc.)
	Type string `json:"type"`

	// Text is the text content (for Type "text")
	Text string `json:"text,omitempty"`

	// Data is additional data for other content types
	Data map[string]interface{} `json:"data,omitempty"`
}

// Transport represents the communication protocol used to connect to the MCP server.
type Transport string

const (
	// TransportStdio uses standard input/output for communication (subprocess)
	TransportStdio Transport = "stdio"

	// TransportSSE uses Server-Sent Events over HTTP
	TransportSSE Transport = "sse"

	// TransportWebSocket uses WebSocket for bidirectional communication
	TransportWebSocket Transport = "websocket"
)

// client is the concrete implementation of the Client interface.
// TODO: Implement in future change
type client struct {
	config    config.MCPConfig
	transport Transport
	connected bool
}

// NewClient creates a new MCP client with the given configuration.
// TODO: Implement in future change
func NewClient(cfg config.MCPConfig) (Client, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// Request represents an MCP protocol request message.
type Request struct {
	// JSONRPC is the JSON-RPC version (always "2.0")
	JSONRPC string `json:"jsonrpc"`

	// ID is the request ID for matching responses
	ID interface{} `json:"id"`

	// Method is the MCP method name (e.g., "tools/list", "tools/call")
	Method string `json:"method"`

	// Params are the method parameters
	Params interface{} `json:"params,omitempty"`
}

// Response represents an MCP protocol response message.
type Response struct {
	// JSONRPC is the JSON-RPC version (always "2.0")
	JSONRPC string `json:"jsonrpc"`

	// ID is the request ID this response is for
	ID interface{} `json:"id"`

	// Result contains the successful result (if no error)
	Result interface{} `json:"result,omitempty"`

	// Error contains error information (if request failed)
	Error *ResponseError `json:"error,omitempty"`
}

// ResponseError represents an error in an MCP response.
type ResponseError struct {
	// Code is the error code
	Code int `json:"code"`

	// Message is the error message
	Message string `json:"message"`

	// Data contains additional error data
	Data interface{} `json:"data,omitempty"`
}
