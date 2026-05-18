package testutil

import (
	"testing"
)

// MockDiscordInteraction represents a mock Discord interaction for testing.
type MockDiscordInteraction struct {
	t              *testing.T
	CommandName    string
	Options        map[string]interface{}
	UserID         string
	GuildID        string
	ChannelID      string
	ResponseSent   bool
	ResponseData   *MockInteractionResponse
}

// MockInteractionResponse represents a mock interaction response.
type MockInteractionResponse struct {
	Type      int
	Content   string
	Ephemeral bool
	Embeds    []MockEmbed
}

// MockEmbed represents a mock Discord embed.
type MockEmbed struct {
	Title       string
	Description string
	Fields      []MockEmbedField
}

// MockEmbedField represents a mock embed field.
type MockEmbedField struct {
	Name   string
	Value  string
	Inline bool
}

// MockApplicationCommand represents a Discord application command definition.
type MockApplicationCommand struct {
	Name        string
	Description string
	Options     []MockCommandOption
}

// MockCommandOption represents a command option.
type MockCommandOption struct {
	Name        string
	Description string
	Type        int // Discord option type
	Required    bool
	Choices     []MockChoice
}

// MockChoice represents a command option choice.
type MockChoice struct {
	Name  string
	Value interface{}
}

// NewMockInteraction creates a new mock Discord interaction.
func NewMockInteraction(t *testing.T, commandName string, options map[string]interface{}) *MockDiscordInteraction {
	t.Helper()
	return &MockDiscordInteraction{
		t:           t,
		CommandName: commandName,
		Options:     options,
		UserID:      "test-user-123",
		GuildID:     "test-guild-456",
		ChannelID:   "test-channel-789",
	}
}

// Respond simulates sending an interaction response.
func (m *MockDiscordInteraction) Respond(content string) {
	m.ResponseSent = true
	m.ResponseData = &MockInteractionResponse{
		Type:    4, // CHANNEL_MESSAGE_WITH_SOURCE
		Content: content,
	}
}

// RespondEphemeral simulates sending an ephemeral interaction response.
func (m *MockDiscordInteraction) RespondEphemeral(content string) {
	m.ResponseSent = true
	m.ResponseData = &MockInteractionResponse{
		Type:      4,
		Content:   content,
		Ephemeral: true,
	}
}

// RespondWithEmbed simulates sending an interaction response with an embed.
func (m *MockDiscordInteraction) RespondWithEmbed(embed MockEmbed) {
	m.ResponseSent = true
	m.ResponseData = &MockInteractionResponse{
		Type:   4,
		Embeds: []MockEmbed{embed},
	}
}

// GetOption retrieves an option value from the interaction.
func (m *MockDiscordInteraction) GetOption(name string) (interface{}, bool) {
	val, ok := m.Options[name]
	return val, ok
}

// VerifyResponse checks if a response was sent and optionally verifies content.
func (m *MockDiscordInteraction) VerifyResponse(t *testing.T, expectedContent string) {
	t.Helper()
	if !m.ResponseSent {
		t.Error("Expected interaction response to be sent, but it wasn't")
		return
	}
	if expectedContent != "" && m.ResponseData.Content != expectedContent {
		t.Errorf("Response content = %q, want %q", m.ResponseData.Content, expectedContent)
	}
}

// NewMockApplicationCommand creates a mock application command definition.
func NewMockApplicationCommand(name, description string) *MockApplicationCommand {
	return &MockApplicationCommand{
		Name:        name,
		Description: description,
		Options:     []MockCommandOption{},
	}
}

// AddOption adds an option to the command.
func (c *MockApplicationCommand) AddOption(name, description string, optionType int, required bool) {
	c.Options = append(c.Options, MockCommandOption{
		Name:        name,
		Description: description,
		Type:        optionType,
		Required:    required,
	})
}

// AddOptionWithChoices adds an option with choices to the command.
func (c *MockApplicationCommand) AddOptionWithChoices(name, description string, optionType int, required bool, choices []MockChoice) {
	c.Options = append(c.Options, MockCommandOption{
		Name:        name,
		Description: description,
		Type:        optionType,
		Required:    required,
		Choices:     choices,
	})
}
