package mcp

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request message.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      *int64      `json:"id,omitempty"` // Pointer to support notifications (no ID)
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response message.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error object.
type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// requestIDCounter generates unique request IDs.
var requestIDCounter int64

// NewRequest creates a new JSON-RPC request with a unique ID.
func NewRequest(method string, params interface{}) *JSONRPCRequest {
	id := atomic.AddInt64(&requestIDCounter, 1)
	return &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      &id,
		Method:  method,
		Params:  params,
	}
}

// EncodeRequest marshals a JSON-RPC request to JSON bytes.
func EncodeRequest(req *JSONRPCRequest) ([]byte, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON-RPC request: %w", err)
	}
	return data, nil
}

// DecodeResponse unmarshals JSON bytes into a JSON-RPC response.
func DecodeResponse(data []byte) (*JSONRPCResponse, error) {
	var resp JSONRPCResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON-RPC response: %w", err)
	}
	return &resp, nil
}

// CheckError checks if a JSON-RPC response contains an error and returns it.
func CheckError(resp *JSONRPCResponse) error {
	if resp.Error != nil {
		return fmt.Errorf("JSON-RPC error %d: %s", resp.Error.Code, resp.Error.Message)
	}
	return nil
}

// ListToolsParams represents parameters for the tools/list method.
type ListToolsParams struct {
	// Currently no parameters required for tools/list
}

// ListToolsResult represents the result of a tools/list call.
type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// CallToolParams represents parameters for the tools/call method.
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult represents the result of a tools/call.
type CallToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}
