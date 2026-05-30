package bot

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestDiscordgoSessionAdapter tests the adapter wrapping a real session
func TestDiscordgoSessionAdapter(t *testing.T) {
	// Create a real discordgo session (won't actually connect)
	dg, err := discordgo.New("Bot test-token")
	if err != nil {
		t.Fatalf("Failed to create discordgo session: %v", err)
	}

	// Wrap it with our adapter
	adapter := NewDiscordgoSessionAdapter(dg)

	// Test that the adapter implements the interface
	var _ DiscordSession = adapter

	// Test GetUser when session state is nil (safe check)
	user := adapter.GetUser()
	if user != nil {
		t.Error("Expected nil user when state is uninitialized")
	}

	// Test GetUser when state is initialized with a user
	dg.State = discordgo.NewState()
	dg.State.User = &discordgo.User{
		ID:       "123",
		Username: "testbot",
	}
	adapter = NewDiscordgoSessionAdapter(dg)
	user = adapter.GetUser()
	if user == nil {
		t.Error("Expected user when state is initialized")
	}
	if user.ID != "123" || user.Username != "testbot" {
		t.Error("GetUser returned wrong user data")
	}
}

// TestDiscordgoSessionAdapter_Methods tests adapter method delegation
func TestDiscordgoSessionAdapter_Methods(t *testing.T) {
	dg, err := discordgo.New("Bot test-token")
	if err != nil {
		t.Fatalf("Failed to create discordgo session: %v", err)
	}

	adapter := NewDiscordgoSessionAdapter(dg).(*discordgoSessionAdapter)

	// Verify the adapter has the session
	if adapter.session != dg {
		t.Error("Adapter should wrap the provided session")
	}

	// Test AddHandler (won't actually register since not connected)
	removeFunc := adapter.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Handler function
	})

	if removeFunc == nil {
		t.Error("AddHandler should return a remove function")
	}

	// Call the remove function (should not panic)
	removeFunc()
}
