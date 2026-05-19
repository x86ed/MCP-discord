package translator

import (
	"testing"
)

func TestParseCSVDetailed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
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
			name:     "single item",
			input:    "apple",
			expected: []string{"apple"},
		},
		{
			name:     "single item with spaces",
			input:    "  apple  ",
			expected: []string{"apple"},
		},
		{
			name:     "two items",
			input:    "apple,banana",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "two items with spaces",
			input:    "apple , banana",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "three items",
			input:    "apple,banana,cherry",
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "items with extra spaces",
			input:    "  apple  ,  banana  ,  cherry  ",
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "trailing comma",
			input:    "apple,banana,",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "leading comma",
			input:    ",apple,banana",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "multiple consecutive commas",
			input:    "apple,,banana",
			expected: []string{"apple", "banana"},
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: []string{},
		},
		{
			name:     "numbers",
			input:    "1,2,3,4,5",
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			name:     "mixed content",
			input:    "user123, john@example.com, active",
			expected: []string{"user123", "john@example.com", "active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCSV(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("ParseCSV(%q) returned %d items, expected %d", tt.input, len(result), len(tt.expected))
				t.Errorf("Got: %v", result)
				t.Errorf("Expected: %v", tt.expected)
				return
			}
			
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ParseCSV(%q)[%d] = %q, expected %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}
