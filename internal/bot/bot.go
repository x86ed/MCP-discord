// Package bot implements the Discord bot functionality for the MCP-Discord bridge.
//
// This package provides the core Discord integration, handling bot initialization,
// slash command registration, and interaction processing. It acts as the interface
// between Discord users and the MCP server.
//
// The bot automatically discovers available tools from the connected MCP server
// and registers them as Discord slash commands. When users invoke these commands,
// the bot coordinates with the translator and MCP client packages to execute
// the requested tools and return formatted responses.
//
// Key responsibilities:
//   - Discord Gateway connection management
//   - Slash command registration (guild-specific or global)
//   - Interaction handling (commands, options, responses)
//   - Error reporting to users
//   - Graceful shutdown coordination
//
// Example usage:
//
//	cfg := config.Load("config.json")
//	client := mcp.NewClient(cfg.MCP)
//	translator := translator.New()
//	bot := bot.New(cfg.Discord, client, translator)
//
//	if err := bot.Start(); err != nil {
//	    log.Fatal(err)
//	}
//	defer bot.Stop()
package bot

import (
	"context"
	"fmt"

	"mcpdiscord/internal/config"
)

// Bot represents a Discord bot instance that bridges to an MCP server.
// It manages the Discord connection, command registration, and interaction handling.
type Bot interface {
	// Start initializes the Discord connection, registers commands, and begins
	// listening for interactions. Returns an error if connection fails.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the bot, disconnecting from Discord and
	// cleaning up resources.
	Stop() error

	// RegisterCommands registers all available MCP tools as Discord slash commands.
	// If guildID is empty, registers commands globally (takes up to 1 hour).
	// If guildID is set, registers to that specific guild (instant).
	RegisterCommands(guildID string) error

	// HandleInteraction processes a Discord interaction event and returns a response.
	// This is called automatically when users invoke slash commands.
	HandleInteraction(interaction *Interaction) (*Response, error)
}

// Interaction represents a Discord slash command interaction from a user.
type Interaction struct {
	// ID is the unique interaction ID from Discord
	ID string

	// CommandName is the name of the invoked slash command
	CommandName string

	// Options contains the command parameters provided by the user
	Options map[string]interface{}

	// UserID is the Discord user ID who invoked the command
	UserID string

	// GuildID is the Discord server ID where the command was invoked
	GuildID string

	// ChannelID is the Discord channel ID where the command was invoked
	ChannelID string
}

// Response represents a response to send back to Discord after handling an interaction.
type Response struct {
	// Content is the text content of the response
	Content string

	// Ephemeral indicates whether the response should only be visible to the user
	Ephemeral bool

	// Embed is an optional rich embed for formatted responses
	Embed *Embed
}

// Embed represents a Discord rich embed for formatted responses.
type Embed struct {
	// Title is the embed title
	Title string

	// Description is the main embed content
	Description string

	// Color is the embed color (integer RGB)
	Color int

	// Fields are optional additional fields
	Fields []EmbedField
}

// EmbedField represents a field in a Discord embed.
type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

// Command represents a Discord slash command definition.
type Command struct {
	// Name is the command name (must be lowercase, no spaces)
	Name string

	// Description is the command description shown in Discord
	Description string

	// Options are the command parameters
	Options []CommandOption
}

// CommandOption represents a parameter for a Discord slash command.
type CommandOption struct {
	// Name is the option name (must be lowercase, no spaces)
	Name string

	// Description is the option description shown in Discord
	Description string

	// Type is the option type (string, integer, boolean, etc.)
	Type OptionType

	// Required indicates whether this option must be provided
	Required bool

	// Choices are predefined values for this option (optional)
	Choices []Choice
}

// OptionType represents a Discord command option type.
type OptionType int

const (
	OptionTypeString OptionType = iota
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeNumber
)

// Choice represents a predefined choice for a command option.
type Choice struct {
	Name  string
	Value interface{}
}

// bot is the concrete implementation of the Bot interface.
// TODO: Implement in future change
type bot struct {
	config     config.DiscordConfig
	mcpClient  interface{} // TODO: Replace with actual MCP client interface
	translator interface{} // TODO: Replace with actual translator interface
}

// New creates a new Discord bot instance.
// TODO: Implement in future change
func New(cfg config.DiscordConfig, mcpClient, translator interface{}) (Bot, error) {
	return nil, fmt.Errorf("not yet implemented")
}
