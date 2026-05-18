package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a configuration validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid configuration: %s: %s", e.Field, e.Message)
}

// Validate checks if the configuration is valid and returns descriptive errors.
func (c *Config) Validate() error {
	var errors []string

	// Validate Discord configuration
	if c.Discord.Token == "" {
		errors = append(errors, "discord.token is required")
	} else if strings.Contains(c.Discord.Token, "${") {
		errors = append(errors, "discord.token contains unresolved environment variable")
	}

	// Validate MCP configuration
	if c.MCP.Command == "" {
		errors = append(errors, "mcp.command is required")
	}

	// Validate transport type if specified
	if c.MCP.Transport != "" {
		validTransports := map[string]bool{
			"stdio":     true,
			"sse":       true,
			"websocket": true,
		}
		if !validTransports[c.MCP.Transport] {
			errors = append(errors, fmt.Sprintf("mcp.transport must be one of: stdio, sse, websocket (got %q)", c.MCP.Transport))
		}
	}

	// Check for unresolved environment variables in MCP env vars
	for key, value := range c.MCP.Env {
		if strings.Contains(value, "${") {
			errors = append(errors, fmt.Sprintf("mcp.env.%s contains unresolved environment variable", key))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}
