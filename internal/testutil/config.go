// Package testutil provides testing utilities for the MCP-Discord bot.
// It includes helpers for creating test configurations, mock servers,
// and mock Discord interactions.
package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"mcpdiscord/internal/config"
)

// TestConfig creates a valid test configuration with sensible defaults.
func TestConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Discord: config.DiscordConfig{
			Token:   "test-discord-token",
			GuildID: "test-guild-123",
		},
		MCP: config.MCPConfig{
			Command:   "test-command",
			Args:      []string{"arg1", "arg2"},
			Transport: "stdio",
		},
	}
}

// TestConfigWithEnv creates a test configuration with environment variables.
func TestConfigWithEnv(t *testing.T, env map[string]string) *config.Config {
	t.Helper()
	cfg := TestConfig(t)
	cfg.MCP.Env = env
	return cfg
}

// WriteTestConfig writes a configuration to a temporary file and returns the path.
func WriteTestConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	return configPath
}

// WriteTestConfigJSON writes raw JSON to a temporary config file.
func WriteTestConfigJSON(t *testing.T, jsonContent string) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test config JSON: %v", err)
	}

	return configPath
}
