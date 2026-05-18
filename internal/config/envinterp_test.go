package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_EnvVarInterpolation(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_DISCORD_TOKEN", "secret-discord-token")
	os.Setenv("TEST_API_KEY", "secret-api-key")
	defer func() {
		os.Unsetenv("TEST_DISCORD_TOKEN")
		os.Unsetenv("TEST_API_KEY")
	}()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"discord": {
			"token": "${TEST_DISCORD_TOKEN}",
			"guildId": "static-guild-id"
		},
		"mcp": {
			"command": "npx",
			"args": ["-y", "@modelcontextprotocol/server-weather"],
			"env": {
				"API_KEY": "${TEST_API_KEY}",
				"STATIC_VAR": "static-value"
			},
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

	// Verify interpolated values
	if cfg.Discord.Token != "secret-discord-token" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "secret-discord-token")
	}

	if cfg.Discord.GuildID != "static-guild-id" {
		t.Errorf("Discord.GuildID = %q, want %q", cfg.Discord.GuildID, "static-guild-id")
	}

	if cfg.MCP.Env["API_KEY"] != "secret-api-key" {
		t.Errorf("MCP.Env[API_KEY] = %q, want %q", cfg.MCP.Env["API_KEY"], "secret-api-key")
	}

	if cfg.MCP.Env["STATIC_VAR"] != "static-value" {
		t.Errorf("MCP.Env[STATIC_VAR] = %q, want %q", cfg.MCP.Env["STATIC_VAR"], "static-value")
	}
}

func TestLoad_EnvVarNotSet(t *testing.T) {
	// Ensure the variable is not set
	os.Unsetenv("NONEXISTENT_VAR")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

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

	// Should keep the original ${...} syntax when var doesn't exist
	if cfg.Discord.Token != "${NONEXISTENT_VAR}" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "${NONEXISTENT_VAR}")
	}
}

func TestLoad_MultipleEnvVarsInString(t *testing.T) {
	os.Setenv("TEST_HOST", "localhost")
	os.Setenv("TEST_PORT", "8080")
	defer func() {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_PORT")
	}()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"discord": {
			"token": "http://${TEST_HOST}:${TEST_PORT}/webhook"
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

	expected := "http://localhost:8080/webhook"
	if cfg.Discord.Token != expected {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, expected)
	}
}

func TestInterpolateEnvVars(t *testing.T) {
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "simple string with env var",
			input:    "${TEST_VAR}",
			expected: "test-value",
		},
		{
			name:     "string without env var",
			input:    "plain-string",
			expected: "plain-string",
		},
		{
			name:     "string with missing env var",
			input:    "${MISSING_VAR}",
			expected: "${MISSING_VAR}",
		},
		{
			name:     "number",
			input:    42,
			expected: 42,
		},
		{
			name:     "boolean",
			input:    true,
			expected: true,
		},
		{
			name: "map with env vars",
			input: map[string]interface{}{
				"key1": "${TEST_VAR}",
				"key2": "static",
			},
			expected: map[string]interface{}{
				"key1": "test-value",
				"key2": "static",
			},
		},
		{
			name:     "array with env vars",
			input:    []interface{}{"${TEST_VAR}", "static"},
			expected: []interface{}{"test-value", "static"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateEnvVars(tt.input)
			
			// Compare based on type
			switch expected := tt.expected.(type) {
			case string:
				if result != expected {
					t.Errorf("interpolateEnvVars() = %v, want %v", result, expected)
				}
			case map[string]interface{}:
				resultMap := result.(map[string]interface{})
				for k, v := range expected {
					if resultMap[k] != v {
						t.Errorf("interpolateEnvVars()[%q] = %v, want %v", k, resultMap[k], v)
					}
				}
			case []interface{}:
				resultSlice := result.([]interface{})
				for i, v := range expected {
					if resultSlice[i] != v {
						t.Errorf("interpolateEnvVars()[%d] = %v, want %v", i, resultSlice[i], v)
					}
				}
			default:
				if result != expected {
					t.Errorf("interpolateEnvVars() = %v, want %v", result, expected)
				}
			}
		})
	}
}
