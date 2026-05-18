// Package config provides configuration management for the MCP-Discord bot.
// It handles loading, parsing, and validating configuration from JSON files
// with support for environment variable interpolation.
package config

// Config represents the complete configuration for the MCP-Discord bot.
type Config struct {
	// Discord contains Discord bot configuration
	Discord DiscordConfig `json:"discord"`

	// MCP contains MCP server connection configuration
	MCP MCPConfig `json:"mcp"`
}

// DiscordConfig holds Discord-specific configuration.
type DiscordConfig struct {
	// Token is the Discord bot token (can use ${ENV_VAR} syntax)
	Token string `json:"token"`

	// GuildID is optional - when set, commands register only to this guild (faster for dev)
	GuildID string `json:"guildId,omitempty"`
}

// MCPConfig holds MCP server connection configuration.
type MCPConfig struct {
	// Command is the executable path or command to launch the MCP server
	Command string `json:"command"`

	// Args are the command-line arguments for the MCP server
	Args []string `json:"args,omitempty"`

	// Env contains environment variables to set for the MCP server
	Env map[string]string `json:"env,omitempty"`

	// Transport specifies the connection type: "stdio", "sse", or "websocket"
	// Default is "stdio"
	Transport string `json:"transport,omitempty"`
}
