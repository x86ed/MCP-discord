package translator

import (
	"testing"

	"mcpdiscord/internal/mcp"
)

func TestSanitizeToolName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "list",
			expected: "list",
		},
		{
			name:     "uppercase to lowercase",
			input:    "GetWeather",
			expected: "getweather",
		},
		{
			name:     "spaces to hyphens",
			input:    "Get Data",
			expected: "get-data",
		},
		{
			name:     "underscores to hyphens",
			input:    "get_weather_data",
			expected: "get-weather-data",
		},
		{
			name:     "special characters removed",
			input:    "get@weather!data#",
			expected: "getweatherdata",
		},
		{
			name:     "mixed case and special chars",
			input:    "Get Weather_Data!",
			expected: "get-weather-data",
		},
		{
			name:     "numbers preserved",
			input:    "tool123",
			expected: "tool123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeToolName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeToolName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple list",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with whitespace",
			input:    "a, b, c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "single item",
			input:    "item",
			expected: []string{"item"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "trailing comma",
			input:    "a,b,c,",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "leading comma",
			input:    ",a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "multiple spaces",
			input:    "a  ,  b  ,  c",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCSV(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCSV(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ParseCSV(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "valid object",
			input:       `{"key": "value"}`,
			shouldError: false,
		},
		{
			name:        "valid array",
			input:       `["a", "b", "c"]`,
			shouldError: false,
		},
		{
			name:        "nested object",
			input:       `{"key": {"nested": "value"}}`,
			shouldError: false,
		},
		{
			name:        "empty object",
			input:       `{}`,
			shouldError: false,
		},
		{
			name:        "empty array",
			input:       `[]`,
			shouldError: false,
		},
		{
			name:        "invalid JSON",
			input:       `{invalid}`,
			shouldError: true,
		},
		{
			name:        "empty string",
			input:       ``,
			shouldError: true,
		},
		{
			name:        "plain string",
			input:       `"string"`,
			shouldError: true,
		},
		{
			name:        "plain number",
			input:       `123`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ParseJSON(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseJSON(%q) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ParseJSON(%q) result is nil", tt.input)
				}
			}
		})
	}
}

func TestValidateCommandCount(t *testing.T) {
	tests := []struct {
		name        string
		count       int
		shouldError bool
	}{
		{
			name:        "under limit",
			count:       50,
			shouldError: false,
		},
		{
			name:        "at limit",
			count:       100,
			shouldError: false,
		},
		{
			name:        "over limit",
			count:       101,
			shouldError: true,
		},
		{
			name:        "zero commands",
			count:       0,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommandCount(tt.count)
			if tt.shouldError && err == nil {
				t.Errorf("ValidateCommandCount(%d) expected error", tt.count)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("ValidateCommandCount(%d) unexpected error: %v", tt.count, err)
			}
		})
	}
}

func TestTranslateArguments(t *testing.T) {
	trans := New()

	// Test with a simple tool schema
	tool := mcp.Tool{
		Name: "test_tool",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"name": {
					Type:        "string",
					Description: "A name",
				},
				"count": {
					Type:        "number",
					Description: "A count",
				},
			},
			Required: []string{"name"},
		},
	}

	// Mock Discord options - we can't easily test with real discordgo types
	// This is a placeholder test that would need actual discordgo mocking
	t.Run("missing required parameter", func(t *testing.T) {
		_, err := trans.TranslateArguments(tool, nil)
		if err == nil {
			t.Error("TranslateArguments with missing required parameter should error")
		}
	})
}
