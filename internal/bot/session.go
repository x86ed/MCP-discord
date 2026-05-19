package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"mcpdiscord/internal/config"
	"mcpdiscord/internal/mcp"
	"mcpdiscord/internal/translator"
)

// DiscordBot implements the Bot interface using discordgo.
type DiscordBot struct {
	config     config.DiscordConfig
	mcpClient  mcp.Client
	translator translator.Translator
	discovery  *mcp.DiscoveryService
	session    *discordgo.Session
	logger     *slog.Logger

	// Registered commands tracking
	commands     []*discordgo.ApplicationCommand
	commandTools map[string]mcp.Tool // Maps command name to MCP tool
}

// NewDiscordBot creates a new Discord bot instance.
func NewDiscordBot(
	cfg config.DiscordConfig,
	mcpClient mcp.Client,
	trans translator.Translator,
	discovery *mcp.DiscoveryService,
	logger *slog.Logger,
) (*DiscordBot, error) {
	if logger == nil {
		logger = slog.Default()
	}

	// Create Discord session with bot token
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &DiscordBot{
		config:       cfg,
		mcpClient:    mcpClient,
		translator:   trans,
		discovery:    discovery,
		session:      session,
		logger:       logger,
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	// Register event handlers
	session.AddHandler(bot.onReady)
	session.AddHandler(bot.onInteractionCreate)

	return bot, nil
}

// Start connects to Discord and begins listening for events.
func (b *DiscordBot) Start(ctx context.Context) error {
	b.logger.Info("starting Discord bot")

	// Open websocket connection
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	b.logger.Info("Discord bot connected", "user", b.session.State.User.String())
	return nil
}

// Stop gracefully shuts down the Discord bot.
func (b *DiscordBot) Stop() error {
	b.logger.Info("stopping Discord bot")

	if b.session != nil {
		// Deregister all commands if guild is specified
		if b.config.GuildID != "" {
			b.logger.Info("deregistering commands", "guild", b.config.GuildID)
			if err := b.deregisterCommands(b.config.GuildID); err != nil {
				b.logger.Warn("failed to deregister commands", "error", err)
			}
		}

		if err := b.session.Close(); err != nil {
			return fmt.Errorf("failed to close Discord session: %w", err)
		}
	}

	return nil
}

// RegisterCommands discovers MCP tools and registers them as Discord slash commands.
func (b *DiscordBot) RegisterCommands(guildID string) error {
	b.logger.Info("registering commands", "guild", guildID)

	// Discover tools from MCP server
	tools, err := b.discovery.DiscoverTools(context.Background())
	if err != nil {
		return fmt.Errorf("failed to discover tools: %w", err)
	}

	// Validate command count against Discord limit
	if err := translator.ValidateCommandCount(len(tools)); err != nil {
		return err
	}

	// Convert tools to Discord commands
	commands := make([]*discordgo.ApplicationCommand, 0, len(tools))
	b.commandTools = make(map[string]mcp.Tool)

	for _, tool := range tools {
		cmd, err := b.translator.ToolToSlashCommand(tool)
		if err != nil {
			b.logger.Warn("skipping tool", "name", tool.Name, "error", err)
			continue
		}

		commands = append(commands, cmd)
		b.commandTools[cmd.Name] = tool
		b.logger.Debug("mapped tool to command", "tool", tool.Name, "command", cmd.Name)
	}

	// Deregister old commands first
	if err := b.deregisterCommands(guildID); err != nil {
		b.logger.Warn("failed to deregister old commands", "error", err)
	}

	// Register new commands
	b.commands = make([]*discordgo.ApplicationCommand, 0, len(commands))
	for _, cmd := range commands {
		registered, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, cmd)
		if err != nil {
			b.logger.Error("failed to register command", "name", cmd.Name, "error", err)
			continue
		}
		b.commands = append(b.commands, registered)
	}

	b.logger.Info("command registration complete",
		"registered", len(b.commands),
		"total", len(tools))

	return nil
}

// deregisterCommands removes all registered commands from Discord.
func (b *DiscordBot) deregisterCommands(guildID string) error {
	// Get existing commands
	existing, err := b.session.ApplicationCommands(b.session.State.User.ID, guildID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing commands: %w", err)
	}

	// Delete each command
	for _, cmd := range existing {
		if err := b.session.ApplicationCommandDelete(b.session.State.User.ID, guildID, cmd.ID); err != nil {
			b.logger.Warn("failed to delete command", "id", cmd.ID, "name", cmd.Name, "error", err)
		}
	}

	return nil
}

