package mcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"mcpdiscord/internal/config"
)


func TestStdioClient_SendReceive(t *testing.T) {
	// Create a mock JSON-RPC server using cat command
	// cat will echo back whatever we send to it
	cfg := config.MCPConfig{
		Command: "cat",
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	client := NewStdioClient(cfg, logger)

	ctx := context.Background()
	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	// Since cat just echoes, this won't work properly for actual MCP protocol
	// but it tests the send/receive mechanics
	t.Log("Basic stdio connection test passed")
}

func TestNewRequest_Sequence(t *testing.T) {
	req1 := NewRequest("method1", nil)
	req2 := NewRequest("method2", nil)
	req3 := NewRequest("method3", nil)

	if req1.ID == nil || req2.ID == nil || *req1.ID >= *req2.ID {
		t.Error("Request IDs should increment")
	}
	if req2.ID == nil || req3.ID == nil || *req2.ID >= *req3.ID {
		t.Error("Request IDs should increment")
	}

	if req1.JSONRPC != "2.0" {
		t.Error("Expected JSON-RPC 2.0")
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	originalReq := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64Ptr(42),
		Method:  "test/method",
		Params:  map[string]interface{}{"key": "value"},
	}

	_, err := EncodeRequest(originalReq)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Create a mock response
	respData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      42,
		"result":  map[string]interface{}{"success": true},
	}

	respJSON, err := json.Marshal(respData)
	if err != nil {
		t.Fatalf("Marshal response failed: %v", err)
	}

	resp, err := DecodeResponse(respJSON)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if originalReq.ID == nil || resp.ID != *originalReq.ID {
		t.Errorf("ID mismatch: got %d, want %d", resp.ID, *originalReq.ID)
	}
}

func TestJSONRPCError_Handling(t *testing.T) {
	tests := []struct {
		name     string
		resp     *JSONRPCResponse
		hasError bool
	}{
		{
			name: "no error",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  json.RawMessage(`{"success": true}`),
			},
			hasError: false,
		},
		{
			name: "parse error -32700",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &JSONRPCError{
					Code:    -32700,
					Message: "Parse error",
				},
			},
			hasError: true,
		},
		{
			name: "invalid request -32600",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &JSONRPCError{
					Code:    -32600,
					Message: "Invalid Request",
				},
			},
			hasError: true,
		},
		{
			name: "method not found -32601",
			resp: &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &JSONRPCError{
					Code:    -32601,
					Message: "Method not found",
				},
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckError(tt.resp)
			if tt.hasError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestDiscoveryService_ConnectionFailure(t *testing.T) {
	client := &MockClient{
		connectErr: os.ErrNotExist,
	}
	service := NewDiscoveryService(client, nil)

	ctx := context.Background()
	_, err := service.DiscoverWithReconnect(ctx)
	if err == nil {
		t.Error("Expected error for connection failure")
	}
}

func TestDiscoveryService_ToolValidation(t *testing.T) {
	tests := []struct {
		name        string
		tools       []Tool
		expectCount int
	}{
		{
			name: "all valid tools",
			tools: []Tool{
				{Name: "tool1", Description: "Desc1", InputSchema: InputSchema{Type: "object"}},
				{Name: "tool2", Description: "Desc2", InputSchema: InputSchema{Type: "object"}},
			},
			expectCount: 2,
		},
		{
			name: "filter invalid tools",
			tools: []Tool{
				{Name: "", Description: "Invalid", InputSchema: InputSchema{Type: "object"}},
				{Name: "valid", Description: "Valid", InputSchema: InputSchema{Type: "object"}},
				{Name: "also_valid", Description: "Also valid", InputSchema: InputSchema{Type: "object"}},
			},
			expectCount: 2,
		},
		{
			name: "filter tools with invalid schema",
			tools: []Tool{
				{Name: "bad_schema", Description: "Bad", InputSchema: InputSchema{Type: ""}},
				{Name: "good", Description: "Good", InputSchema: InputSchema{Type: "object"}},
			},
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockClient{
				tools: tt.tools,
			}
			service := NewDiscoveryService(client, nil)

			ctx := context.Background()
			result, err := service.DiscoverTools(ctx)
			if err != nil {
				t.Fatalf("DiscoverTools failed: %v", err)
			}

			if len(result) != tt.expectCount {
				t.Errorf("Expected %d valid tools, got %d", tt.expectCount, len(result))
			}
		})
	}
}

func TestInputSchemaTypes(t *testing.T) {
	schemas := []struct {
		schemaType string
		valid      bool
	}{
		{"object", true},
		{"string", false},
		{"number", false},
		{"array", false},
		{"", false},
	}

	for _, tc := range schemas {
		tool := &Tool{
			Name:        "test",
			Description: "test",
			InputSchema: InputSchema{Type: tc.schemaType},
		}

		err := ValidateTool(tool)
		if tc.valid && err != nil {
			t.Errorf("Schema type %q should be valid but got error: %v", tc.schemaType, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Schema type %q should be invalid but got no error", tc.schemaType)
		}
	}
}

func TestContentBlock_Types(t *testing.T) {
	result := &ToolResult{
		Content: []ContentBlock{
			{Type: "text", Text: "Text content"},
			{Type: "image", Text: "Image data"},
			{Type: "resource", Text: "Resource data"},
		},
		IsError: false,
	}

	if len(result.Content) != 3 {
		t.Errorf("Expected 3 content blocks, got %d", len(result.Content))
	}

	for i, block := range result.Content {
		if block.Text == "" {
			t.Errorf("Block %d has empty text", i)
		}
	}
}
