package bot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// MockDiscordSession is a mock implementation of DiscordSession for testing.
type MockDiscordSession struct {
	mu sync.Mutex

	// Mock state
	IsOpen     bool
	UserID     string
	Username   string
	Commands   []*discordgo.ApplicationCommand
	Handlers   []interface{}
	
	// Call tracking
	OpenCalls               int
	CloseCalls              int
	AddHandlerCalls         int
	GetUserCalls            int
	CreateCommandCalls      int
	ListCommandsCalls       int
	DeleteCommandCalls      int
	RespondCalls            int
	EditResponseCalls       int
	
	// Error injection
	OpenError               error
	CloseError              error
	CreateCommandError      error
	ListCommandsError       error
	DeleteCommandError      error
	RespondError            error
	EditResponseError       error
	
	// Captured data for assertions
	LastInteractionRespond *discordgo.InteractionResponse
	LastWebhookEdit        *discordgo.WebhookEdit
	DeletedCommandIDs      []string
}

// NewMockDiscordSession creates a new mock Discord session for testing.
func NewMockDiscordSession() *MockDiscordSession {
	return &MockDiscordSession{
		UserID:            "mock-bot-id-123",
		Username:          "MockBot",
		Commands:          make([]*discordgo.ApplicationCommand, 0),
		Handlers:          make([]interface{}, 0),
		DeletedCommandIDs: make([]string, 0),
	}
}

func (m *MockDiscordSession) Open() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.OpenCalls++
	if m.OpenError != nil {
		return m.OpenError
	}
	m.IsOpen = true
	return nil
}

func (m *MockDiscordSession) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.CloseCalls++
	if m.CloseError != nil {
		return m.CloseError
	}
	m.IsOpen = false
	return nil
}

func (m *MockDiscordSession) AddHandler(handler interface{}) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.AddHandlerCalls++
	m.Handlers = append(m.Handlers, handler)
	
	// Return a function that removes the handler
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		for i, h := range m.Handlers {
			if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
				m.Handlers = append(m.Handlers[:i], m.Handlers[i+1:]...)
				break
			}
		}
	}
}

func (m *MockDiscordSession) GetUser() *discordgo.User {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.GetUserCalls++
	if m.UserID == "" {
		return nil
	}
	return &discordgo.User{
		ID:       m.UserID,
		Username: m.Username,
	}
}

func (m *MockDiscordSession) ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.CreateCommandCalls++
	if m.CreateCommandError != nil {
		return nil, m.CreateCommandError
	}
	
	// Create a copy with an ID
	registered := &discordgo.ApplicationCommand{
		ID:          fmt.Sprintf("cmd-%d", len(m.Commands)+1),
		Name:        cmd.Name,
		Description: cmd.Description,
		Options:     cmd.Options,
	}
	m.Commands = append(m.Commands, registered)
	return registered, nil
}

func (m *MockDiscordSession) ApplicationCommands(appID, guildID string) ([]*discordgo.ApplicationCommand, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.ListCommandsCalls++
	if m.ListCommandsError != nil {
		return nil, m.ListCommandsError
	}
	
	// Return a copy to prevent external modification
	result := make([]*discordgo.ApplicationCommand, len(m.Commands))
	copy(result, m.Commands)
	return result, nil
}

func (m *MockDiscordSession) ApplicationCommandDelete(appID, guildID, cmdID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.DeleteCommandCalls++
	if m.DeleteCommandError != nil {
		return m.DeleteCommandError
	}
	
	m.DeletedCommandIDs = append(m.DeletedCommandIDs, cmdID)
	
	// Remove from commands list
	for i, cmd := range m.Commands {
		if cmd.ID == cmdID {
			m.Commands = append(m.Commands[:i], m.Commands[i+1:]...)
			break
		}
	}
	
	return nil
}

func (m *MockDiscordSession) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.RespondCalls++
	m.LastInteractionRespond = resp
	
	if m.RespondError != nil {
		return m.RespondError
	}
	return nil
}

func (m *MockDiscordSession) InteractionResponseEdit(interaction *discordgo.Interaction, edit *discordgo.WebhookEdit) (*discordgo.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.EditResponseCalls++
	m.LastWebhookEdit = edit
	
	if m.EditResponseError != nil {
		return nil, m.EditResponseError
	}
	
	// Return a mock message
	return &discordgo.Message{
		ID:      "mock-message-id",
		Content: func() string {
			if edit.Content != nil {
				return *edit.Content
			}
			return ""
		}(),
	}, nil
}

// Helper methods for test assertions

func (m *MockDiscordSession) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.OpenCalls = 0
	m.CloseCalls = 0
	m.AddHandlerCalls = 0
	m.GetUserCalls = 0
	m.CreateCommandCalls = 0
	m.ListCommandsCalls = 0
	m.DeleteCommandCalls = 0
	m.RespondCalls = 0
	m.EditResponseCalls = 0
	
	m.DeletedCommandIDs = make([]string, 0)
	m.LastInteractionRespond = nil
	m.LastWebhookEdit = nil
}

func (m *MockDiscordSession) SetUser(id, username string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UserID = id
	m.Username = username
}

func (m *MockDiscordSession) GetCommandCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Commands)
}
