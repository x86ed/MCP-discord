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
)

// Client defines the interface for interacting with an MCP server.
// Implementations handle transport-specific details (stdio, SSE, WebSocket).
type Client interface {
	// Connect establishes a connection to the MCP server.
	// Returns an error if the connection fails.
	Connect(ctx context.Context) error

	// ListTools retrieves all available tools from the MCP server.
	// Returns an error if the request fails or times out.
	ListTools(ctx context.Context) ([]Tool, error)

	// CallTool executes a tool with the given name and arguments.
	// Returns the tool result or an error if execution fails.
	CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*ToolResult, error)

	// Close terminates the connection to the MCP server.
	Close() error
}

// Tool represents metadata for an MCP tool.
type Tool struct {
	// Name is the unique identifier for the tool.
	Name string `json:"name"`

	// Description provides information about what the tool does.
	Description string `json:"description,omitempty"`

	// InputSchema defines the expected parameters as a JSON Schema.
	// The schema includes type information, required fields, and descriptions.
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the JSON Schema for tool parameters.
type InputSchema struct {
	// Type is typically "object" for tool parameters.
	Type string `json:"type"`

	// Properties maps parameter names to their schema definitions.
	Properties map[string]PropertySchema `json:"properties,omitempty"`

	// Required lists the names of required parameters.
	Required []string `json:"required,omitempty"`
}

// PropertySchema defines the schema for a single parameter.
type PropertySchema struct {
	// Type defines the parameter type (string, number, boolean, array, object).
	Type string `json:"type"`

	// Description explains the parameter's purpose.
	Description string `json:"description,omitempty"`

	// Items defines the schema for array elements (only for type="array").
	Items *PropertySchema `json:"items,omitempty"`

	// Properties defines nested object structure (only for type="object").
	Properties map[string]PropertySchema `json:"properties,omitempty"`
}

// ToolResult represents the result of a tool execution.
type ToolResult struct {
	// Content contains the structured result data from the tool.
	// MCP tools return an array of content blocks.
	Content []ContentBlock `json:"content"`

	// IsError indicates if the result represents an error.
	IsError bool `json:"isError,omitempty"`
}

// ContentBlock represents a single content item in a tool result.
type ContentBlock struct {
	// Type is the content type (typically "text").
	Type string `json:"type"`

	// Text contains the textual content.
	Text string `json:"text,omitempty"`
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
