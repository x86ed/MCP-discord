package mcp

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEncodeRequest_Error(t *testing.T) {
	// Create a request with a channel (unmarshalable type)
	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      int64Ptr(1),
		Method:  "test",
		Params:  make(chan int), // Channels can't be marshaled to JSON
	}

	_, err := EncodeRequest(req)
	if err == nil {
		t.Error("Expected error encoding request with channel")
	}
}

func TestDecodeResponse_Error(t *testing.T) {
	// Invalid JSON
	invalidJSON := []byte("{invalid json")
	_, err := DecodeResponse(invalidJSON)
	if err == nil {
		t.Error("Expected error decoding invalid JSON")
	}
}

func TestDiscoveryService_ListToolsError(t *testing.T) {
	client := &MockClient{
		err: context.DeadlineExceeded,
	}
	service := NewDiscoveryService(client, nil)

	ctx := context.Background()
	_, err := service.DiscoverTools(ctx)
	if err == nil {
		t.Error("Expected error when ListTools fails")
	}
}

func TestDiscoveryService_ReconnectMaxAttempts(t *testing.T) {
	client := &MockClient{
		tools: []Tool{}, // Empty list will trigger reconnect
	}
	service := NewDiscoveryService(client, nil)

	ctx := context.Background()
	_, err := service.DiscoverWithReconnect(ctx)
	// Should fail after max reconnect attempts
	if err == nil {
		t.Error("Expected error after max reconnect attempts")
	}
}

func TestJSONRPCResponse_NilError(t *testing.T) {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"success": true}`),
		Error:   nil,
	}

	err := CheckError(resp)
	if err != nil {
		t.Errorf("CheckError should return nil for response without error, got: %v", err)
	}
}

func TestToolValidation_MissingName(t *testing.T) {
	tool := &Tool{
		Name:        "",
		Description: "Test tool",
		InputSchema: InputSchema{Type: "object"},
	}

	err := ValidateTool(tool)
	if err == nil {
		t.Error("Expected error for tool with empty name")
	}
}

func TestToolValidation_MissingDescription(t *testing.T) {
	tool := &Tool{
		Name:        "test",
		Description: "",
		InputSchema: InputSchema{Type: "object"},
	}

	// Empty description should be valid (will get default)
	err := ValidateTool(tool)
	if err != nil {
		t.Errorf("ValidateTool should allow empty description: %v", err)
	}
}

func TestToolValidation_InvalidSchemaType(t *testing.T) {
	tool := &Tool{
		Name:        "test",
		Description: "Test",
		InputSchema: InputSchema{Type: "invalid_type"},
	}

	err := ValidateTool(tool)
	if err == nil {
		t.Error("Expected error for invalid schema type")
	}
}

func TestSanitizeToolName_RemoveInvalidChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"valid-tool", "valid-tool"},
		{"tool_name", "tool-name"},
		{"Tool Name", "tool-name"},
		{"tool@#$%name", "toolname"},
		{"123tool", "123tool"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		result := sanitizeToolName(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeToolName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestTruncateDescription_Short(t *testing.T) {
	desc := "Short description"
	result := TruncateDescription(desc)
	if result != desc {
		t.Errorf("Short description should not be truncated")
	}
}

func TestTruncateDescription_Long(t *testing.T) {
	desc := "This is a very long description that exceeds the maximum length allowed and should be truncated with ellipsis"
	result := TruncateDescription(desc)
	if len(result) > 100 {
		t.Errorf("Description should be truncated to 100 chars, got %d", len(result))
	}
}

func TestCheckNameCollisions_WithCollision(t *testing.T) {
	tools := []Tool{
		{Name: "test_tool"},
		{Name: "test-tool"}, // Will sanitize to same name
	}

	err := CheckNameCollisions(tools)
	if err == nil {
		t.Error("Expected error for name collision")
	}
}

func TestCheckNameCollisions_NoCollision(t *testing.T) {
	tools := []Tool{
		{Name: "tool1"},
		{Name: "tool2"},
		{Name: "tool3"},
	}

	err := CheckNameCollisions(tools)
	if err != nil {
		t.Errorf("Should not error when no collisions: %v", err)
	}
}

func TestParseInputSchema_EmptyProperties(t *testing.T) {
	schema := InputSchema{
		Type:       "object",
		Properties: nil,
		Required:   []string{},
	}

	result, err := ParseInputSchema(schema)
	if err != nil {
		t.Fatalf("ParseInputSchema failed: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestParseInputSchema_WithRequired(t *testing.T) {
	schema := InputSchema{
		Type: "object",
		Properties: map[string]PropertySchema{
			"name": {Type: "string", Description: "Name field"},
			"age":  {Type: "integer", Description: "Age field"},
		},
		Required: []string{"name"},
	}

	result, err := ParseInputSchema(schema)
	if err != nil {
		t.Fatalf("ParseInputSchema failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(result))
	}

	// Check that required field is marked correctly
	foundRequired := false
	for _, param := range result {
		if param.Name == "name" && param.Required {
			foundRequired = true
		}
	}
	if !foundRequired {
		t.Error("Required field 'name' not marked as required")
	}
}
