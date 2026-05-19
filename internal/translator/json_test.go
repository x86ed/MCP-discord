package translator

import (
	"testing"
)

func TestParseJSONDetailed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "empty string",
			input:       "",
			shouldError: true,
			errorMsg:    "empty JSON input",
		},
		{
			name:        "whitespace only",
			input:       "   ",
			shouldError: true,
			errorMsg:    "empty JSON input",
		},
		{
			name:        "simple object",
			input:       `{"name": "John"}`,
			shouldError: false,
		},
		{
			name:        "nested object",
			input:       `{"user": {"name": "John", "age": 30}}`,
			shouldError: false,
		},
		{
			name:        "array",
			input:       `["a", "b", "c"]`,
			shouldError: false,
		},
		{
			name:        "array of objects",
			input:       `[{"id": 1}, {"id": 2}]`,
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
			name:        "invalid JSON - missing quotes",
			input:       `{name: "John"}`,
			shouldError: true,
		},
		{
			name:        "invalid JSON - trailing comma",
			input:       `{"name": "John",}`,
			shouldError: true,
		},
		{
			name:        "invalid JSON - single quotes",
			input:       `{'name': 'John'}`,
			shouldError: true,
		},
		{
			name:        "plain string",
			input:       `"just a string"`,
			shouldError: true,
		},
		{
			name:        "plain number",
			input:       `42`,
			shouldError: true,
		},
		{
			name:        "plain boolean",
			input:       `true`,
			shouldError: true,
		},
		{
			name:        "plain null",
			input:       `null`,
			shouldError: true,
		},
		{
			name:        "object with various types",
			input:       `{"string": "value", "number": 123, "bool": true, "null": null}`,
			shouldError: false,
		},
		{
			name:        "object with whitespace",
			input:       `  {"name": "John"}  `,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseJSON(tt.input)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("ParseJSON(%q) expected error but got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseJSON(%q) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ParseJSON(%q) returned nil result", tt.input)
				}
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid object",
			input:   `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "invalid object",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
