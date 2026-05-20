package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// DiscoveryService handles tool discovery from MCP servers.
type DiscoveryService struct {
	client Client
	logger *slog.Logger
}

// NewDiscoveryService creates a new tool discovery service.
func NewDiscoveryService(client Client, logger *slog.Logger) *DiscoveryService {
	if logger == nil {
		logger = slog.Default()
	}
	return &DiscoveryService{
		client: client,
		logger: logger,
	}
}

// DiscoverTools performs tool discovery with validation and collision detection.
// Returns discovered tools or an error if discovery fails.
func (d *DiscoveryService) DiscoverTools(ctx context.Context) ([]Tool, error) {
	d.logger.Info("starting tool discovery")

	// Set 120 second timeout for tool listing
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// List tools from MCP server
	tools, err := d.client.ListTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	// Validate empty tool list
	if len(tools) == 0 {
		return nil, fmt.Errorf("MCP server returned empty tool list")
	}

	// Validate and process each tool
	validatedTools := make([]Tool, 0, len(tools))
	for i, tool := range tools {
		if err := ValidateTool(&tool); err != nil {
			d.logger.Warn("skipping invalid tool",
				"index", i,
				"name", tool.Name,
				"error", err)
			continue
		}

		// Truncate description if needed
		tool.Description = TruncateDescription(tool.Description)

		validatedTools = append(validatedTools, tool)
	}

	// Check for name collisions
	if err := CheckNameCollisions(validatedTools); err != nil {
		return nil, err
	}

	d.logger.Info("tool discovery complete",
		"total", len(tools),
		"valid", len(validatedTools))

	return validatedTools, nil
}

// DiscoverWithReconnect attempts tool discovery with automatic reconnection on failure.
func (d *DiscoveryService) DiscoverWithReconnect(ctx context.Context) ([]Tool, error) {
	// Try discovery first
	tools, err := d.DiscoverTools(ctx)
	if err == nil {
		return tools, nil
	}

	d.logger.Warn("initial tool discovery failed, attempting reconnection", "error", err)

	// Close existing connection
	if closeErr := d.client.Close(); closeErr != nil {
		d.logger.Error("failed to close client during reconnection", "error", closeErr)
	}

	// Reconnect
	if err := d.client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to reconnect to MCP server: %w", err)
	}

	// Retry discovery
	tools, err = d.DiscoverTools(ctx)
	if err != nil {
		return nil, fmt.Errorf("tool discovery failed after reconnection: %w", err)
	}

	d.logger.Info("tool discovery succeeded after reconnection")
	return tools, nil
}

// ValidateTool checks that a tool has required metadata fields.
func ValidateTool(tool *Tool) error {
	// Validate name presence
	if tool.Name == "" {
		return fmt.Errorf("tool missing required 'name' field")
	}

	// Validate description presence (can be empty but field must exist)
	// Note: Description is a string, not a pointer, so this always passes
	// We're keeping this check for API consistency and future changes

	// Validate input schema
	if tool.InputSchema.Type == "" {
		return fmt.Errorf("tool '%s' missing input schema type", tool.Name)
	}

	// Validate schema type is "object" (MCP tools expect object parameters)
	if tool.InputSchema.Type != "object" {
		return fmt.Errorf("tool '%s' has invalid schema type '%s' (expected 'object')",
			tool.Name, tool.InputSchema.Type)
	}

	// Validate properties map exists (can be empty)
	if tool.InputSchema.Properties == nil {
		tool.InputSchema.Properties = make(map[string]PropertySchema)
	}

	return nil
}

// TruncateDescription truncates a description to 100 characters with "..." suffix.
func TruncateDescription(desc string) string {
	const maxLen = 100
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:97] + "..."
}

// CheckNameCollisions detects if multiple tools would map to the same sanitized command name.
func CheckNameCollisions(tools []Tool) error {
	seen := make(map[string][]string)

	for _, tool := range tools {
		// Sanitize the tool name (this will be implemented in translator package)
		// For now, we use a simple sanitization here
		sanitized := sanitizeToolName(tool.Name)

		seen[sanitized] = append(seen[sanitized], tool.Name)
	}

	// Report collisions
	var collisions []string
	for sanitized, originals := range seen {
		if len(originals) > 1 {
			collisions = append(collisions,
				fmt.Sprintf("'%s' -> %v", sanitized, originals))
		}
	}

	if len(collisions) > 0 {
		return fmt.Errorf("tool name collisions detected after sanitization: %s",
			strings.Join(collisions, "; "))
	}

	return nil
}

// sanitizeToolName converts a tool name to a valid Discord command name.
// This is a simplified version - the full implementation is in the translator package.
func sanitizeToolName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	// Remove other special characters
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ParseInputSchema extracts parameter information from a JSON schema.
func ParseInputSchema(schema InputSchema) ([]Parameter, error) {
	params := make([]Parameter, 0, len(schema.Properties))

	for name, prop := range schema.Properties {
		required := false
		for _, req := range schema.Required {
			if req == name {
				required = true
				break
			}
		}

		param := Parameter{
			Name:        name,
			Type:        prop.Type,
			Description: prop.Description,
			Required:    required,
			Schema:      prop,
		}

		params = append(params, param)
	}

	return params, nil
}

// Parameter represents a parsed tool parameter.
type Parameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Schema      PropertySchema
}
