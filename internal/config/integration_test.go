// +build integration

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIntegration_ConfigLoadAndValidate is an example integration test
// that tests the full configuration loading and validation flow.
// Run with: go test -tags=integration ./...
func TestIntegration_ConfigLoadAndValidate(t *testing.T) {
	// Set up test environment variables
	os.Setenv("INTEGRATION_TEST_TOKEN", "integration-test-token-123")
	defer os.Unsetenv("INTEGRATION_TEST_TOKEN")

	// Create a temporary config file with env var interpolation
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "integration-config.json")

	configJSON := `{
		"discord": {
			"token": "${INTEGRATION_TEST_TOKEN}",
			"guildId": "integration-guild-123"
		},
		"mcp": {
			"command": "/usr/bin/node",
			"args": ["server.js"],
			"env": {
				"NODE_ENV": "test"
			},
			"transport": "stdio"
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write integration test config: %v", err)
	}

	// Test loading the configuration
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify environment variable was interpolated
	if cfg.Discord.Token != "integration-test-token-123" {
		t.Errorf("Discord token not interpolated correctly: got %q, want %q",
			cfg.Discord.Token, "integration-test-token-123")
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}

	// Verify all fields loaded correctly
	if cfg.Discord.GuildID != "integration-guild-123" {
		t.Errorf("Discord.GuildID = %q, want %q", cfg.Discord.GuildID, "integration-guild-123")
	}

	if cfg.MCP.Command != "/usr/bin/node" {
		t.Errorf("MCP.Command = %q, want %q", cfg.MCP.Command, "/usr/bin/node")
	}

	if len(cfg.MCP.Args) != 1 || cfg.MCP.Args[0] != "server.js" {
		t.Errorf("MCP.Args = %v, want [server.js]", cfg.MCP.Args)
	}

	if cfg.MCP.Env["NODE_ENV"] != "test" {
		t.Errorf("MCP.Env[NODE_ENV] = %q, want %q", cfg.MCP.Env["NODE_ENV"], "test")
	}

	if cfg.MCP.Transport != "stdio" {
		t.Errorf("MCP.Transport = %q, want %q", cfg.MCP.Transport, "stdio")
	}
}

// TestIntegration_LoadFromMultipleSources tests the priority order
// of configuration sources in an integration scenario.
func TestIntegration_LoadFromMultipleSources(t *testing.T) {
	tmpDir := t.TempDir()

	// Create default config
	defaultPath := filepath.Join(tmpDir, "default.json")
	defaultJSON := `{
		"discord": {"token": "default-token"},
		"mcp": {"command": "default-cmd"}
	}`
	if err := os.WriteFile(defaultPath, []byte(defaultJSON), 0644); err != nil {
		t.Fatalf("Failed to write default config: %v", err)
	}

	// Create env config
	envPath := filepath.Join(tmpDir, "env.json")
	envJSON := `{
		"discord": {"token": "env-token"},
		"mcp": {"command": "env-cmd"}
	}`
	if err := os.WriteFile(envPath, []byte(envJSON), 0644); err != nil {
		t.Fatalf("Failed to write env config: %v", err)
	}

	// Test 1: Only default path
	cfg, err := LoadFromSource("", "NONEXISTENT_VAR", defaultPath)
	if err != nil {
		t.Fatalf("LoadFromSource (default) failed: %v", err)
	}
	if cfg.Discord.Token != "default-token" {
		t.Errorf("Expected default-token, got %q", cfg.Discord.Token)
	}

	// Test 2: Env var takes priority
	os.Setenv("TEST_CONFIG_PATH", envPath)
	defer os.Unsetenv("TEST_CONFIG_PATH")

	cfg, err = LoadFromSource("", "TEST_CONFIG_PATH", defaultPath)
	if err != nil {
		t.Fatalf("LoadFromSource (env) failed: %v", err)
	}
	if cfg.Discord.Token != "env-token" {
		t.Errorf("Expected env-token, got %q", cfg.Discord.Token)
	}
}
