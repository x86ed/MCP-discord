package translator

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"mcpdiscord/internal/mcp"
)

// TestParameterDescription_EmptyString tests that empty string descriptions default to "No description"
func TestParameterDescription_EmptyString(t *testing.T) {
	trans := New(nil)
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"param1": {
					Type:        "string",
					Description: "", // Empty description
				},
			},
		},
	}

	options, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	if len(options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(options))
	}

	if options[0].Description != "No description" {
		t.Errorf("Expected 'No description', got %q", options[0].Description)
	}
}

// TestParameterDescription_ValidDescription tests that valid descriptions are preserved
func TestParameterDescription_ValidDescription(t *testing.T) {
	trans := New(nil)
	expectedDesc := "This is a valid parameter description"
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"param1": {
					Type:        "string",
					Description: expectedDesc,
				},
			},
		},
	}

	options, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	if len(options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(options))
	}

	if options[0].Description != expectedDesc {
		t.Errorf("Expected %q, got %q", expectedDesc, options[0].Description)
	}
}

// TestParameterDescription_ArrayWithDescription tests that array parameters append type hint
func TestParameterDescription_ArrayWithDescription(t *testing.T) {
	trans := New(nil)
	baseDesc := "List of items to process"
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"items": {
					Type:        "array",
					Description: baseDesc,
				},
			},
		},
	}

	options, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	if len(options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(options))
	}

	expectedDesc := baseDesc + " (Comma-separated list)"
	if options[0].Description != expectedDesc {
		t.Errorf("Expected %q, got %q", expectedDesc, options[0].Description)
	}
}

// TestParameterDescription_ObjectWithDescription tests that object parameters append type hint
func TestParameterDescription_ObjectWithDescription(t *testing.T) {
	trans := New(nil)
	baseDesc := "Configuration object"
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"config": {
					Type:        "object",
					Description: baseDesc,
				},
			},
		},
	}

	options, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	if len(options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(options))
	}

	expectedDesc := baseDesc + " (JSON object)"
	if options[0].Description != expectedDesc {
		t.Errorf("Expected %q, got %q", expectedDesc, options[0].Description)
	}
}

// TestParameterDescription_Truncation tests that long descriptions are truncated to 100 chars
func TestParameterDescription_Truncation(t *testing.T) {
	trans := New(nil)
	longDesc := "This is a very long description that exceeds the Discord limit of 100 characters and should be truncated with ellipsis"
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"param1": {
					Type:        "string",
					Description: longDesc,
				},
			},
		},
	}

	options, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	if len(options) != 1 {
		t.Fatalf("Expected 1 option, got %d", len(options))
	}

	desc := options[0].Description
	if len(desc) != 100 {
		t.Errorf("Expected description length 100, got %d", len(desc))
	}

	if !strings.HasSuffix(desc, "...") {
		t.Errorf("Expected description to end with '...', got %q", desc)
	}
}

// TestParameterDescription_LoggingForMissing tests that warnings are logged for missing descriptions
func TestParameterDescription_LoggingForMissing(t *testing.T) {
	// Create a buffer to capture log output
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	
	trans := New(logger)
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"param_no_desc": {
					Type:        "string",
					Description: "", // Missing description
				},
				"param_with_desc": {
					Type:        "string",
					Description: "Valid description",
				},
			},
		},
	}

	_, err := trans.parametersToOptions(tool)
	if err != nil {
		t.Fatalf("parametersToOptions failed: %v", err)
	}

	logOutput := logBuf.String()
	
	// Verify warning was logged for missing description
	if !strings.Contains(logOutput, "Parameter missing description") {
		t.Error("Expected warning log for missing description")
	}
	
	if !strings.Contains(logOutput, "test_tool") {
		t.Error("Expected tool name in log output")
	}
	
	if !strings.Contains(logOutput, "param_no_desc") {
		t.Error("Expected parameter name in log output")
	}
	
	// Verify only one warning (not for param_with_desc)
	warnCount := strings.Count(logOutput, "Parameter missing description")
	if warnCount != 1 {
		t.Errorf("Expected 1 warning, got %d", warnCount)
	}
}
