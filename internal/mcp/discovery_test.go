package mcp

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	tools       []Tool
	err         error
	connectErr  error
	callToolErr error
	callResult  *ToolResult
}

func (m *MockClient) Connect(ctx context.Context) error {
	return m.connectErr
}

func (m *MockClient) ListTools(ctx context.Context) ([]Tool, error) {
	return m.tools, m.err
}

func (m *MockClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error) {
	return m.callResult, m.callToolErr
}

func (m *MockClient) Close() error {
	return nil
}

func TestDiscoveryService_DiscoverTools(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name        string
		tools       []Tool
		err         error
		shouldError bool
		expectCount int
	}{
		{
			name: "successful discovery",
			tools: []Tool{
				{
					Name:        "test_tool",
					Description: "A test tool",
					InputSchema: InputSchema{
						Type:       "object",
						Properties: map[string]PropertySchema{},
					},
				},
			},
			err:         nil,
			shouldError: false,
			expectCount: 1,
		},
		{
			name:        "empty tool list",
			tools:       []Tool{},
			err:         nil,
			shouldError: true,
			expectCount: 0,
		},
		{
			name:        "client error",
			tools:       nil,
			err:         errors.New("connection failed"),
			shouldError: true,
			expectCount: 0,
		},
		{
			name: "invalid tool skipped",
			tools: []Tool{
				{
					Name:        "",
					Description: "Invalid tool",
				},
				{
					Name:        "valid_tool",
					Description: "Valid tool",
					InputSchema: InputSchema{
						Type:       "object",
						Properties: map[string]PropertySchema{},
					},
				},
			},
			err:         nil,
			shouldError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockClient{
				tools: tt.tools,
				err:   tt.err,
			}
			service := NewDiscoveryService(client, logger)

			tools, err := service.DiscoverTools(context.Background())

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tools) != tt.expectCount {
				t.Errorf("Expected %d tools, got %d", tt.expectCount, len(tools))
			}
		})
	}
}

func TestDiscoveryService_DiscoverWithReconnect(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tools := []Tool{
		{
			Name:        "test_tool",
			Description: "A test tool",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]PropertySchema{},
			},
		},
	}

	client := &MockClient{
		tools: tools,
		err:   nil,
	}
	service := NewDiscoveryService(client, logger)

	result, err := service.DiscoverWithReconnect(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(result))
	}
}

func TestParseInputSchema(t *testing.T) {
	schema := InputSchema{
		Type: "object",
		Properties: map[string]PropertySchema{
			"name": {
				Type:        "string",
				Description: "User name",
			},
			"age": {
				Type:        "number",
				Description: "User age",
			},
			"active": {
				Type:        "boolean",
				Description: "Is active",
			},
		},
		Required: []string{"name"},
	}

	params, err := ParseInputSchema(schema)
	if err != nil {
		t.Fatalf("ParseInputSchema failed: %v", err)
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	// Check required parameter
	nameParam := findParam(params, "name")
	if nameParam == nil {
		t.Fatal("Could not find 'name' parameter")
	}
	if !nameParam.Required {
		t.Error("Expected 'name' to be required")
	}

	// Check optional parameter
	ageParam := findParam(params, "age")
	if ageParam == nil {
		t.Fatal("Could not find 'age' parameter")
	}
	if ageParam.Required {
		t.Error("Expected 'age' to be optional")
	}
}

func findParam(params []Parameter, name string) *Parameter {
	for _, p := range params {
		if p.Name == name {
			return &p
		}
	}
	return nil
}
