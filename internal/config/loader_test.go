package config

import (
	"os"
	"path/filepath"
	"strings"
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

func TestLoadFromJSON(t *testing.T) {
	raw := `{
		"discord": {
			"token": "json-token"
		},
		"mcp": {
			"command": "json-cmd",
			"transport": "stdio"
		}
	}`

	cfg, err := LoadFromJSON(raw)
	if err != nil {
		t.Fatalf("LoadFromJSON() failed: %v", err)
	}

	if cfg.Discord.Token != "json-token" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "json-token")
	}
	if cfg.MCP.Command != "json-cmd" {
		t.Errorf("MCP.Command = %q, want %q", cfg.MCP.Command, "json-cmd")
	}
}

func TestLoadFromJSON_Empty(t *testing.T) {
	_, err := LoadFromJSON("   ")
	if err == nil {
		t.Error("LoadFromJSON() should fail for empty input")
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

func TestLoadFromJSONOrSource_Priority(t *testing.T) {
	tmpDir := t.TempDir()

	defaultConfig := filepath.Join(tmpDir, "default.json")
	if err := os.WriteFile(defaultConfig, []byte(`{"discord":{"token":"default-token"},"mcp":{"command":"default-cmd"}}`), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	t.Setenv("TEST_JSON_CONFIG", `{"discord":{"token":"json-token"},"mcp":{"command":"json-cmd"}}`)
	t.Setenv("TEST_CONFIG_PATH", defaultConfig)

	cfg, err := LoadFromJSONOrSource("", "TEST_CONFIG_PATH", "TEST_JSON_CONFIG", defaultConfig)
	if err != nil {
		t.Fatalf("LoadFromJSONOrSource() failed: %v", err)
	}

	if cfg.Discord.Token != "json-token" {
		t.Errorf("Token = %q, want %q", cfg.Discord.Token, "json-token")
	}
	if cfg.MCP.Command != "json-cmd" {
		t.Errorf("Command = %q, want %q", cfg.MCP.Command, "json-cmd")
	}
}

func TestInterpolateEnvVars_Success(t *testing.T) {
	os.Setenv("TEST_TOKEN", "secret-token-123")
	os.Setenv("TEST_GUILD", "guild-789")
	defer os.Unsetenv("TEST_TOKEN")
	defer os.Unsetenv("TEST_GUILD")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	configJSON := `{
		"discord": {
			"token": "${TEST_TOKEN}",
			"guildId": "${TEST_GUILD}"
		},
		"mcp": {
			"command": "test"
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Discord.Token != "secret-token-123" {
		t.Errorf("Token = %q, want %q", cfg.Discord.Token, "secret-token-123")
	}
	if cfg.Discord.GuildID != "guild-789" {
		t.Errorf("GuildID = %q, want %q", cfg.Discord.GuildID, "guild-789")
	}
}

func TestInterpolateEnvVars_MissingVar(t *testing.T) {
	os.Unsetenv("NONEXISTENT_VAR")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	configJSON := `{
		"discord": {
			"token": "${NONEXISTENT_VAR}"
		},
		"mcp": {
			"command": "test"
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// The variable should remain unresolved
	if !strings.Contains(cfg.Discord.Token, "${NONEXISTENT_VAR}") {
		t.Errorf("Expected unresolved env var, got: %s", cfg.Discord.Token)
	}

	// Validation should catch it
	if err := cfg.Validate(); err == nil {
		t.Error("Validate() should fail for unresolved env var")
	}
}

func TestLoad_ValidationFailure(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.json")

	// Config with empty token
	configJSON := `{
		"discord": {
			"token": ""
		},
		"mcp": {
			"command": "test",
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

	// Validation should fail
	err = cfg.Validate()
	if err == nil {
		t.Error("Validate() should fail when token is empty")
	}
	if err != nil && !strings.Contains(err.Error(), "discord.token is required") {
		t.Errorf("Expected token validation error, got: %v", err)
	}
}
