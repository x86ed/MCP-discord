package bot

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"mcpdiscord/internal/config"
	"mcpdiscord/internal/mcp"
	"mcpdiscord/internal/translator"

	"github.com/bwmarrin/discordgo"
)

// mockMCPClient implements mcp.Client for testing
type mockMCPClient struct {
	tools       []mcp.Tool
	toolResult  *mcp.ToolResult
	listErr     error
	callErr     error
	connectErr  error
}

func (m *mockMCPClient) Connect(ctx context.Context) error {
	return m.connectErr
}

func (m *mockMCPClient) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	return m.tools, m.listErr
}

func (m *mockMCPClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.ToolResult, error) {
	if m.toolResult == nil {
		return &mcp.ToolResult{
			Content: []mcp.ContentBlock{{Type: "text", Text: "OK"}},
		}, m.callErr
	}
	return m.toolResult, m.callErr
}

func (m *mockMCPClient) Close() error {
	return nil
}

// TestOnReady_HandlerExecution tests the onReady event handler
func TestOnReady_HandlerExecution(t *testing.T) {
	mockSession := NewMockDiscordSession()
	logger := slog.Default()
	mcpClient := &mockMCPClient{
		tools: []mcp.Tool{
			{
				Name:        "test_tool",
				Description: "Test tool",
				InputSchema: mcp.InputSchema{
					Type:       "object",
					Properties: map[string]mcp.PropertySchema{},
				},
			},
		},
	}
	trans := translator.New(logger)
	discovery := mcp.NewDiscoveryService(mcpClient, logger)

	bot := &DiscordBot{
		config: config.DiscordConfig{
			Token:   "test-token",
			GuildID: "test-guild",
		},
		mcpClient:    mcpClient,
		translator:   trans,
		discovery:    discovery,
		session:      mockSession,
		logger:       logger,
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	ready := &discordgo.Ready{
		User: &discordgo.User{
			ID:       "bot-123",
			Username: "TestBot",
		},
		Guilds: []*discordgo.Guild{
			{ID: "guild-1"},
			{ID: "guild-2"},
		},
	}

	// Call onReady handler
	bot.onReady(nil, ready)

	// Verify commands were attempted to be registered
	if mockSession.CreateCommandCalls == 0 {
		t.Error("Expected commands to be registered via CreateCommand")
	}
}

// TestOnInteractionCreate_CommandExecution tests interaction handling
func TestOnInteractionCreate_CommandExecution(t *testing.T) {
	mockSession := NewMockDiscordSession()
	logger := slog.Default()
	mcpClient := &mockMCPClient{}
	trans := translator.New(logger)

	tool := mcp.Tool{
		Name:        "test_cmd",
		Description: "Test",
		InputSchema: mcp.InputSchema{
			Type:       "object",
			Properties: map[string]mcp.PropertySchema{},
		},
	}

	bot := &DiscordBot{
		session:    mockSession,
		mcpClient:  mcpClient,
		logger:     logger,
		translator: trans,
		commandTools: map[string]mcp.Tool{
			"test_cmd": tool,
		},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name:    "test_cmd",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{},
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID:       "user-123",
					Username: "testuser",
				},
			},
			GuildID: "guild-123",
		},
	}

	// Call interaction handler
	bot.onInteractionCreate(nil, interaction)

	// Verify interaction was acknowledged
	if mockSession.RespondCalls == 0 {
		t.Error("Expected interaction to be acknowledged")
	}

	if mockSession.LastInteractionRespond.Type != discordgo.InteractionResponseDeferredChannelMessageWithSource {
		t.Errorf("Expected deferred response, got %v", mockSession.LastInteractionRespond.Type)
	}
}

// TestOnInteractionCreate_NonCommandIgnored tests that non-commands are ignored
func TestOnInteractionCreate_NonCommandIgnored(t *testing.T) {
	mockSession := NewMockDiscordSession()
	logger := slog.Default()

	bot := &DiscordBot{
		session: mockSession,
		logger:  logger,
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent, // Not a command
		},
	}

	bot.onInteractionCreate(nil, interaction)

	// Verify no response was sent
	if mockSession.RespondCalls != 0 {
		t.Error("Expected non-command to be ignored")
	}
}

// TestOnInteractionCreate_RespondError tests error handling when respond fails
func TestOnInteractionCreate_RespondError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.RespondError = fmt.Errorf("respond failed")
	logger := slog.Default()

	bot := &DiscordBot{
		session: mockSession,
		logger:  logger,
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
		},
	}

	bot.onInteractionCreate(nil, interaction)

	// Verify respond was attempted
	if mockSession.RespondCalls != 1 {
		t.Errorf("Expected 1 respond call, got %d", mockSession.RespondCalls)
	}
	// The bot should log the error and return without calling handleCommand
}

// TestLooksLikeJSON_EdgeCases tests JSON detection
func TestLooksLikeJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{" ", false},
		{"\t\n", false},
		{"{}", true},
		{"[]", true},
		{" { } ", true},
		{" [ ] ", true},
		{"not json", false},
		{"{incomplete", true},  // Looks like JSON (starts with {)
		{"[incomplete", true},  // Looks like JSON (starts with [)
	}

	for _, tt := range tests {
		result := looksLikeJSON(tt.input)
		if result != tt.expected {
			t.Errorf("looksLikeJSON(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

// TestImageURLFromResourceLink_EdgeCases tests image URL extraction
func TestImageURLFromResourceLink_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		link     *mcp.ResourceLink
		expected string
	}{
		{
			name:     "nil",
			link:     nil,
			expected: "",
		},
		{
			name:     "empty",
			link:     &mcp.ResourceLink{},
			expected: "",
		},
		{
			name: "url priority",
			link: &mcp.ResourceLink{
				URL: "https://example.com/a.png",
				URI: "https://example.com/b.png",
			},
			expected: "https://example.com/a.png",
		},
		{
			name: "uri fallback",
			link: &mcp.ResourceLink{
				URI: "https://example.com/c.jpg",
			},
			expected: "https://example.com/c.jpg",
		},
		{
			name: "href fallback",
			link: &mcp.ResourceLink{
				Href: "https://example.com/d.gif",
			},
			expected: "https://example.com/d.gif",
		},
		{
			name: "file scheme rejected",
			link: &mcp.ResourceLink{
				URL: "file:///local.png",
			},
			expected: "",
		},
		{
			name: "non-image extension",
			link: &mcp.ResourceLink{
				URL: "https://example.com/doc.pdf",
			},
			expected: "",
		},
		{
			name: "image mime type",
			link: &mcp.ResourceLink{
				URL:      "https://example.com/img",
				MimeType: "image/png",
			},
			expected: "https://example.com/img",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := imageURLFromResourceLink(tt.link)
			if result != tt.expected {
				t.Errorf("got %q, expected %q", result, tt.expected)
			}
		})
	}
}
