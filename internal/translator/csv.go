package translator

import (
	"strings"
)

// ParseCSV parses a comma-separated string into an array of strings.
// Whitespace around items is trimmed.
// Empty strings result in empty arrays.
func ParseCSV(input string) []string {
	// Handle empty input
	if strings.TrimSpace(input) == "" {
		return []string{}
	}

	// Split by comma
	items := strings.Split(input, ",")

	// Trim whitespace from each item
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
