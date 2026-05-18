// Package translator provides translation between MCP tool definitions and
// Discord command definitions, as well as between Discord interactions and
// MCP tool calls.
//
// This package bridges the gap between the Model Context Protocol's JSON Schema
// based tool definitions and Discord's slash command system. It handles:
//
//   - Converting MCP tool schemas to Discord command definitions
//   - Mapping JSON Schema types to Discord option types
//   - Translating Discord interaction options to MCP tool arguments
//   - Formatting MCP tool results as Discord messages
//
// The translator implements intelligent schema mapping to handle various
// JSON Schema features including nested objects, arrays, enums, and complex
// validation rules. When direct translation isn't possible (e.g., deeply
// nested objects), it falls back to JSON string representations.
//
// Key responsibilities:
//   - MCP Tool → Discord Command translation
//   - Discord Interaction → MCP Tool Call translation
//   - MCP Tool Result → Discord Response formatting
//   - Schema type mapping and validation
//   - Error message formatting
//
// Example usage:
//
//	translator := translator.New()
//
//	// Translate MCP tool to Discord command
//	mcpTool := mcp.Tool{
//	    Name: "get_weather",
//	    Description: "Get current weather",
//	    InputSchema: mcp.Schema{...},
//	}
//	command := translator.ToolToCommand(mcpTool)
//
//	// Translate Discord interaction to MCP tool call
//	interaction := bot.Interaction{
//	    CommandName: "get_weather",
//	    Options: map[string]interface{}{"location": "Seattle"},
//	}
//	args := translator.InteractionToArguments(interaction)
//
//	// Format MCP result as Discord response
//	result := mcp.ToolResult{Content: [...]}
//	response := translator.ResultToResponse(result)
package translator

import (
	"fmt"

	"mcpdiscord/internal/bot"
	"mcpdiscord/internal/mcp"
)

// Translator handles bidirectional translation between MCP and Discord formats.
type Translator interface {
	// ToolToCommand converts an MCP tool definition to a Discord slash command.
	// It translates the JSON Schema to Discord command options.
	ToolToCommand(tool mcp.Tool) (*bot.Command, error)

	// InteractionToArguments converts Discord interaction options to MCP tool arguments.
	// It transforms Discord's option format to the JSON structure expected by the MCP tool.
	InteractionToArguments(interaction *bot.Interaction) (map[string]interface{}, error)

	// ResultToResponse converts an MCP tool result to a Discord response.
	// It formats the content blocks as Discord messages with appropriate formatting.
	ResultToResponse(result *mcp.ToolResult) (*bot.Response, error)

	// ErrorToResponse converts an error to a user-friendly Discord response.
	// It formats errors in a way that's helpful for Discord users.
	ErrorToResponse(err error) *bot.Response
}

// translator is the concrete implementation of the Translator interface.
// TODO: Implement in future change
type translator struct {
	maxChoices int // Maximum number of enum values to translate as Discord choices (25 limit)
}

// New creates a new translator instance.
// TODO: Implement in future change
func New() Translator {
	return &translator{
		maxChoices: 25, // Discord's maximum number of choices
	}
}

// Schema Translation Rules
//
// The following table describes how JSON Schema types are mapped to Discord option types:
//
// | JSON Schema Type | Discord Option Type | Notes                                    |
// |------------------|---------------------|------------------------------------------|
// | string           | STRING              | Direct mapping                           |
// | integer          | INTEGER             | Direct mapping                           |
// | number           | NUMBER              | Supports decimals                        |
// | boolean          | BOOLEAN             | True/false values                        |
// | enum (≤25 items) | STRING with choices | Discord limit: 25 choices                |
// | enum (>25 items) | STRING              | No choices, user types value             |
// | object (simple)  | Multiple options    | Flatten one level                        |
// | object (nested)  | STRING              | JSON string fallback                     |
// | array            | STRING              | JSON string fallback (e.g., "[1,2,3]")   |
//
// Limitations:
//   - Discord commands support max 25 options
//   - Discord choices support max 25 values
//   - Option names must be lowercase, alphanumeric with underscores/hyphens
//   - Descriptions limited to 100 characters
//
// When a schema cannot be directly translated, the translator uses JSON string
// fallback, allowing users to provide JSON-formatted values that are parsed
// before calling the MCP tool.

// ToolToCommand converts an MCP tool to a Discord command.
// TODO: Implement in future change
func (t *translator) ToolToCommand(tool mcp.Tool) (*bot.Command, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// InteractionToArguments converts a Discord interaction to MCP tool arguments.
// TODO: Implement in future change
func (t *translator) InteractionToArguments(interaction *bot.Interaction) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// ResultToResponse converts an MCP tool result to a Discord response.
// TODO: Implement in future change
func (t *translator) ResultToResponse(result *mcp.ToolResult) (*bot.Response, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// ErrorToResponse converts an error to a Discord response.
// TODO: Implement in future change
func (t *translator) ErrorToResponse(err error) *bot.Response {
	return &bot.Response{
		Content:   fmt.Sprintf("Error: %v", err),
		Ephemeral: true,
	}
}

// Formatting Utilities
//
// These helper functions format content for Discord's markdown-style formatting:
//
//   - Bold: **text**
//   - Italic: *text* or _text_
//   - Code: `code`
//   - Code Block: ```language\ncode\n```
//   - Quote: > quote
//   - Spoiler: ||spoiler||

// formatText applies basic formatting to text content.
// TODO: Implement in future change
func formatText(text string) string {
	return text
}

// truncateText truncates text to fit Discord's message limits.
// Discord limits: 2000 chars for message content, 4096 for embed description
// TODO: Implement in future change
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

// sanitizeCommandName converts a tool name to a valid Discord command name.
// Discord requirements: lowercase, alphanumeric with hyphens/underscores, 1-32 chars
// TODO: Implement in future change
func sanitizeCommandName(name string) string {
	return name
}

// sanitizeOptionName converts a parameter name to a valid Discord option name.
// Same requirements as command names
// TODO: Implement in future change
func sanitizeOptionName(name string) string {
	return name
}

// truncateDescription truncates a description to fit Discord's limits.
// Discord limits: 100 chars for command/option descriptions
// TODO: Implement in future change
func truncateDescription(description string) string {
	if len(description) <= 100 {
		return description
	}
	return description[:97] + "..."
}
