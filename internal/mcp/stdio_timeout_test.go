package mcp

import (
	"context"
	"log/slog"
	"mcpdiscord/internal/config"
	"testing"
	"time"
)

// TestStdioClient_CancelledContextRequest tests behavior when context is already cancelled
func TestStdioClient_CancelledContextRequest(t *testing.T) {
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

	// Create a context that's already cancelled
	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately

	// Try to list tools with cancelled context
	_, err := client.ListTools(cancelledCtx)
	if err == nil {
		t.Error("Expected error with cancelled context")
	}
}

// TestStdioClient_VeryShortTimeout tests behavior with very short timeout
func TestStdioClient_VeryShortTimeout(t *testing.T) {
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

	// Create context with very short timeout (may or may not trigger depending on speed)
	shortCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()

	// Try to list tools - may succeed or fail depending on timing
	_, _ = client.ListTools(shortCtx)
}

// TestStdioClient_ConcurrentRequests tests multiple concurrent requests
func TestStdioClient_ConcurrentRequests(t *testing.T) {
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

	// Launch multiple concurrent requests
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := client.ListTools(ctx)
			if err != nil {
				t.Logf("ListTools failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestStdioClient_RequestAfterContextCancel tests calling after context cancel
func TestStdioClient_RequestAfterContextCancel(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx, cancel := context.WithCancel(context.Background())
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	// Get tools successfully first
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	// Cancel the context
	cancel()

	// Try to call a tool with cancelled context
	if len(tools) > 0 {
		_, err := client.CallTool(ctx, tools[0].Name, map[string]interface{}{})
		if err == nil {
			t.Error("Expected error with cancelled context")
		}
	}
}

// TestStdioClient_RapidConnectClose tests rapid connect/close cycles
func TestStdioClient_RapidConnectClose(t *testing.T) {
	serverPath := testMCPServerPath(t)

	for i := 0; i < 5; i++ {
		cfg := config.MCPConfig{
			Command: serverPath,
		}
		client := NewStdioClient(cfg, slog.Default())

		ctx := context.Background()
		if err := client.Connect(ctx); err != nil {
			t.Fatalf("Connect %d failed: %v", i, err)
		}

		// Immediately close
		if err := client.Close(); err != nil {
			t.Fatalf("Close %d failed: %v", i, err)
		}
	}
}

// TestStdioClient_CloseWhileRequestInFlight tests closing while request is pending
func TestStdioClient_CloseWhileRequestInFlight(t *testing.T) {
	serverPath := testMCPServerPath(t)
	cfg := config.MCPConfig{
		Command: serverPath,
	}
	client := NewStdioClient(cfg, slog.Default())

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Launch request in background
	done := make(chan bool)
	go func() {
		_, _ = client.ListTools(ctx)
		done <- true
	}()

	// Wait a tiny bit then close
	time.Sleep(10 * time.Millisecond)
	client.Close()

	// Wait for goroutine to finish
	<-done
}
