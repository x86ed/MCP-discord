package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"os/exec"
	"testing"
	"time"
)

// TestStdioClient_CloseWithSlowServer tests Close with a server that doesn't exit quickly
func TestStdioClient_CloseWithSlowServer(t *testing.T) {
	// Use 'sleep' command which will take time to exit
	cfg := config.MCPConfig{
		Command: "sleep",
		Args:    []string{"10"}, // Sleep for 10 seconds
	}
	client := NewStdioClient(cfg, slog.Default())

	// Manually set up a basic process without full connect
	ctx := context.Background()
	client.cmd = exec.CommandContext(ctx, cfg.Command, cfg.Args...)
	
	// Create pipes manually
	stdin, err := client.cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	client.stdin = stdin

	// Start the process
	if err := client.cmd.Start(); err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}

	// Close should handle the slow exit (will hit 5 second timeout and kill)
	start := time.Now()
	if err := client.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	elapsed := time.Since(start)

	// Should have timed out and killed around 5 seconds
	if elapsed < 4*time.Second || elapsed > 7*time.Second {
		t.Logf("Close took %v (expected ~5s for timeout and kill)", elapsed)
	}
}

// TestStdioClient_MultipleCloseNoError tests multiple Close calls don't error
func TestStdioClient_MultipleCloseNoError(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// First close
	if err := client.Close(); err != nil {
		t.Errorf("First close failed: %v", err)
	}

	// Second close should be no-op (cmd is nil)
	if err := client.Close(); err != nil {
		t.Errorf("Second close failed: %v", err)
	}

	// Third close for good measure
	if err := client.Close(); err != nil {
		t.Errorf("Third close failed: %v", err)
	}
}

// TestStdioClient_ConnectInitializeFlow tests the full connect -> initialize flow
func TestStdioClient_ConnectInitializeFlow(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{"--test-mode"},
		Env: map[string]string{
			"MCP_TEST": "true",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	
	// Connect should call initialize internally
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect (with initialize) failed: %v", err)
	}
	defer client.Close()

	// Verify we can list tools (proves initialization worked)
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools after connect/initialize failed: %v", err)
	}

	if len(tools) == 0 {
		t.Error("Expected at least one tool after successful initialization")
	}
}

// TestStdioClient_StdinNilCheck tests the stdin nil check in sendRequest
func TestStdioClient_StdinNilCheck(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "echo",
	}
	client := NewStdioClient(cfg, slog.Default())

	// Don't connect, try to call tool directly
	ctx := context.Background()
	_, err := client.CallTool(ctx, "test", map[string]interface{}{})
	
	if err == nil {
		t.Error("Expected error when stdin is nil")
	}
	
	if err != nil && err.Error() != "client not connected" {
		t.Logf("Got error: %v", err)
	}
}

// TestStdioClient_LargeNumberOfRequests tests handling many sequential requests
func TestStdioClient_LargeNumberOfRequests(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	// Make many requests to exercise sendRequest/readResponses
	for i := 0; i < 20; i++ {
		_, err := client.ListTools(ctx)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
	}
}
