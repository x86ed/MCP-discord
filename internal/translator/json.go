package translator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseJSON parses a JSON string into a Go value (map or slice).
// Returns user-friendly error messages for invalid JSON.
func ParseJSON(input string) (interface{}, error) {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Handle empty input
	if input == "" {
		return nil, fmt.Errorf("empty JSON input")
	}

	// Parse JSON
	var result interface{}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w. Example: {\"key\": \"value\"}", err)
	}

	// Validate result is an object or array
	switch result.(type) {
	case map[string]interface{}, []interface{}:
		return result, nil
	default:
		return nil, fmt.Errorf("JSON must be an object {} or array [], got: %T", result)
	}
}

// ValidateJSON checks if a string is valid JSON without parsing it fully.
func ValidateJSON(input string) error {
	_, err := ParseJSON(input)
	return err
}
