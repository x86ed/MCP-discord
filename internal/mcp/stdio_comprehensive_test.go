package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"testing"
	"time"
)

// TestStdioClient_ComprehensiveWorkflow tests a realistic workflow
func TestStdioClient_ComprehensiveWorkflow(t *testing.T) {
	serverPath := testMCPServerPath(t)
	
	// Test 1: Connect with environment variables and args
	cfg := config.MCPConfig{
		Command: serverPath,
		Args:    []string{"--verbose"},
		Env: map[string]string{
			"LOG_LEVEL": "debug",
			"TEST_MODE": "true",
		},
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Test 2: List tools immediately after connect
	tools1, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("First ListTools failed: %v", err)
	}

	// Test 3: List tools again (exercises caching/repeat calls)
	tools2, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("Second ListTools failed: %v", err)
	}

	if len(tools1) != len(tools2) {
		t.Error("Tool list should be consistent")
	}

	// Test 4: Call tools with various argument patterns
	if len(tools1) > 0 {
		// Empty args
		_, _ = client.CallTool(ctx, tools1[0].Name, map[string]interface{}{})
		
		// String arg
		_, _ = client.CallTool(ctx, tools1[0].Name, map[string]interface{}{
			"arg1": "value1",
		})
		
		// Multiple args
		_, _ = client.CallTool(ctx, tools1[0].Name, map[string]interface{}{
			"arg1": "value1",
			"arg2": 42,
			"arg3": true,
		})
		
		// Nested object
		_, _ = client.CallTool(ctx, tools1[0].Name, map[string]interface{}{
			"nested": map[string]interface{}{
				"key": "value",
			},
		})
	}

	// Test 5: Call ListTools with short timeout context
	ctxShort, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.ListTools(ctxShort)
	if err != nil {
		t.Fatalf("ListTools with timeout failed: %v", err)
	}

	// Test 6: Multiple rapid sequential calls
	for i := 0; i < 10; i++ {
		_, _ = client.ListTools(ctx)
	}

	// Test 7: Close and verify cleanup
	if err := client.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Test 8: Verify operations fail after close
	_, err = client.ListTools(ctx)
	if err == nil {
		t.Error("Expected error when calling ListTools after Close")
	}

	_, err = client.CallTool(ctx, "test", nil)
	if err == nil {
		t.Error("Expected error when calling CallTool after Close")
	}
}

// TestStdioClient_EdgeCases tests edge cases and boundary conditions
func TestStdioClient_EdgeCases(t *testing.T) {
	serverPath := testMCPServerPath(t)

	// Test with nil args map
	t.Run("nil_args", func(t *testing.T) {
		cfg := config.MCPConfig{
			Command: serverPath,
		}
		client := NewStdioClient(cfg, slog.Default())

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			t.Fatalf("Connect failed: %v", err)
		}
		defer client.Close()

		tools, _ := client.ListTools(ctx)
		if len(tools) > 0 {
			_, _ = client.CallTool(ctx, tools[0].Name, nil)
		}
	})

	// Test with empty string args
	t.Run("empty_string_args", func(t *testing.T) {
		cfg := config.MCPConfig{
			Command: serverPath,
		}
		client := NewStdioClient(cfg, slog.Default())

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			t.Fatalf("Connect failed: %v", err)
		}
		defer client.Close()

		tools, _ := client.ListTools(ctx)
		if len(tools) > 0 {
			_, _ = client.CallTool(ctx, tools[0].Name, map[string]interface{}{
				"": "",
			})
		}
	})

	// Test close immediately after connect
	t.Run("immediate_close", func(t *testing.T) {
		cfg := config.MCPConfig{
			Command: serverPath,
		}
		client := NewStdioClient(cfg, slog.Default())

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			t.Fatalf("Connect failed: %v", err)
		}

		if err := client.Close(); err != nil {
			t.Fatalf("Immediate close failed: %v", err)
		}
	})

	// Test double close
	t.Run("double_close", func(t *testing.T) {
		cfg := config.MCPConfig{
			Command: serverPath,
		}
		client := NewStdioClient(cfg, slog.Default())

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			t.Fatalf("Connect failed: %v", err)
		}

		client.Close()
		// Second close should not error
		if err := client.Close(); err != nil {
			t.Errorf("Second close should not error: %v", err)
		}
	})
}
