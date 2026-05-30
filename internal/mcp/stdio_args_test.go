package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"testing"
)

// TestStdioClient_ConnectWithBothArgsAndEnv tests with both args and env
func TestStdioClient_ConnectWithBothArgsAndEnv(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{"--mode", "test"},
		Env: map[string]string{
			"TEST_ENV": "test_value",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with args and env failed: %v", err)
	}
	defer client.Close()

	// Verify connection works
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(tools) == 0 {
		t.Error("Expected at least one tool")
	}

	// Call multiple tools if available
	for i := 0; i < len(tools) && i < 3; i++ {
		_, err := client.CallTool(ctx, tools[i].Name, map[string]interface{}{})
		// Ignore errors as tools may require args
		_ = err
	}
}

// TestStdioClient_EmptyEnvMap tests with empty env map (should use default env)
func TestStdioClient_EmptyEnvMap(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Env:     map[string]string{}, // Empty map, should not set env vars
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with empty env map failed: %v", err)
	}
	defer client.Close()

	// Verify connection works
	_, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
}
