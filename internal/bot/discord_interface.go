package bot

import "github.com/bwmarrin/discordgo"

// DiscordSession abstracts the Discord session operations needed by the bot.
// This interface enables mocking of Discord interactions for testing.
type DiscordSession interface {
	// Connection management
	Open() error
	Close() error

	// Event handlers
	AddHandler(handler interface{}) func()

	// User information
	GetUser() *discordgo.User

	// Command management
	ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	ApplicationCommands(appID, guildID string) ([]*discordgo.ApplicationCommand, error)
	ApplicationCommandDelete(appID, guildID, cmdID string) error

	// Interaction handling
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error
	InteractionResponseEdit(interaction *discordgo.Interaction, edit *discordgo.WebhookEdit) (*discordgo.Message, error)
}

// discordgoSessionAdapter wraps a real discordgo.Session to implement DiscordSession.
type discordgoSessionAdapter struct {
	session *discordgo.Session
}

// NewDiscordgoSessionAdapter creates an adapter for a real discordgo.Session.
func NewDiscordgoSessionAdapter(session *discordgo.Session) DiscordSession {
	return &discordgoSessionAdapter{session: session}
}

func (a *discordgoSessionAdapter) Open() error {
	return a.session.Open()
}

func (a *discordgoSessionAdapter) Close() error {
	return a.session.Close()
}

func (a *discordgoSessionAdapter) AddHandler(handler interface{}) func() {
	return a.session.AddHandler(handler)
}

func (a *discordgoSessionAdapter) GetUser() *discordgo.User {
	if a.session.State != nil && a.session.State.User != nil {
		return a.session.State.User
	}
	return nil
}

func (a *discordgoSessionAdapter) ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	return a.session.ApplicationCommandCreate(appID, guildID, cmd)
}

func (a *discordgoSessionAdapter) ApplicationCommands(appID, guildID string) ([]*discordgo.ApplicationCommand, error) {
	return a.session.ApplicationCommands(appID, guildID)
}

func (a *discordgoSessionAdapter) ApplicationCommandDelete(appID, guildID, cmdID string) error {
	return a.session.ApplicationCommandDelete(appID, guildID, cmdID)
}

func (a *discordgoSessionAdapter) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return a.session.InteractionRespond(interaction, resp)
}

func (a *discordgoSessionAdapter) InteractionResponseEdit(interaction *discordgo.Interaction, edit *discordgo.WebhookEdit) (*discordgo.Message, error) {
	return a.session.InteractionResponseEdit(interaction, edit)
}
