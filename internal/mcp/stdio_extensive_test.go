package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"testing"
	"time"
)

// TestStdioClient_InitializeAndNotification tests the full initialize flow
func TestStdioClient_InitializeAndNotification(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	
	// Connect calls initialize internally
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect (with initialize) failed: %v", err)
	}
	defer client.Close()

	// Verify initialization succeeded by making a request
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools after initialization failed: %v", err)
	}

	if len(tools) == 0 {
		t.Log("Warning: No tools returned, but initialization succeeded")
	}

	// Call a tool if available
	if len(tools) > 0 {
		_, _ = client.CallTool(ctx, tools[0].Name, map[string]interface{}{})
	}
}

// TestStdioClient_RepeatedOperations tests many repeated operations
func TestStdioClient_RepeatedOperations(t *testing.T) {
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

	// Do many operations to exercise different code paths
	for i := 0; i < 15; i++ {
		tools, err := client.ListTools(ctx)
		if err != nil {
			t.Fatalf("ListTools iteration %d failed: %v", i, err)
		}

		// Try calling tools with different argument patterns
		if len(tools) > 0 {
			// Vary the arguments
			args := map[string]interface{}{
				"iteration": i,
				"timestamp": time.Now().Unix(),
			}
			_, _ = client.CallTool(ctx, tools[0].Name, args)
		}
	}
}

// TestStdioClient_ContextVariations tests different context scenarios
func TestStdioClient_ContextVariations(t *testing.T) {
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

	// Test 1: Normal context
	_, err := client.ListTools(ctx)
	if err != nil {
		t.Errorf("Normal context failed: %v", err)
	}

	// Test 2: Context with timeout (long enough to succeed)
	ctx2, cancel2 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel2()
	_, err = client.ListTools(ctx2)
	if err != nil {
		t.Errorf("Context with timeout failed: %v", err)
	}

	// Test 3: Create and immediately cancel context
	ctx3, cancel3 := context.WithCancel(ctx)
	cancel3()
	_, err = client.ListTools(ctx3)
	// Should fail with cancelled error
	if err == nil {
		t.Error("Expected error with immediately cancelled context")
	}

	// Test 4: Very short timeout (may race)
	ctx4, cancel4 := context.WithTimeout(ctx, 1*time.Microsecond)
	defer cancel4()
	time.Sleep(5 * time.Millisecond) // Ensure timeout fires
	_, _ = client.ListTools(ctx4)
}

// TestStdioClient_ExtensiveToolCalls tests calling various tools extensively
func TestStdioClient_ExtensiveToolCalls(t *testing.T) {
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

	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(tools) == 0 {
		t.Skip("No tools available for testing")
	}

	// Call each tool multiple times with different args
	for _, tool := range tools {
		for i := 0; i < 3; i++ {
			args := map[string]interface{}{
				"test_param": i,
				"tool_name":  tool.Name,
			}
			_, _ = client.CallTool(ctx, tool.Name, args)
		}
	}

	// Call a non-existent tool to exercise error path
	_, err = client.CallTool(ctx, "nonexistent_tool_12345", map[string]interface{}{})
	// This should fail, which is expected
	_ = err
}

// TestStdioClient_MixedOperations tests mixing ListTools and CallTool
func TestStdioClient_MixedOperations(t *testing.T) {
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

	// Interleave ListTools and CallTool calls
	for i := 0; i < 10; i++ {
		tools, err := client.ListTools(ctx)
		if err != nil {
			t.Fatalf("ListTools failed: %v", err)
		}

		if len(tools) > 0 {
			// Call first tool
			_, _ = client.CallTool(ctx, tools[0].Name, map[string]interface{}{
				"iteration": i,
			})

			// List tools again
			_, _ = client.ListTools(ctx)

			// Call another tool if available
			if len(tools) > 1 {
				_, _ = client.CallTool(ctx, tools[1].Name, map[string]interface{}{})
			}
		}
	}
}
