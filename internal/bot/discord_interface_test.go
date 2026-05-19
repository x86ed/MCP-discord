package bot

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestDiscordgoSessionAdapter_Open tests the Open method delegation
func TestDiscordgoSessionAdapter_Open(t *testing.T) {
	mock := NewMockDiscordSession()

	err := mock.Open()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mock.OpenCalls != 1 {
		t.Errorf("Expected 1 Open call, got %d", mock.OpenCalls)
	}
}

// TestDiscordgoSessionAdapter_OpenError tests Open with error
func TestDiscordgoSessionAdapter_OpenError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.OpenError = errors.New("connection failed")

	err := mock.Open()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.OpenCalls != 1 {
		t.Errorf("Expected 1 Open call, got %d", mock.OpenCalls)
	}
}

// TestDiscordgoSessionAdapter_Close tests the Close method delegation
func TestDiscordgoSessionAdapter_Close(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.IsOpen = true

	err := mock.Close()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mock.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mock.CloseCalls)
	}

	if mock.IsOpen {
		t.Error("Expected session to be closed")
	}
}

// TestDiscordgoSessionAdapter_CloseError tests Close with error
func TestDiscordgoSessionAdapter_CloseError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.CloseError = errors.New("close failed")

	err := mock.Close()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mock.CloseCalls)
	}
}