// onReady is called when the Discord bot connects successfully.
func (b *DiscordBot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	b.logger.Info("Discord bot ready",
		"guilds", len(event.Guilds),
		"user", event.User.String())

	// Register commands for configured guild or globally
	guildID := b.config.GuildID
	if err := b.RegisterCommands(guildID); err != nil {
		b.logger.Error("failed to register commands on ready", "error", err)
	}
}

// onInteractionCreate handles Discord interaction events (slash command invocations).
func (b *DiscordBot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Immediately acknowledge to prevent timeout (3 second limit)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("failed to acknowledge interaction", "error", err)
		return
	}

	// Process the command
	go b.handleCommand(s, i)
}

// handleCommand processes a slash command interaction.
func (b *DiscordBot) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	commandName := data.Name

	b.logger.Info("handling command",
		"command", commandName,
		"user", i.Member.User.Username,
		"guild", i.GuildID)

	// Look up the MCP tool for this command
	tool, exists := b.commandTools[commandName]
	if !exists {
		b.sendError(s, i, fmt.Errorf("command '%s' not found", commandName), true)
		return
	}

	// Translate Discord options to MCP arguments
	args, err := b.translator.TranslateArguments(tool, data.Options)
	if err != nil {
		b.sendError(s, i, fmt.Errorf("invalid arguments: %w", err), false)
		return
	}

	b.logger.Debug("translated arguments", "command", commandName, "args", args)

	// Execute MCP tool
	ctx := context.Background()
	result, err := b.mcpClient.CallTool(ctx, tool.Name, args)
	if err != nil {
		b.sendError(s, i, fmt.Errorf("tool execution failed: %w", err), true)
		return
	}

	// Format and send response
	response := b.formatResult(tool.Name, result)
	b.sendResponse(s, i, response)
}

// formatResult converts an MCP tool result into a Discord response.
func (b *DiscordBot) formatResult(toolName string, result *mcp.ToolResult) *Response {
	// Determine color based on error status
	color := 0x00FF00 // Green for success
	title := fmt.Sprintf("✓ %s", toolName)
	if result.IsError {
		color = 0xFF0000 // Red for error
		title = fmt.Sprintf("✗ %s", toolName)
	}

	// Combine content blocks
	var contentBuilder strings.Builder
	for _, block := range result.Content {
		if block.Type == "text" {
			contentBuilder.WriteString(block.Text)
			contentBuilder.WriteString("\n")
		}
	}

	content := strings.TrimSpace(contentBuilder.String())

	// Truncate if exceeds Discord limit (4096 chars for embed description)
	const maxLength = 4096
	if len(content) > maxLength {
		content = content[:maxLength-20] + "\n\n...(truncated)"
	}

	// Format structured data with code blocks
	if strings.Contains(content, "{") || strings.Contains(content, "[") {
		content = "```json\n" + content + "\n```"
	}

	return &Response{
		Embed: &Embed{
			Title:       title,
			Description: content,
			Color:       color,
		},
		Ephemeral: false,
	}
}

// sendResponse sends a response to a Discord interaction.
func (b *DiscordBot) sendResponse(s *discordgo.Session, i *discordgo.InteractionCreate, resp *Response) {
	var embeds []*discordgo.MessageEmbed
	if resp.Embed != nil {
		embeds = []*discordgo.MessageEmbed{
			{
				Title:       resp.Embed.Title,
				Description: resp.Embed.Description,
				Color:       resp.Embed.Color,
			},
		}
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &resp.Content,
		Embeds:  &embeds,
	})
	if err != nil {
		b.logger.Error("failed to send response", "error", err)
	}
}

// sendError sends an error message to a Discord interaction.
func (b *DiscordBot) sendError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, isInternal bool) {
	b.logger.Error("command error",
		"command", i.ApplicationCommandData().Name,
		"user", i.Member.User.Username,
		"error", err)

	message := err.Error()
	if isInternal {
		message = "An internal error occurred. Please try again later."
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Error",
		Description: message,
		Color:       0xFF0000, // Red
	}

	_, sendErr := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if sendErr != nil {
		b.logger.Error("failed to send error response", "error", sendErr)
	}
}

// HandleInteraction processes an interaction (stub for interface compliance).
func (b *DiscordBot) HandleInteraction(interaction *Interaction) (*Response, error) {
	// This is handled by the event handler in practice
	return nil, fmt.Errorf("not implemented: use event handlers")
}
