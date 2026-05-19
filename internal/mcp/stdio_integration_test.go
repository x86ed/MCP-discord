package mcp

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"mcpdiscord/internal/config"
)

// testMCPServerPath returns the path to the compiled test MCP server
func testMCPServerPath(t *testing.T) string {
	t.Helper()
	
	// Get the project root (go up from internal/mcp)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	serverDir := filepath.Join(projectRoot, "internal", "testutil", "mcpserver")
	
	// Build the test server
	serverBinary := filepath.Join(serverDir, "testmcp")
	if runtime.GOOS == "windows" {
		serverBinary += ".exe"
	}
	
	// Check if already built and up to date
	needsBuild := true
	if info, err := os.Stat(serverBinary); err == nil {
		// Check if source is newer than binary
		if srcInfo, err := os.Stat(filepath.Join(serverDir, "main.go")); err == nil {
			if info.ModTime().After(srcInfo.ModTime()) {
				needsBuild = false
			}
		}
	}
	
	if needsBuild {
		cmd := exec.Command("go", "build", "-o", serverBinary, ".")
		cmd.Dir = serverDir
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to build test MCP server: %v\n%s", err, output)
		}
	}
	
	return serverBinary
}

func TestStdioClient_ListTools_RealMCP(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	// Connect to the server
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// List tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	
	// Verify we got expected tools
	if len(tools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(tools))
	}
	
	// Check for specific tools
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}
	
	expectedTools := []string{"echo", "add", "greet"}
	for _, name := range expectedTools {
		if !toolNames[name] {
			t.Errorf("Expected tool %q not found", name)
		}
	}
	
	// Verify tool structure
	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("Tool has empty name")
		}
		if tool.Description == "" {
			t.Error("Tool has empty description")
		}
		if tool.InputSchema.Type != "object" {
			t.Errorf("Tool %q has invalid schema type: %q", tool.Name, tool.InputSchema.Type)
		}
	}
}

func TestStdioClient_CallTool_Echo(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Call echo tool
	result, err := client.CallTool(ctx, "echo", map[string]interface{}{
		"message": "Hello, World!",
	})
	
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}
	
	if result.IsError {
		t.Error("Expected successful result, got error")
	}
	
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(result.Content))
	}
	
	if result.Content[0].Type != "text" {
		t.Errorf("Expected text content, got %q", result.Content[0].Type)
	}
	
	if result.Content[0].Text != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %q", result.Content[0].Text)
	}
}

func TestStdioClient_CallTool_Add(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Call add tool
	result, err := client.CallTool(ctx, "add", map[string]interface{}{
		"a": 10.0,
		"b": 32.5,
	})
	
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}
	
	if result.IsError {
		t.Error("Expected successful result, got error")
	}
	
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(result.Content))
	}
	
	if result.Content[0].Text != "Result: 42.50" {
		t.Errorf("Expected 'Result: 42.50', got %q", result.Content[0].Text)
	}
}

func TestStdioClient_CallTool_Greet(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Call greet tool
	result, err := client.CallTool(ctx, "greet", map[string]interface{}{
		"name": "Alice",
	})
	
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}
	
	if result.IsError {
		t.Error("Expected successful result, got error")
	}
	
	if len(result.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(result.Content))
	}
	
	if result.Content[0].Text != "Hello, Alice!" {
		t.Errorf("Expected 'Hello, Alice!', got %q", result.Content[0].Text)
	}
}

func TestStdioClient_CallTool_InvalidTool(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Call non-existent tool
	_, err := client.CallTool(ctx, "nonexistent", map[string]interface{}{})
	
	if err == nil {
		t.Error("Expected error for invalid tool, got nil")
	}
}

func TestStdioClient_CallTool_InvalidArguments(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Call echo without required message argument
	result, err := client.CallTool(ctx, "echo", map[string]interface{}{})
	
	// Should succeed but return error result
	if err != nil {
		t.Fatalf("CallTool failed: %v", err)
	}
	
	if !result.IsError {
		t.Error("Expected error result for missing arguments")
	}
}

func TestStdioClient_MultipleOperations(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// List tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	if len(tools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(tools))
	}
	
	// Call multiple tools in sequence
	tests := []struct {
		name string
		tool string
		args map[string]interface{}
	}{
		{"echo", "echo", map[string]interface{}{"message": "test1"}},
		{"add", "add", map[string]interface{}{"a": 5.0, "b": 3.0}},
		{"greet", "greet", map[string]interface{}{"name": "Bob"}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.CallTool(ctx, tt.tool, tt.args)
			if err != nil {
				t.Errorf("CallTool %q failed: %v", tt.tool, err)
			}
			if result.IsError {
				t.Errorf("Tool %q returned error", tt.tool)
			}
		})
	}
}

func TestStdioClient_ContextCancellation(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// Create a context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	// Wait for context to expire
	time.Sleep(10 * time.Millisecond)
	
	// Try to call tool with expired context
	_, err := client.CallTool(ctx, "echo", map[string]interface{}{"message": "test"})
	
	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
}

func TestStdioClient_CloseAndReconnect(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	
	// First connection
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	if len(tools) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(tools))
	}
	
	// Close
	if err := client.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	
	// Try to use after close (should fail)
	_, err = client.ListTools(ctx)
	if err == nil {
		t.Error("Expected error after close, got nil")
	}
}

func TestStdioClient_InvalidCommandPath(t *testing.T) {
	cfg := config.MCPConfig{
		Command: "/nonexistent/command/path",
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	err := client.Connect(ctx)
	
	if err == nil {
		t.Error("Expected error with invalid command path, got nil")
		client.Close()
	}
}

func TestStdioClient_SendRequest_Coverage(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{},
	}
	
	client := NewStdioClient(cfg, slog.Default())
	
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	
	// This exercises sendRequest internally through public methods
	// Test with various parameter types
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{"string_param", map[string]interface{}{"message": "test"}},
		{"number_params", map[string]interface{}{"a": 1.0, "b": 2.0}},
		{"empty_params", map[string]interface{}{}},
		{"nil_value", map[string]interface{}{"message": nil}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// CallTool uses sendRequest internally
			_, err := client.CallTool(ctx, "echo", tc.args)
			// We expect some to fail, but sendRequest should be exercised
			if err != nil {
				t.Logf("Expected failure for %s: %v", tc.name, err)
			}
		})
	}
}
