package mcp

import (
	"context"
	"log/slog"
	"testing"

	"mcpdiscord/internal/config"
)

func TestStdioClient_AlreadyConnected(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("First connect failed: %v", err)
	}
	defer client.Close()

	// Try to connect again
	err := client.Connect(ctx)
	if err == nil {
		t.Error("Expected error when connecting twice")
	}
	if err.Error() != "client already connected" {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestStdioClient_ConnectWithEnv(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Env: map[string]string{
			"TEST_VAR": "test_value",
			"CUSTOM":   "env_var",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with env failed: %v", err)
	}
	defer client.Close()

	// Verify client is connected
	if client.cmd == nil {
		t.Error("Command should not be nil after connect")
	}
}

func TestStdioClient_InvalidCommandError(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "/nonexistent/command/that/does/not/exist",
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	err := client.Connect(ctx)
	if err == nil {
		t.Error("Expected error for invalid command")
		client.Close()
	}
}

func TestStdioClient_ListToolsError(t *testing.T) {
	// Create a client without connecting
	cfg := config.MCPConfig{
		Command: "echo",
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	_, err := client.ListTools(ctx)
	if err == nil {
		t.Error("Expected error when calling ListTools without connection")
	}
}

func TestStdioClient_CallToolError(t *testing.T) {
	// Create a client without connecting
	cfg := config.MCPConfig{
		Command: "echo",
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	_, err := client.CallTool(ctx, "test", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error when calling CallTool without connection")
	}
}

func TestStdioClient_DoubleClose(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Close once
	if err := client.Close(); err != nil {
		t.Errorf("First close failed: %v", err)
	}

	// Close again - should not error
	if err := client.Close(); err != nil {
		t.Errorf("Second close should not error: %v", err)
	}
}

func TestStdioClient_CallToolWithResult(t *testing.T) {
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

	// Call a tool that exists
	result, err := client.CallTool(ctx, "echo", map[string]interface{}{
		"message": "test",
	})
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
	if result.IsError {
		t.Error("Result should not be an error")
	}
}

func TestStdioClient_CallToolWithError(t *testing.T) {
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

	// Call with invalid arguments
	result, err := client.CallTool(ctx, "echo", map[string]interface{}{
		"a": 1,
		"b": 2,
	})
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if !result.IsError {
		t.Error("Expected result to be marked as error")
	}
}

func TestStdioClient_CallToolMissingArgs(t *testing.T) {
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

	// Call without required arguments
	result, err := client.CallTool(ctx, "echo", nil)
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if !result.IsError {
		t.Error("Expected result to be marked as error for missing args")
	}
}

func TestStdioClient_CallToolNullArg(t *testing.T) {
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

	// Call with null argument
	result, err := client.CallTool(ctx, "echo", map[string]interface{}{
		"message": nil,
	})
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}

	if !result.IsError {
		t.Error("Expected result to be marked as error for null arg")
	}
}

func TestNewStdioClient_NilLogger(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "echo",
	}
	client := NewStdioClient(cfg, nil)

	if client.logger == nil {
		t.Error("Logger should default to slog.Default() when nil")
	}
}

func TestStdioClient_CloseWithoutConnect(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "echo",
	}
	client := NewStdioClient(cfg, slog.Default())

	// Close without connecting should not error
	if err := client.Close(); err != nil {
		t.Errorf("Close without connect should not error: %v", err)
	}
}

// TestStdioClient_MultipleToolCalls tests calling tools multiple times
func TestStdioClient_MultipleToolCalls(t *testing.T) {
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

	// Call ListTools multiple times to exercise sendRequest/readResponses
	for i := 0; i < 5; i++ {
		tools, err := client.ListTools(ctx)
		if err != nil {
			t.Fatalf("ListTools call %d failed: %v", i+1, err)
		}
		if len(tools) == 0 {
			t.Error("Expected at least one tool")
		}
	}

	// Call a tool multiple times if available
	tools, _ := client.ListTools(ctx)
	if len(tools) > 0 {
		for i := 0; i < 3; i++ {
			_, err := client.CallTool(ctx, tools[0].Name, map[string]interface{}{
				"iteration": i,
			})
			// Ignore errors as tool may not accept these args
			_ = err
		}
	}
}

// TestStdioClient_ConnectWithMultipleEnvVars tests env var handling
func TestStdioClient_ConnectWithMultipleEnvVars(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
		Env: map[string]string{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"VAR3":      "value3",
			"EMPTY_VAR": "",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect with multiple env vars failed: %v", err)
	}
	defer client.Close()

	// Verify connection works
	_, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
}

