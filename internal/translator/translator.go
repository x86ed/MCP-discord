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
	"strings"

	"github.com/bwmarrin/discordgo"
	"mcpdiscord/internal/mcp"
)

// Translator converts between MCP tools and Discord slash commands.
type Translator interface {
	// ToolToSlashCommand converts an MCP tool to a Discord slash command.
	ToolToSlashCommand(tool mcp.Tool) (*discordgo.ApplicationCommand, error)

	// TranslateArguments converts Discord interaction options to MCP tool arguments.
	TranslateArguments(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error)
}

// DefaultTranslator implements the Translator interface.
type DefaultTranslator struct{}

// New creates a new DefaultTranslator.
func New() *DefaultTranslator {
	return &DefaultTranslator{}
}

// ToolToSlashCommand converts an MCP tool to a Discord slash command.
func (t *DefaultTranslator) ToolToSlashCommand(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
	// Validate Discord limits
	if err := validateDiscordLimits(tool); err != nil {
		return nil, err
	}

	// Sanitize tool name for Discord (1-to-1 mapping: MCP "list" → Discord "/list")
	commandName := SanitizeToolName(tool.Name)

	// Truncate description to Discord's 100 character limit
	description := tool.Description
	if len(description) > 100 {
		description = description[:97] + "..."
	}
	if description == "" {
		description = "No description provided"
	}

	// Convert parameters to Discord options
	options, err := parametersToOptions(tool.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to convert parameters: %w", err)
	}

	return &discordgo.ApplicationCommand{
		Name:        commandName,
		Description: description,
		Options:     options,
	}, nil
}

// SanitizeToolName converts a tool name to a valid Discord command name.
// Implements 1-to-1 mapping: MCP "list" → "/list", "Get Data" → "/get-data"
func SanitizeToolName(name string) string {
	// Lowercase
	name = strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove other special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// parametersToOptions converts MCP input schema properties to Discord options.
func parametersToOptions(schema mcp.InputSchema) ([]*discordgo.ApplicationCommandOption, error) {
	if len(schema.Properties) == 0 {
		return nil, nil
	}

	options := make([]*discordgo.ApplicationCommandOption, 0, len(schema.Properties))

	for name, prop := range schema.Properties {
		// Check if parameter is required
		required := false
		for _, req := range schema.Required {
			if req == name {
				required = true
				break
			}
		}

		// Map MCP type to Discord type
		optionType, hint := mapTypeToDiscord(prop.Type)

		// Build description with hint for arrays/objects
		description := prop.Description
		if hint != "" {
			if description != "" {
				description = description + " " + hint
			} else {
				description = hint
			}
		}
		if description == "" {
			description = "No description"
		}

		// Truncate description to Discord's limit (100 chars)
		if len(description) > 100 {
			description = description[:97] + "..."
		}

		option := &discordgo.ApplicationCommandOption{
			Type:        optionType,
			Name:        name,
			Description: description,
			Required:    required,
		}

		options = append(options, option)
	}

	return options, nil
}

// mapTypeToDiscord maps MCP JSON Schema types to Discord option types.
// Returns the Discord type and a hint to append to the description.
func mapTypeToDiscord(mcpType string) (discordgo.ApplicationCommandOptionType, string) {
	switch mcpType {
	case "string":
		return discordgo.ApplicationCommandOptionString, ""
	case "number", "integer":
		return discordgo.ApplicationCommandOptionNumber, ""
	case "boolean":
		return discordgo.ApplicationCommandOptionBoolean, ""
	case "array":
		return discordgo.ApplicationCommandOptionString, "(Comma-separated list)"
	case "object":
		return discordgo.ApplicationCommandOptionString, "(JSON object)"
	default:
		// Default to string for unknown types
		return discordgo.ApplicationCommandOptionString, ""
	}
}

// validateDiscordLimits checks that the tool fits within Discord's constraints.
func validateDiscordLimits(tool mcp.Tool) error {
	// Discord allows max 25 options per command
	if len(tool.InputSchema.Properties) > 25 {
		return fmt.Errorf("tool '%s' has %d parameters (Discord limit: 25)",
			tool.Name, len(tool.InputSchema.Properties))
	}

	return nil
}

// ValidateCommandCount checks if the number of commands exceeds Discord's limit.
func ValidateCommandCount(count int) error {
	const maxCommands = 100
	if count > maxCommands {
		return fmt.Errorf("cannot register %d commands (Discord limit: %d)", count, maxCommands)
	}
	return nil
}

// TranslateArguments converts Discord interaction options to MCP tool arguments.
func (t *DefaultTranslator) TranslateArguments(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
	args := make(map[string]interface{})

	// Build a map of option values by name
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Process each parameter in the schema
	for paramName, paramSchema := range tool.InputSchema.Properties {
		opt, exists := optionMap[paramName]

		// Handle missing optional parameters (omit from MCP call)
		if !exists {
			isRequired := false
			for _, req := range tool.InputSchema.Required {
				if req == paramName {
					isRequired = true
					break
				}
			}
			if !isRequired {
				continue
			}
			return nil, fmt.Errorf("missing required parameter: %s", paramName)
		}

		// Translate based on parameter type
		value, err := translateValue(paramSchema.Type, opt)
		if err != nil {
			return nil, fmt.Errorf("parameter '%s': %w", paramName, err)
		}

		args[paramName] = value
	}

	return args, nil
}

// translateValue converts a Discord option value to the appropriate Go type based on MCP schema type.
func translateValue(mcpType string, opt *discordgo.ApplicationCommandInteractionDataOption) (interface{}, error) {
	switch mcpType {
	case "string":
		return opt.StringValue(), nil

	case "number", "integer":
		// Discord returns float64 for numbers
		return opt.FloatValue(), nil

	case "boolean":
		return opt.BoolValue(), nil

	case "array":
		// Parse CSV string into array
		csvStr := opt.StringValue()
		return ParseCSV(csvStr), nil

	case "object":
		// Parse JSON string into object
		jsonStr := opt.StringValue()
		return ParseJSON(jsonStr)

	default:
		// Unknown type, return as string
		return opt.StringValue(), nil
	}
}

