package config

import (
	"strings"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				Discord: DiscordConfig{
					Token:   "valid-token",
					GuildID: "guild-123",
				},
				MCP: MCPConfig{
					Command:   "npx",
					Args:      []string{"-y", "server"},
					Transport: "stdio",
				},
			},
			wantErr: false,
		},
		{
			name: "valid config without optional fields",
			config: &Config{
				Discord: DiscordConfig{
					Token: "valid-token",
				},
				MCP: MCPConfig{
					Command: "npx",
				},
			},
			wantErr: false,
		},
		{
			name: "missing discord token",
			config: &Config{
				Discord: DiscordConfig{
					Token: "",
				},
				MCP: MCPConfig{
					Command: "npx",
				},
			},
			wantErr: true,
			errMsg:  "discord.token is required",
		},
		{
			name: "missing mcp command",
			config: &Config{
				Discord: DiscordConfig{
					Token: "valid-token",
				},
				MCP: MCPConfig{
					Command: "",
				},
			},
			wantErr: true,
			errMsg:  "mcp.command is required",
		},
		{
			name: "unresolved env var in discord token",
			config: &Config{
				Discord: DiscordConfig{
					Token: "${UNRESOLVED_TOKEN}",
				},
				MCP: MCPConfig{
					Command: "npx",
				},
			},
			wantErr: true,
			errMsg:  "discord.token contains unresolved environment variable",
		},
		{
			name: "invalid transport type",
			config: &Config{
				Discord: DiscordConfig{
					Token: "valid-token",
				},
				MCP: MCPConfig{
					Command:   "npx",
					Transport: "invalid-transport",
				},
			},
			wantErr: true,
			errMsg:  "mcp.transport must be one of: stdio, sse, websocket",
		},
		{
			name: "valid transport types",
			config: &Config{
				Discord: DiscordConfig{
					Token: "valid-token",
				},
				MCP: MCPConfig{
					Command:   "npx",
					Transport: "sse",
				},
			},
			wantErr: false,
		},
		{
			name: "unresolved env var in mcp env",
			config: &Config{
				Discord: DiscordConfig{
					Token: "valid-token",
				},
				MCP: MCPConfig{
					Command: "npx",
					Env: map[string]string{
						"API_KEY": "${UNRESOLVED_KEY}",
					},
				},
			},
			wantErr: true,
			errMsg:  "mcp.env.API_KEY contains unresolved environment variable",
		},
		{
			name: "multiple validation errors",
			config: &Config{
				Discord: DiscordConfig{
					Token: "",
				},
				MCP: MCPConfig{
					Command:   "",
					Transport: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "discord.token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Validate() should return error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() should not return error, got: %v", err)
				}
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "test.field",
		Message: "test error message",
	}

	expected := "invalid configuration: test.field: test error message"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %q, want %q", err.Error(), expected)
	}
}
