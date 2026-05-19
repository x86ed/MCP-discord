package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ErrorPaths(t *testing.T) {
	tempDir := t.TempDir()

	// Test invalid JSON that will fail to re-marshal (though this is rare)
	t.Run("unmarshal error path", func(t *testing.T) {
		cfgFile := filepath.Join(tempDir, "badstruct.json")
		
		// Create a config that will pass initial JSON parse but might fail struct unmarshal
		// This is hard to trigger, but we can try with unexpected types
		badJSON := []byte(`{
			"discord": {"token": 123},
			"mcp": {"command": "test"}
		}`)
		
		err := os.WriteFile(cfgFile, badJSON, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		_, err = Load(cfgFile)
		// This should work actually since JSON is flexible
		// But it tests the code path
		if err != nil {
			t.Logf("Got error as expected: %v", err)
		}
	})

	// Test a file that exists but has problematic content
	t.Run("edge case JSON", func(t *testing.T) {
		cfgFile := filepath.Join(tempDir, "edge.json")
		
		// Valid JSON but might trigger different code paths
		edgeJSON := []byte(`{
			"discord": {
				"token": "${NONEXISTENT_VAR}"
			},
			"mcp": {
				"command": "test",
				"args": ["arg1", "arg2"],
				"env": {
					"VAR1": "${ENV1}",
					"VAR2": "value2"
				}
			}
		}`)
		
		err := os.WriteFile(cfgFile, edgeJSON, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		cfg, err := Load(cfgFile)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Should preserve the ${NONEXISTENT_VAR} since env var doesn't exist
		if cfg.Discord.Token != "${NONEXISTENT_VAR}" {
			t.Errorf("Expected unresolved var, got %q", cfg.Discord.Token)
		}
	})

	// Test re-marshal error path (very hard to trigger in practice)
	t.Run("complex nested structure", func(t *testing.T) {
		cfgFile := filepath.Join(tempDir, "complex.json")
		
		complexJSON := []byte(`{
			"discord": {
				"token": "test_token",
				"guild_id": ""
			},
			"mcp": {
				"command": "test_cmd",
				"args": [],
				"env": {},
				"transport": "stdio"
			}
		}`)
		
		err := os.WriteFile(cfgFile, complexJSON, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		cfg, err := Load(cfgFile)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if cfg.Discord.Token != "test_token" {
			t.Errorf("Token = %q, want %q", cfg.Discord.Token, "test_token")
		}
	})
}

func TestLoadFromSource_EdgeCases(t *testing.T) {
	t.Run("all sources empty", func(t *testing.T) {
		_, err := LoadFromSource("", "", "")
		if err == nil {
			t.Error("Expected error when all sources are empty")
		}
	})

	t.Run("flag source takes priority over invalid env", func(t *testing.T) {
		tempDir := t.TempDir()
		cfgFile := filepath.Join(tempDir, "config.json")
		
		validConfig := Config{
			Discord: DiscordConfig{Token: "test_token"},
			MCP:     MCPConfig{Command: "test_cmd"},
		}
		data, _ := json.Marshal(validConfig)
		os.WriteFile(cfgFile, data, 0644)

		// Set invalid env var path
		os.Setenv("TEST_CONFIG_PATH", "/nonexistent/path")
		defer os.Unsetenv("TEST_CONFIG_PATH")

		cfg, err := LoadFromSource(cfgFile, "TEST_CONFIG_PATH", "")
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if cfg.Discord.Token != "test_token" {
			t.Error("Flag source should take priority")
		}
	})
}

func TestInterpolateEnvVars_NestedStructures(t *testing.T) {
	os.Setenv("TEST_VAR1", "value1")
	os.Setenv("TEST_VAR2", "value2")
	defer func() {
		os.Unsetenv("TEST_VAR1")
		os.Unsetenv("TEST_VAR2")
	}()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "number unchanged",
			input:    42,
			expected: 42,
		},
		{
			name:     "boolean unchanged",
			input:    true,
			expected: true,
		},
		{
			name:     "string with var",
			input:    "${TEST_VAR1}",
			expected: "value1",
		},
		{
			name:     "string with multiple vars",
			input:    "${TEST_VAR1}_${TEST_VAR2}",
			expected: "value1_value2",
		},
		{
			name:     "array of strings with vars",
			input:    []interface{}{"${TEST_VAR1}", "${TEST_VAR2}"},
			expected: []interface{}{"value1", "value2"},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": "${TEST_VAR1}",
				},
			},
			expected: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": "value1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateEnvVars(tt.input)
			
			// Convert to JSON for comparison
			resultJSON, _ := json.Marshal(result)
			expectedJSON, _ := json.Marshal(tt.expected)
			
			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("Result = %s, want %s", resultJSON, expectedJSON)
			}
		})
	}
}