// TestDiscordgoSessionAdapter_GetUser tests the GetUser method delegation
func TestDiscordgoSessionAdapter_GetUser(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.SetUser("user-123", "testuser")

	user := mock.GetUser()

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.ID != "user-123" {
		t.Errorf("Expected user ID %s, got %s", "user-123", user.ID)
	}

	if mock.GetUserCalls != 1 {
		t.Errorf("Expected 1 GetUser call, got %d", mock.GetUserCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandCreate tests command creation delegation
func TestDiscordgoSessionAdapter_ApplicationCommandCreate(t *testing.T) {
	mock := NewMockDiscordSession()

	cmd := &discordgo.ApplicationCommand{
		Name:        "test",
		Description: "Test command",
	}

	result, err := mock.ApplicationCommandCreate("app-123", "guild-123", cmd)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected command result, got nil")
	}

	if result.Name != "test" {
		t.Errorf("Expected command name 'test', got %s", result.Name)
	}

	if mock.CreateCommandCalls != 1 {
		t.Errorf("Expected 1 CreateCommand call, got %d", mock.CreateCommandCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandCreateError tests command creation with error
func TestDiscordgoSessionAdapter_ApplicationCommandCreateError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.CreateCommandError = errors.New("create failed")

	cmd := &discordgo.ApplicationCommand{Name: "test"}

	_, err := mock.ApplicationCommandCreate("app-123", "guild-123", cmd)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.CreateCommandCalls != 1 {
		t.Errorf("Expected 1 CreateCommand call, got %d", mock.CreateCommandCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommands tests listing commands
func TestDiscordgoSessionAdapter_ApplicationCommands(t *testing.T) {
	expectedCommands := []*discordgo.ApplicationCommand{
		{ID: "cmd-1", Name: "command1"},
		{ID: "cmd-2", Name: "command2"},
	}
	mock := NewMockDiscordSession()
	mock.Commands = expectedCommands

	commands, err := mock.ApplicationCommands("app-123", "guild-123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	if mock.ListCommandsCalls != 1 {
		t.Errorf("Expected 1 ListCommands call, got %d", mock.ListCommandsCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandsEmpty tests listing with no commands
func TestDiscordgoSessionAdapter_ApplicationCommandsEmpty(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.Commands = []*discordgo.ApplicationCommand{}

	commands, err := mock.ApplicationCommands("app-123", "guild-123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(commands) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(commands))
	}

	if mock.ListCommandsCalls != 1 {
		t.Errorf("Expected 1 ListCommands call, got %d", mock.ListCommandsCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandsError tests listing with error
func TestDiscordgoSessionAdapter_ApplicationCommandsError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.ListCommandsError = errors.New("list failed")

	_, err := mock.ApplicationCommands("app-123", "guild-123")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.ListCommandsCalls != 1 {
		t.Errorf("Expected 1 ListCommands call, got %d", mock.ListCommandsCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandDelete tests command deletion
func TestDiscordgoSessionAdapter_ApplicationCommandDelete(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.Commands = []*discordgo.ApplicationCommand{
		{ID: "cmd-123", Name: "test"},
	}

	err := mock.ApplicationCommandDelete("app-123", "guild-123", "cmd-123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mock.DeleteCommandCalls != 1 {
		t.Errorf("Expected 1 DeleteCommand call, got %d", mock.DeleteCommandCalls)
	}
}

// TestDiscordgoSessionAdapter_ApplicationCommandDeleteError tests deletion with error
func TestDiscordgoSessionAdapter_ApplicationCommandDeleteError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.DeleteCommandError = errors.New("delete failed")

	err := mock.ApplicationCommandDelete("app-123", "guild-123", "cmd-123")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.DeleteCommandCalls != 1 {
		t.Errorf("Expected 1 DeleteCommand call, got %d", mock.DeleteCommandCalls)
	}
}

// TestDiscordgoSessionAdapter_InteractionRespond tests interaction response
func TestDiscordgoSessionAdapter_InteractionRespond(t *testing.T) {
	mock := NewMockDiscordSession()

	interaction := &discordgo.Interaction{ID: "interaction-123"}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}

	err := mock.InteractionRespond(interaction, response)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mock.RespondCalls != 1 {
		t.Errorf("Expected 1 Respond call, got %d", mock.RespondCalls)
	}
}

// TestDiscordgoSessionAdapter_InteractionRespondError tests respond with error
func TestDiscordgoSessionAdapter_InteractionRespondError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.RespondError = errors.New("respond failed")

	interaction := &discordgo.Interaction{ID: "interaction-123"}
	response := &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource}

	err := mock.InteractionRespond(interaction, response)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.RespondCalls != 1 {
		t.Errorf("Expected 1 Respond call, got %d", mock.RespondCalls)
	}
}

// TestDiscordgoSessionAdapter_InteractionResponseEdit tests editing interaction response
func TestDiscordgoSessionAdapter_InteractionResponseEdit(t *testing.T) {
	mock := NewMockDiscordSession()

	interaction := &discordgo.Interaction{ID: "interaction-123"}
	content := "Updated content"
	edit := &discordgo.WebhookEdit{Content: &content}

	message, err := mock.InteractionResponseEdit(interaction, edit)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if message == nil {
		t.Fatal("Expected message result, got nil")
	}

	if mock.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mock.EditResponseCalls)
	}
}

// TestDiscordgoSessionAdapter_InteractionResponseEditError tests edit with error
func TestDiscordgoSessionAdapter_InteractionResponseEditError(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.EditResponseError = errors.New("edit failed")

	interaction := &discordgo.Interaction{ID: "interaction-123"}
	content := "Test"
	edit := &discordgo.WebhookEdit{Content: &content}

	_, err := mock.InteractionResponseEdit(interaction, edit)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mock.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mock.EditResponseCalls)
	}
}

// TestDiscordgoSessionAdapter_AddHandler tests handler registration
func TestDiscordgoSessionAdapter_AddHandler(t *testing.T) {
	mock := NewMockDiscordSession()

	handler := func() {}

	removeFunc := mock.AddHandler(handler)

	if removeFunc == nil {
		t.Error("Expected remove function, got nil")
	}

	if mock.AddHandlerCalls != 1 {
		t.Errorf("Expected 1 AddHandler call, got %d", mock.AddHandlerCalls)
	}

	if len(mock.Handlers) != 1 {
		t.Errorf("Expected 1 handler registered, got %d", len(mock.Handlers))
	}

	// Test the remove function
	removeFunc()

	if len(mock.Handlers) != 0 {
		t.Errorf("Expected 0 handlers after removal, got %d", len(mock.Handlers))
	}
}
