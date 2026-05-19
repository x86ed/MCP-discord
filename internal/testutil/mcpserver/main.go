package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSON-RPC 2.0 structures
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP Protocol structures
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

type ToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse request: %v\n", err)
			continue
		}

		resp := handleRequest(req)
		
		data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal response: %v\n", err)
			continue
		}

		fmt.Println(string(data))
	}
}

func handleRequest(req Request) Response {
	switch req.Method {
	case "initialize":
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "0.1.0",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "test-mcp-server",
					"version": "1.0.0",
				},
			},
		}

	case "tools/list":
		tools := []Tool{
			{
				Name:        "echo",
				Description: "Echoes back the input message",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"message": map[string]interface{}{
							"type":        "string",
							"description": "Message to echo",
						},
					},
					Required: []string{"message"},
				},
			},
			{
				Name:        "add",
				Description: "Adds two numbers together",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number",
						},
					},
					Required: []string{"a", "b"},
				},
			},
			{
				Name:        "greet",
				Description: "Greets a person by name",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Name to greet",
						},
					},
					Required: []string{"name"},
				},
			},
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": tools,
			},
		}

	case "tools/call":
		return handleToolCall(req)

	default:
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func handleToolCall(req Request) Response {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	var result ToolResult

	switch params.Name {
	case "echo":
		message, ok := params.Arguments["message"].(string)
		if !ok {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: "Error: message must be a string"},
				},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: message},
				},
				IsError: false,
			}
		}

	case "add":
		a, aOk := params.Arguments["a"].(float64)
		b, bOk := params.Arguments["b"].(float64)
		if !aOk || !bOk {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: "Error: a and b must be numbers"},
				},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: fmt.Sprintf("Result: %.2f", a+b)},
				},
				IsError: false,
			}
		}

	case "greet":
		name, ok := params.Arguments["name"].(string)
		if !ok {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: "Error: name must be a string"},
				},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ContentBlock{
					{Type: "text", Text: fmt.Sprintf("Hello, %s!", strings.TrimSpace(name))},
				},
				IsError: false,
			}
		}

	default:
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: fmt.Sprintf("Unknown tool: %s", params.Name),
			},
		}
	}

	return Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}
