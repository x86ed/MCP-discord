package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"testing"
)

// TestStdioClient_FullLifecycle tests a complete client lifecycle
func TestStdioClient_FullLifecycle(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Env: map[string]string{
			"TEST_MODE": "true",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	
	// Connect
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// List tools multiple times to exercise readResponses
	for i := 0; i < 3; i++ {
		tools, err := client.ListTools(ctx)
		if err != nil {
			t.Fatalf("ListTools call %d failed: %v", i+1, err)
		}
		if len(tools) == 0 {
			t.Error("Expected at least one tool")
		}
	}

	// Call a tool if one exists
	tools, _ := client.ListTools(ctx)
	if len(tools) > 0 {
		// Try calling with empty args
		_, err := client.CallTool(ctx, tools[0].Name, map[string]interface{}{})
		// Result may succeed or fail depending on the tool's requirements
		_ = err // Ignore error, we just want to exercise the code path
		
		// Try calling with some args
		_, err = client.CallTool(ctx, tools[0].Name, map[string]interface{}{
			"test_param": "test_value",
		})
		_ = err // Ignore error
	}

	// Close
	if err := client.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Try operations after close (should fail)
	_, err := client.ListTools(ctx)
	if err == nil {
		t.Error("Expected error when calling ListTools after Close")
	}

	_, err = client.CallTool(ctx, "test", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error when calling CallTool after Close")
	}
}

// TestStdioClient_ConnectWithComplexEnv tests environment variable handling
func TestStdioClient_ConnectWithComplexEnv(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Env: map[string]string{
			"VAR1":        "value1",
			"VAR2":        "value2",
			"EMPTY_VAR":   "",
			"SPECIAL_VAR": "value with spaces",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with complex env failed: %v", err)
	}
	defer client.Close()

	// Verify connection works
	_, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
}

// TestStdioClient_ConnectWithArgs tests command with arguments
func TestStdioClient_ConnectWithArgs(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{"arg1", "arg2", "--flag"},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with args failed: %v", err)
	}
	defer client.Close()

	// Verify connection works
	_, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
}
