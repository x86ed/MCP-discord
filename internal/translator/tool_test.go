package translator

import (
	"testing"

	"mcpdiscord/internal/mcp"

	"github.com/bwmarrin/discordgo"
)

func TestToolToSlashCommand_Comprehensive(t *testing.T) {
	trans := New(nil)

	tests := []struct {
		name        string
		tool        mcp.Tool
		shouldError bool
		checkName   string
		checkOpts   int
	}{
		{
			name: "simple string parameter",
			tool: mcp.Tool{
				Name:        "echo",
				Description: "Echo a message",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"message": {Type: "string", Description: "The message"},
					},
				},
			},
			shouldError: false,
			checkName:   "echo",
			checkOpts:   1,
		},
		{
			name: "multiple parameters with different types",
			tool: mcp.Tool{
				Name:        "create_user",
				Description: "Create a new user",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"name":    {Type: "string", Description: "User name"},
						"age":     {Type: "number", Description: "User age"},
						"active":  {Type: "boolean", Description: "Is active"},
						"tags":    {Type: "array", Description: "Tags"},
						"profile": {Type: "object", Description: "Profile data"},
					},
					Required: []string{"name"},
				},
			},
			shouldError: false,
			checkName:   "create-user",
			checkOpts:   5,
		},
		{
			name: "tool name sanitization",
			tool: mcp.Tool{
				Name:        "Get Weather Data",
				Description: "Get weather",
				InputSchema: mcp.InputSchema{
					Type:       "object",
					Properties: map[string]mcp.PropertySchema{},
				},
			},
			shouldError: false,
			checkName:   "get-weather-data",
			checkOpts:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := trans.ToolToSlashCommand(tt.tool)
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if cmd.Name != tt.checkName {
				t.Errorf("Command name = %q, want %q", cmd.Name, tt.checkName)
			}
			if len(cmd.Options) != tt.checkOpts {
				t.Errorf("Options count = %d, want %d", len(cmd.Options), tt.checkOpts)
			}
		})
	}
}

func TestTranslateArguments_Comprehensive(t *testing.T) {
	trans := New(nil)

	tests := []struct {
		name        string
		tool        mcp.Tool
		options     []*discordgo.ApplicationCommandInteractionDataOption
		shouldError bool
		checkResult func(map[string]interface{}) bool
	}{
		{
			name: "string parameter",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"message": {Type: "string"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "message", Type: discordgo.ApplicationCommandOptionString, Value: "hello"},
			},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				return args["message"] == "hello"
			},
		},
		{
			name: "number parameter",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"count": {Type: "number"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "count", Type: discordgo.ApplicationCommandOptionNumber, Value: 42.0},
			},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				return args["count"] == 42.0
			},
		},
		{
			name: "boolean parameter",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"enabled": {Type: "boolean"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "enabled", Type: discordgo.ApplicationCommandOptionBoolean, Value: true},
			},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				return args["enabled"] == true
			},
		},
		{
			name: "array parameter with CSV",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"items": {Type: "array"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "items", Type: discordgo.ApplicationCommandOptionString, Value: "a, b, c"},
			},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				arr, ok := args["items"].([]string)
				return ok && len(arr) == 3 && arr[0] == "a"
			},
		},
		{
			name: "object parameter with JSON",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"data": {Type: "object"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "data", Type: discordgo.ApplicationCommandOptionString, Value: `{"key":"value"}`},
			},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				obj, ok := args["data"].(map[string]interface{})
				return ok && obj["key"] == "value"
			},
		},
		{
			name: "invalid JSON for object",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"data": {Type: "object"},
					},
				},
			},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{Name: "data", Type: discordgo.ApplicationCommandOptionString, Value: "{invalid}"},
			},
			shouldError: true,
		},
		{
			name: "missing required parameter",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"required_field": {Type: "string"},
					},
					Required: []string{"required_field"},
				},
			},
			options:     []*discordgo.ApplicationCommandInteractionDataOption{},
			shouldError: true,
		},
		{
			name: "optional parameter omitted",
			tool: mcp.Tool{
				Name: "test",
				InputSchema: mcp.InputSchema{
					Type: "object",
					Properties: map[string]mcp.PropertySchema{
						"optional_field": {Type: "string"},
					},
				},
			},
			options:     []*discordgo.ApplicationCommandInteractionDataOption{},
			shouldError: false,
			checkResult: func(args map[string]interface{}) bool {
				_, exists := args["optional_field"]
				return !exists
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := trans.TranslateArguments(tt.tool, tt.options)
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkResult != nil && !tt.checkResult(args) {
				t.Errorf("Result check failed: %+v", args)
			}
		})
	}
}

func TestParametersToOptions_EdgeCases(t *testing.T) {
	trans := New(nil)

	// Test with 26 parameters (exceeds Discord limit of 25)
	tool := mcp.Tool{
		Name:        "complex",
		Description: "Complex tool",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: make(map[string]mcp.PropertySchema),
		},
	}

	for i := 0; i < 26; i++ {
		key := string(rune('a' + i))
		tool.InputSchema.Properties[key] = mcp.PropertySchema{
			Type:        "string",
			Description: "Param " + key,
		}
	}

	_, err := trans.ToolToSlashCommand(tool)
	if err == nil {
		t.Error("Expected error for too many parameters")
	}
}
