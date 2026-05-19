package mcp

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest("test/method", map[string]interface{}{"key": "value"})

	if req.JSONRPC != "2.0" {
		t.Errorf("JSONRPC version = %q, want %q", req.JSONRPC, "2.0")
	}

	if req.Method != "test/method" {
		t.Errorf("Method = %q, want %q", req.Method, "test/method")
	}

	if req.ID == 0 {
		t.Error("Expected non-zero ID")
	}

	// Test that IDs increment
	req2 := NewRequest("another/method", nil)
	if req2.ID <= req.ID {
		t.Errorf("Expected ID to increment: %d <= %d", req2.ID, req.ID)
	}
}

func TestEncodeRequest(t *testing.T) {
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      123,
		Method:  "test/method",
		Params:  map[string]interface{}{"key": "value"},
	}

	data, err := EncodeRequest(req)
	if err != nil {
		t.Fatalf("EncodeRequest failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty encoded data")
	}

	// Should be valid JSON
	if data[0] != '{' {
		t.Error("Expected JSON object to start with '{'")
	}
}

func TestDecodeResponse(t *testing.T) {
	data := []byte(`{"jsonrpc":"2.0","id":123,"result":{"success":true}}`)

	resp, err := DecodeResponse(data)
	if err != nil {
		t.Fatalf("DecodeResponse failed: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}

	if resp.ID != 123 {
		t.Errorf("ID = %d, want %d", resp.ID, 123)
	}

	if len(resp.Result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestCheckError(t *testing.T) {
	tests := []struct {
		name        string
		resp        *JSONRPCResponse
		shouldError bool
	}{
		{
			name: "no error",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  []byte(`{}`),
			},
			shouldError: false,
		},
		{
			name: "with error",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &JSONRPCError{
					Code:    -32600,
					Message: "Invalid request",
				},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckError(tt.resp)
			if tt.shouldError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateTool(t *testing.T) {
	tests := []struct {
		name        string
		tool        *Tool
		shouldError bool
	}{
		{
			name: "valid tool",
			tool: &Tool{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: InputSchema{
					Type:       "object",
					Properties: map[string]PropertySchema{},
				},
			},
			shouldError: false,
		},
		{
			name: "missing name",
			tool: &Tool{
				Description: "A test tool",
				InputSchema: InputSchema{
					Type:       "object",
					Properties: map[string]PropertySchema{},
				},
			},
			shouldError: true,
		},
		{
			name: "missing schema type",
			tool: &Tool{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: InputSchema{},
			},
			shouldError: true,
		},
		{
			name: "invalid schema type",
			tool: &Tool{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: InputSchema{
					Type: "array",
				},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTool(tt.tool)
			if tt.shouldError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestTruncateDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "short description",
			input:    "Short",
			expected: "Short",
		},
		{
			name:     "exactly 100 chars",
			input:    "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			expected: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		},
		{
			name:     "over 100 chars",
			input:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			expected: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateDescription(tt.input)
			if result != tt.expected {
				t.Errorf("TruncateDescription() length = %d, want %d", len(result), len(tt.expected))
			}
			if len(result) > 100 {
				t.Errorf("Result exceeds 100 characters: %d", len(result))
			}
		})
	}
}

func TestCheckNameCollisions(t *testing.T) {
	tests := []struct {
		name        string
		tools       []Tool
		shouldError bool
	}{
		{
			name: "no collisions",
			tools: []Tool{
				{Name: "tool1"},
				{Name: "tool2"},
				{Name: "tool3"},
			},
			shouldError: false,
		},
		{
			name: "collision with space and underscore",
			tools: []Tool{
				{Name: "get data"},
				{Name: "get_data"},
			},
			shouldError: true,
		},
		{
			name: "collision with case difference",
			tools: []Tool{
				{Name: "GetData"},
				{Name: "getdata"},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckNameCollisions(tt.tools)
			if tt.shouldError && err == nil {
				t.Error("Expected collision error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
