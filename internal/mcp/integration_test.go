package mcp

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"mcpdiscord/internal/config"
)

func TestStdioClient_Lifecycle(t *testing.T) {
	serverPath := testMCPServerPath(t)

	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	client := NewStdioClient(cfg, logger)

	ctx := context.Background()
	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Close should work
	err = client.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestStdioClient_InvalidCommand(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "/invalid/command/that/does/not/exist",
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	client := NewStdioClient(cfg, logger)

	ctx := context.Background()
	err := client.Connect(ctx)
	if err == nil {
		t.Error("Expected error for invalid command")
		client.Close()
	}
}

func TestDiscoveryService_EmptyToolList(t *testing.T) {
	client := &MockClient{
		tools: []Tool{},
		err:   nil,
	}
	service := NewDiscoveryService(client, nil)

	ctx := context.Background()
	_, err := service.DiscoverTools(ctx)
	if err == nil {
		t.Error("Expected error for empty tool list")
	}
}

func TestDiscoveryService_InvalidTool(t *testing.T) {
	client := &MockClient{
		tools: []Tool{
			{Name: "", Description: "Invalid"},
			{Name: "valid", Description: "Valid", InputSchema: InputSchema{Type: "object"}},
		},
	}
	service := NewDiscoveryService(client, nil)

	ctx := context.Background()
	tools, err := service.DiscoverTools(ctx)
	if err != nil {
		t.Fatalf("DiscoverTools failed: %v", err)
	}

	// Should only return valid tool
	if len(tools) != 1 {
		t.Errorf("Expected 1 valid tool, got %d", len(tools))
	}
}

func TestValidateTool_Cases(t *testing.T) {
	tests := []struct {
		name        string
		tool        Tool
		shouldError bool
	}{
		{
			name: "valid tool",
			tool: Tool{
				Name:        "test",
				Description: "Test tool",
				InputSchema: InputSchema{Type: "object"},
			},
			shouldError: false,
		},
		{
			name: "missing name",
			tool: Tool{
				Name:        "",
				Description: "Test",
				InputSchema: InputSchema{Type: "object"},
			},
			shouldError: true,
		},
		{
			name: "missing schema type",
			tool: Tool{
				Name:        "test",
				Description: "Test",
				InputSchema: InputSchema{},
			},
			shouldError: true,
		},
		{
			name: "invalid schema type",
			tool: Tool{
				Name:        "test",
				Description: "Test",
				InputSchema: InputSchema{Type: "invalid"},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTool(&tt.tool)
			if tt.shouldError && err == nil {
				t.Error("Expected error")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSanitizeToolName_Internal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple", "simple"},
		{"with spaces", "with-spaces"},
		{"with_underscores", "with-underscores"},
		{"With@Special#Chars!", "withspecialchars"},
		{"123numbers", "123numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeToolName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeToolName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseInputSchema_Complex(t *testing.T) {
	schema := InputSchema{
		Type: "object",
		Properties: map[string]PropertySchema{
			"string_field":  {Type: "string", Description: "A string"},
			"number_field":  {Type: "number", Description: "A number"},
			"boolean_field": {Type: "boolean", Description: "A boolean"},
			"array_field":   {Type: "array", Description: "An array"},
			"object_field":  {Type: "object", Description: "An object"},
		},
		Required: []string{"string_field", "number_field"},
	}

	params, err := ParseInputSchema(schema)
	if err != nil {
		t.Fatalf("ParseInputSchema failed: %v", err)
	}

	if len(params) != 5 {
		t.Errorf("Expected 5 parameters, got %d", len(params))
	}

	// Check required fields
	stringParam := findParam(params, "string_field")
	if stringParam == nil || !stringParam.Required {
		t.Error("string_field should be required")
	}

	numberParam := findParam(params, "number_field")
	if numberParam == nil || !numberParam.Required {
		t.Error("number_field should be required")
	}

	// Check optional fields
	boolParam := findParam(params, "boolean_field")
	if boolParam == nil || boolParam.Required {
		t.Error("boolean_field should be optional")
	}
}

func TestCheckNameCollisions_Comprehensive(t *testing.T) {
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
			name: "space vs underscore collision",
			tools: []Tool{
				{Name: "get data"},
				{Name: "get_data"},
			},
			shouldError: true,
		},
		{
			name: "case collision",
			tools: []Tool{
				{Name: "GetData"},
				{Name: "getdata"},
			},
			shouldError: true,
		},
		{
			name: "special chars collision",
			tools: []Tool{
				{Name: "get-data"},
				{Name: "get_data"},
			},
			shouldError: true,
		},
		{
			name: "three-way collision",
			tools: []Tool{
				{Name: "Get Data"},
				{Name: "get_data"},
				{Name: "get-data"},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckNameCollisions(tt.tools)
			if tt.shouldError && err == nil {
				t.Error("Expected collision error")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTruncateDescription_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "short string",
			input:    "Short",
			maxLen:   100,
			expected: "Short",
		},
		{
			name:     "exactly 100 chars",
			input:    string(make([]byte, 100)),
			maxLen:   100,
			expected: string(make([]byte, 100)),
		},
		{
			name:     "over 100 chars",
			input:    string(make([]byte, 150)),
			maxLen:   100,
			expected: string(make([]byte, 97)) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateDescription(tt.input)
			if len(result) > 100 {
				t.Errorf("Result exceeds 100 chars: %d", len(result))
			}
		})
	}
}
