package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	configJSON := `{
		"discord": {
			"token": "test-token-123",
			"guildId": "guild-456"
		},
		"mcp": {
			"command": "npx",
			"args": ["-y", "@modelcontextprotocol/server-weather"],
			"transport": "stdio"
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify Discord config
	if cfg.Discord.Token != "test-token-123" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "test-token-123")
	}
	if cfg.Discord.GuildID != "guild-456" {
		t.Errorf("Discord.GuildID = %q, want %q", cfg.Discord.GuildID, "guild-456")
	}

	// Verify MCP config
	if cfg.MCP.Command != "npx" {
		t.Errorf("MCP.Command = %q, want %q", cfg.MCP.Command, "npx")
	}
	if len(cfg.MCP.Args) != 2 {
		t.Errorf("MCP.Args length = %d, want 2", len(cfg.MCP.Args))
	}
	if cfg.MCP.Transport != "stdio" {
		t.Errorf("MCP.Transport = %q, want %q", cfg.MCP.Transport, "stdio")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Load() should fail for nonexistent file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() should fail for invalid JSON")
	}
}

func TestLoadFromSource_Priority(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test config files
	defaultConfig := filepath.Join(tmpDir, "default.json")
	envConfig := filepath.Join(tmpDir, "env.json")
	flagConfig := filepath.Join(tmpDir, "flag.json")

	testConfigs := []struct {
		path  string
		token string
	}{
		{defaultConfig, "default-token"},
		{envConfig, "env-token"},
		{flagConfig, "flag-token"},
	}

	for _, tc := range testConfigs {
		content := `{"discord":{"token":"` + tc.token + `"},"mcp":{"command":"test"}}`
		if err := os.WriteFile(tc.path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}
	}

	tests := []struct {
		name         string
		flagPath     string
		envValue     string
		defaultPath  string
		expectedPath string
		wantToken    string
	}{
		{
			name:         "flag takes priority",
			flagPath:     flagConfig,
			envValue:     envConfig,
			defaultPath:  defaultConfig,
			expectedPath: flagConfig,
			wantToken:    "flag-token",
		},
		{
			name:         "env takes priority over default",
			flagPath:     "",
			envValue:     envConfig,
			defaultPath:  defaultConfig,
			expectedPath: envConfig,
			wantToken:    "env-token",
		},
		{
			name:         "default is used when nothing else set",
			flagPath:     "",
			envValue:     "",
			defaultPath:  defaultConfig,
			expectedPath: defaultConfig,
			wantToken:    "default-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if specified
			if tt.envValue != "" {
				os.Setenv("TEST_CONFIG_PATH", tt.envValue)
				defer os.Unsetenv("TEST_CONFIG_PATH")
			} else {
				os.Unsetenv("TEST_CONFIG_PATH")
			}

			cfg, err := LoadFromSource(tt.flagPath, "TEST_CONFIG_PATH", tt.defaultPath)
			if err != nil {
				t.Fatalf("LoadFromSource() failed: %v", err)
			}

			if cfg.Discord.Token != tt.wantToken {
				t.Errorf("Token = %q, want %q", cfg.Discord.Token, tt.wantToken)
			}
		})
	}
}

func TestLoadFromSource_NoPathSpecified(t *testing.T) {
	os.Unsetenv("TEST_CONFIG_PATH")
	_, err := LoadFromSource("", "TEST_CONFIG_PATH", "")
	if err == nil {
		t.Error("LoadFromSource() should fail when no path is specified")
	}
}
