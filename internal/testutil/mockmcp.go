package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

// MockMCPServer is a mock MCP server for testing.
// It simulates the MCP protocol over stdio.
type MockMCPServer struct {
	t       *testing.T
	Tools   []MCPTool
	Calls   []MCPToolCall
	stdin   io.Reader
	stdout  io.Writer
}

// MCPTool represents an MCP tool definition.
type MCPTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// MCPToolCall represents a call to an MCP tool.
type MCPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPToolsListResponse represents the response to tools/list.
type MCPToolsListResponse struct {
	Tools []MCPTool `json:"tools"`
}

// MCPToolCallRequest represents a request to call a tool.
type MCPToolCallRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPToolCallResponse represents the response from a tool call.
type MCPToolCallResponse struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent represents content in an MCP response.
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// NewMockMCPServer creates a new mock MCP server.
func NewMockMCPServer(t *testing.T) *MockMCPServer {
	t.Helper()
	return &MockMCPServer{
		t:     t,
		Tools: []MCPTool{},
		Calls: []MCPToolCall{},
	}
}

// AddTool adds a tool to the mock server.
func (m *MockMCPServer) AddTool(name, description string, inputSchema json.RawMessage) {
	m.Tools = append(m.Tools, MCPTool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
	})
}

// AddSimpleTool adds a simple tool with a basic string parameter schema.
func (m *MockMCPServer) AddSimpleTool(name, description string) {
	schema := json.RawMessage(`{
		"type": "object",
		"properties": {
			"input": {
				"type": "string",
				"description": "Input parameter"
			}
		},
		"required": ["input"]
	}`)
	m.AddTool(name, description, schema)
}

// GetToolsList returns the tools/list response.
func (m *MockMCPServer) GetToolsList() MCPToolsListResponse {
	return MCPToolsListResponse{
		Tools: m.Tools,
	}
}

// HandleToolCall simulates handling a tool call and returns a mock response.
func (m *MockMCPServer) HandleToolCall(req MCPToolCallRequest) MCPToolCallResponse {
	// Record the call
	m.Calls = append(m.Calls, MCPToolCall{
		Name:      req.Name,
		Arguments: req.Arguments,
	})

	// Return a mock success response
	return MCPToolCallResponse{
		Content: []MCPContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Mock response for tool %s", req.Name),
			},
		},
		IsError: false,
	}
}

// GetCalls returns all recorded tool calls.
func (m *MockMCPServer) GetCalls() []MCPToolCall {
	return m.Calls
}

// Reset clears all recorded calls.
func (m *MockMCPServer) Reset() {
	m.Calls = []MCPToolCall{}
}
