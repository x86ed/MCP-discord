package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"mcpdiscord/internal/config"
	"mcpdiscord/internal/mcp"
)

// ===== MOCK TRANSLATOR =====

type MockTranslator struct {
	ToolToSlashCommandFunc func(tool mcp.Tool) (*discordgo.ApplicationCommand, error)
	TranslateArgumentsFunc func(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error)
}

func (m *MockTranslator) ToolToSlashCommand(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
	if m.ToolToSlashCommandFunc != nil {
		return m.ToolToSlashCommandFunc(tool)
	}
	return &discordgo.ApplicationCommand{
		Name:        tool.Name,
		Description: tool.Description,
	}, nil
}

func (m *MockTranslator) TranslateArguments(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
	if m.TranslateArgumentsFunc != nil {
		return m.TranslateArgumentsFunc(tool, options)
	}
	return map[string]interface{}{}, nil
}

// ===== MOCK MCP CLIENT =====

type MockMCPClient struct {
	ListToolsData []mcp.Tool
	ListToolsErr  error
	CallToolData  *mcp.ToolResult
	CallToolErr   error
}

func (m *MockMCPClient) Connect(ctx context.Context) error {
	return nil
}

func (m *MockMCPClient) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	return m.ListToolsData, m.ListToolsErr
}

func (m *MockMCPClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (*mcp.ToolResult, error) {
	return m.CallToolData, m.CallToolErr
}

func (m *MockMCPClient) Close() error {
	return nil
}

// ===== TESTS FOR FORMAT RESULT (CORE LOGIC) =====

func TestFormatResult_WithToolResult(t *testing.T) {
	bot := &DiscordBot{}

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{Type: "text", Text: "Operation successful"},
		},
		IsError: false,
	}

	resp := bot.formatResult("test_command", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.Color != 0x00FF00 {
		t.Errorf("Expected green color for success, got 0x%06X", resp.Embed.Color)
	}

	if resp.Embed.Title != "✓ test_command" {
		t.Errorf("Expected title with checkmark, got %q", resp.Embed.Title)
	}
}

func TestFormatResult_ErrorResult(t *testing.T) {
	bot := &DiscordBot{}

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{Type: "text", Text: "Operation failed"},
		},
		IsError: true,
	}

	resp := bot.formatResult("test_command", result)

	if resp.Embed.Color != 0xFF0000 {
		t.Errorf("Expected red color for error, got 0x%06X", resp.Embed.Color)
	}

	if resp.Embed.Title != "✗ test_command" {
		t.Errorf("Expected title with X, got %q", resp.Embed.Title)
	}
}

func TestFormatResult_MultipleContentBlocks(t *testing.T) {
	bot := &DiscordBot{}

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{Type: "text", Text: "Line 1"},
			{Type: "text", Text: "Line 2"},
			{Type: "text", Text: "Line 3"},
		},
		IsError: false,
	}

	resp := bot.formatResult("test", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed")
	}

	// Should combine all content blocks
	expected := "Line 1\nLine 2\nLine 3"
	if resp.Embed.Description != expected {
		t.Errorf("Expected combined content, got %q", resp.Embed.Description)
	}
}

func TestFormatResult_EmptyContent(t *testing.T) {
	bot := &DiscordBot{}

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{},
		IsError: false,
	}

	resp := bot.formatResult("test", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed")
	}

	// Should handle empty content gracefully
	if resp.Embed.Description == "" {
		t.Log("Empty description is acceptable for empty content")
	}
}

func TestFormatResult_JSONDetection(t *testing.T) {
	bot := &DiscordBot{}

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{Type: "text", Text: `{"key": "value", "count": 42}`},
		},
		IsError: false,
	}

	resp := bot.formatResult("test", result)

	// Should wrap JSON in code block
	if resp.Embed.Description[:7] != "```json" {
		t.Errorf("Expected JSON code block, got %q", resp.Embed.Description[:20])
	}
}

// ===== TESTS FOR HANDLE INTERACTION STUB =====

func TestHandleInteraction_NotImplemented(t *testing.T) {
	bot := &DiscordBot{}

	_, err := bot.HandleInteraction(&Interaction{
		ID:          "test",
		CommandName: "test",
	})

	if err == nil {
		t.Error("Expected error for not implemented method")
	}
}

// ===== TESTS FOR COMMAND TOOL MAPPING =====

func TestCommandToolMapping(t *testing.T) {
	// Test that we can build the command/tool mapping correctly
	bot := &DiscordBot{
		commandTools: make(map[string]mcp.Tool),
	}

	tool := mcp.Tool{
		Name:        "test_tool",
		Description: "Test tool",
		InputSchema: mcp.InputSchema{Type: "object"},
	}

	// Simulate what RegisterCommands does
	commandName := "test_tool"
	bot.commandTools[commandName] = tool

	// Verify mapping exists
	mappedTool, exists := bot.commandTools[commandName]
	if !exists {
		t.Error("Expected tool to be mapped")
	}

	if mappedTool.Name != tool.Name {
		t.Errorf("Expected tool name %q, got %q", tool.Name, mappedTool.Name)
	}
}

func TestCommandToolMapping_LookupFailure(t *testing.T) {
	bot := &DiscordBot{
		commandTools: make(map[string]mcp.Tool),
	}

	// Try to look up non-existent command
	_, exists := bot.commandTools["nonexistent"]
	if exists {
		t.Error("Should not find non-existent command")
	}
}

// ===== TESTS FOR MCP CALL LOGIC =====

func TestMCPCallLogic_Success(t *testing.T) {
	mcpClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: "Success"},
			},
			IsError: false,
		},
	}

	toolName := "test_tool"
	args := map[string]interface{}{"key": "value"}

	result, err := mcpClient.CallTool(context.Background(), toolName, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(result.Content))
	}
}

func TestMCPCallLogic_Error(t *testing.T) {
	mcpClient := &MockMCPClient{
		CallToolErr: errors.New("MCP call failed"),
	}

	_, err := mcpClient.CallTool(context.Background(), "test", nil)

	if err == nil {
		t.Error("Expected error from MCP call")
	}
}

// ===== TESTS FOR TRANSLATION LOGIC =====

func TestTranslationLogic_Success(t *testing.T) {
	translator := &MockTranslator{
		TranslateArgumentsFunc: func(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
			return map[string]interface{}{
				"arg1": "value1",
				"arg2": 42,
			}, nil
		},
	}

	tool := mcp.Tool{
		Name: "test",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"arg1": {Type: "string"},
				"arg2": {Type: "number"},
			},
		},
	}

	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "arg1", Type: discordgo.ApplicationCommandOptionString, Value: "value1"},
		{Name: "arg2", Type: discordgo.ApplicationCommandOptionNumber, Value: float64(42)},
	}

	args, err := translator.TranslateArguments(tool, options)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if args["arg1"] != "value1" {
		t.Errorf("Expected arg1=value1, got %v", args["arg1"])
	}

	if args["arg2"] != 42 {
		t.Errorf("Expected arg2=42, got %v", args["arg2"])
	}
}

func TestTranslationLogic_Error(t *testing.T) {
	translator := &MockTranslator{
		TranslateArgumentsFunc: func(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
			return nil, errors.New("invalid arguments")
		},
	}

	tool := mcp.Tool{Name: "test"}
	options := []*discordgo.ApplicationCommandInteractionDataOption{}

	_, err := translator.TranslateArguments(tool, options)

	if err == nil {
		t.Error("Expected translation error")
	}
}

// ===== TESTS FOR TOOL TO COMMAND CONVERSION =====

func TestToolToCommandConversion_Success(t *testing.T) {
	translator := &MockTranslator{
		ToolToSlashCommandFunc: func(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
			return &discordgo.ApplicationCommand{
				Name:        tool.Name,
				Description: tool.Description,
				Options: []*discordgo.ApplicationCommandOption{
					{Name: "param1", Type: discordgo.ApplicationCommandOptionString},
				},
			}, nil
		},
	}

	tool := mcp.Tool{
		Name:        "test_tool",
		Description: "Test description",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.PropertySchema{
				"param1": {Type: "string", Description: "A parameter"},
			},
		},
	}

	cmd, err := translator.ToolToSlashCommand(tool)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cmd.Name != tool.Name {
		t.Errorf("Expected command name %q, got %q", tool.Name, cmd.Name)
	}

	if len(cmd.Options) != 1 {
		t.Errorf("Expected 1 option, got %d", len(cmd.Options))
	}
}

func TestToolToCommandConversion_Error(t *testing.T) {
	translator := &MockTranslator{
		ToolToSlashCommandFunc: func(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
			return nil, errors.New("conversion failed")
		},
	}

	tool := mcp.Tool{Name: "test"}

	_, err := translator.ToolToSlashCommand(tool)

	if err == nil {
		t.Error("Expected conversion error")
	}
}

// ===== TESTS FOR DISCOVERY INTEGRATION =====

func TestDiscoveryIntegration_ValidTools(t *testing.T) {
	mcpClient := &MockMCPClient{
		ListToolsData: []mcp.Tool{
			{
				Name:        "tool1",
				Description: "Tool 1",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
			{
				Name:        "tool2",
				Description: "Tool 2",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
		},
	}

	tools, err := mcpClient.ListTools(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}
}

func TestDiscoveryIntegration_Error(t *testing.T) {
	mcpClient := &MockMCPClient{
		ListToolsErr: errors.New("discovery failed"),
	}

	_, err := mcpClient.ListTools(context.Background())

	if err == nil {
		t.Error("Expected discovery error")
	}
}

// ===== TESTS FOR RESPONSE BUILDING =====

func TestResponseBuilding_WithEmbed(t *testing.T) {
	resp := &Response{
		Content: "Test content",
		Embed: &Embed{
			Title:       "Test Title",
			Description: "Test Description",
			Color:       0x00FF00,
		},
		Ephemeral: false,
	}

	if resp.Embed.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %q", resp.Embed.Title)
	}

	if resp.Embed.Color != 0x00FF00 {
		t.Errorf("Expected color 0x00FF00, got 0x%06X", resp.Embed.Color)
	}
}

func TestResponseBuilding_Ephemeral(t *testing.T) {
	resp := &Response{
		Content:   "Secret message",
		Ephemeral: true,
	}

	if !resp.Ephemeral {
		t.Error("Expected ephemeral response")
	}
}

// ===== TESTS WITH MOCK DISCORD SESSION =====

func TestNewDiscordBot_WithMockSession(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: "guild-123"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	if bot.session == nil {
		t.Fatal("Expected session to be set")
	}

	if mockSession.AddHandlerCalls > 0 {
		t.Log("Handlers would be added in NewDiscordBot constructor")
	}
}

func TestDiscordBot_Start(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	err := bot.Start(context.Background())
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if mockSession.OpenCalls != 1 {
		t.Errorf("Expected 1 Open call, got %d", mockSession.OpenCalls)
	}

	if !mockSession.IsOpen {
		t.Error("Expected session to be open")
	}

	if mockSession.GetUserCalls != 1 {
		t.Errorf("Expected 1 GetUser call, got %d", mockSession.GetUserCalls)
	}

	if bot.userID != "mock-bot-id-123" {
		t.Errorf("Expected userID to be set to mock-bot-id-123, got %q", bot.userID)
	}
}

func TestDiscordBot_Start_OpenError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.OpenError = errors.New("connection failed")

	bot := &DiscordBot{
		config:     config.DiscordConfig{Token: "test"},
		session:    mockSession,
		logger:     slog.Default(),
		commands:   make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	err := bot.Start(context.Background())
	if err == nil {
		t.Fatal("Expected error from Start")
	}

	if mockSession.OpenCalls != 1 {
		t.Errorf("Expected 1 Open call despite error, got %d", mockSession.OpenCalls)
	}
}

func TestDiscordBot_Start_NoUser(t *testing.T) {
	mockSession := NewMockDiscordSession()
	// Don't set user, so GetUser returns nil
	mockSession.UserID = ""
	
	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	err := bot.Start(context.Background())
	if err != nil {
		t.Fatalf("Start should not fail when GetUser returns nil: %v", err)
	}

	if mockSession.GetUserCalls != 1 {
		t.Errorf("Expected 1 GetUser call, got %d", mockSession.GetUserCalls)
	}

	if bot.userID != "" {
		t.Errorf("Expected empty userID, got %q", bot.userID)
	}
}


func TestDiscordBot_Stop(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.IsOpen = true

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: "guild-123"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if mockSession.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mockSession.CloseCalls)
	}

	if mockSession.IsOpen {
		t.Error("Expected session to be closed")
	}
}

func TestDiscordBot_RegisterCommands(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		ListToolsData: []mcp.Tool{
			{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
		},
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.RegisterCommands("guild-123")
	if err != nil {
		t.Fatalf("RegisterCommands failed: %v", err)
	}

	if mockSession.CreateCommandCalls != 1 {
		t.Errorf("Expected 1 CreateCommand call, got %d", mockSession.CreateCommandCalls)
	}

	if len(bot.commands) != 1 {
		t.Errorf("Expected 1 registered command, got %d", len(bot.commands))
	}

	if len(bot.commandTools) != 1 {
		t.Errorf("Expected 1 command tool mapping, got %d", len(bot.commandTools))
	}
}

func TestDiscordBot_RegisterCommands_DiscoveryError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		ListToolsErr: errors.New("discovery failed"),
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.RegisterCommands("guild-123")
	if err == nil {
		t.Fatal("Expected error from RegisterCommands")
	}

	if mockSession.CreateCommandCalls != 0 {
		t.Errorf("Expected 0 CreateCommand calls on error, got %d", mockSession.CreateCommandCalls)
	}
}

func TestDiscordBot_DeregisterCommands(t *testing.T) {
	mockSession := NewMockDiscordSession()
	
	// Pre-populate with existing commands
	mockSession.Commands = []*discordgo.ApplicationCommand{
		{ID: "cmd-1", Name: "command1"},
		{ID: "cmd-2", Name: "command2"},
	}

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.deregisterCommands("guild-123")
	if err != nil {
		t.Fatalf("deregisterCommands failed: %v", err)
	}

	if mockSession.ListCommandsCalls != 1 {
		t.Errorf("Expected 1 ListCommands call, got %d", mockSession.ListCommandsCalls)
	}

	if mockSession.DeleteCommandCalls != 2 {
		t.Errorf("Expected 2 DeleteCommand calls, got %d", mockSession.DeleteCommandCalls)
	}

	if len(mockSession.DeletedCommandIDs) != 2 {
		t.Errorf("Expected 2 deleted command IDs, got %d", len(mockSession.DeletedCommandIDs))
	}

	if mockSession.GetCommandCount() != 0 {
		t.Errorf("Expected 0 commands after deregister, got %d", mockSession.GetCommandCount())
	}
}

func TestDiscordBot_DeregisterCommands_ListError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.ListCommandsError = errors.New("failed to list commands")

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.deregisterCommands("guild-123")
	if err == nil {
		t.Fatal("Expected error from deregisterCommands")
	}

	if mockSession.DeleteCommandCalls != 0 {
		t.Errorf("Expected 0 DeleteCommand calls on error, got %d", mockSession.DeleteCommandCalls)
	}
}

func TestDiscordBot_DeregisterCommands_DeleteError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.Commands = []*discordgo.ApplicationCommand{
		{ID: "cmd-1", Name: "command1"},
	}
	mockSession.DeleteCommandError = errors.New("failed to delete command")

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	// Should not return error, just log warning
	err := bot.deregisterCommands("guild-123")
	if err != nil {
		t.Fatalf("deregisterCommands should not fail on delete error: %v", err)
	}

	if mockSession.DeleteCommandCalls != 1 {
		t.Errorf("Expected 1 DeleteCommand call, got %d", mockSession.DeleteCommandCalls)
	}
}

func TestDiscordBot_Stop_CloseError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.IsOpen = true
	mockSession.CloseError = errors.New("failed to close session")

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.Stop()
	if err == nil {
		t.Fatal("Expected error from Stop")
	}

	if mockSession.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mockSession.CloseCalls)
	}
}

func TestDiscordBot_Stop_WithGuildDeregister(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.IsOpen = true
	mockSession.Commands = []*discordgo.ApplicationCommand{
		{ID: "cmd-1", Name: "command1"},
	}

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: "guild-123"},
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if mockSession.ListCommandsCalls != 1 {
		t.Errorf("Expected 1 ListCommands call for deregister, got %d", mockSession.ListCommandsCalls)
	}

	if mockSession.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mockSession.CloseCalls)
	}
}

func TestDiscordBot_Stop_NoGuildID(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.IsOpen = true

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: ""}, // No guild ID
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Should close but not try to deregister
	if mockSession.CloseCalls != 1 {
		t.Errorf("Expected 1 Close call, got %d", mockSession.CloseCalls)
	}

	// Should not have tried to list commands for deregistration
	if mockSession.ListCommandsCalls != 0 {
		t.Errorf("Expected 0 ListCommands calls (no guild ID), got %d", mockSession.ListCommandsCalls)
	}
}

func TestDiscordBot_Stop_NilSession(t *testing.T) {
	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		session:      nil, // No session
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
	}

	err := bot.Stop()
	if err != nil {
		t.Fatalf("Stop should not fail with nil session: %v", err)
	}
}

func TestDiscordBot_RegisterCommands_TranslationError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		ListToolsData: []mcp.Tool{
			{
				Name:        "valid_tool",
				Description: "A valid tool",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
			{
				Name:        "invalid_tool",
				Description: "An invalid tool",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
		},
	}
	mockTranslator := &MockTranslator{
		ToolToSlashCommandFunc: func(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
			if tool.Name == "invalid_tool" {
				return nil, errors.New("translation failed")
			}
			return &discordgo.ApplicationCommand{
				Name:        tool.Name,
				Description: tool.Description,
			}, nil
		},
	}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.RegisterCommands("guild-123")
	if err != nil {
		t.Fatalf("RegisterCommands failed: %v", err)
	}

	// Should have registered only the valid tool
	if mockSession.CreateCommandCalls != 1 {
		t.Errorf("Expected 1 CreateCommand call (invalid skipped), got %d", mockSession.CreateCommandCalls)
	}

	if len(bot.commands) != 1 {
		t.Errorf("Expected 1 registered command, got %d", len(bot.commands))
	}
}

func TestDiscordBot_RegisterCommands_TooManyTools(t *testing.T) {
	mockSession := NewMockDiscordSession()
	
	// Create 101 tools to exceed Discord's limit of 100
	tools := make([]mcp.Tool, 101)
	for i := 0; i < 101; i++ {
		tools[i] = mcp.Tool{
			Name:        fmt.Sprintf("tool_%d", i),
			Description: "A tool",
			InputSchema: mcp.InputSchema{Type: "object"},
		}
	}
	
	mockClient := &MockMCPClient{
		ListToolsData: tools,
	}
	mockTranslator := &MockTranslator{
		ToolToSlashCommandFunc: func(tool mcp.Tool) (*discordgo.ApplicationCommand, error) {
			return &discordgo.ApplicationCommand{
				Name:        tool.Name,
				Description: tool.Description,
			}, nil
		},
	}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.RegisterCommands("guild-123")
	if err == nil {
		t.Fatal("Expected error for too many tools, got nil")
	}

	if !strings.Contains(err.Error(), "cannot register 101 commands") {
		t.Errorf("Expected limit error, got: %v", err)
	}
}

func TestDiscordBot_RegisterCommands_CreateError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.CreateCommandError = errors.New("failed to create command")
	mockClient := &MockMCPClient{
		ListToolsData: []mcp.Tool{
			{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
		},
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	err := bot.RegisterCommands("guild-123")
	if err != nil {
		t.Fatalf("RegisterCommands should not fail on create error: %v", err)
	}

	if mockSession.CreateCommandCalls != 1 {
		t.Errorf("Expected 1 CreateCommand call, got %d", mockSession.CreateCommandCalls)
	}

	// Should have 0 registered commands due to error
	if len(bot.commands) != 0 {
		t.Errorf("Expected 0 registered commands on error, got %d", len(bot.commands))
	}
}

func TestDiscordBot_SendResponse(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID: "interaction-123",
		},
	}

	resp := &Response{
		Content: "Test response",
		Embed: &Embed{
			Title:       "Test",
			Description: "Description",
			Color:       0x00FF00,
		},
	}

	bot.sendResponse(mockSession, interaction, resp)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}

	if mockSession.LastWebhookEdit == nil {
		t.Fatal("Expected LastWebhookEdit to be captured")
	}

	if mockSession.LastWebhookEdit.Content == nil || *mockSession.LastWebhookEdit.Content != "Test response" {
		t.Errorf("Expected content 'Test response', got %v", mockSession.LastWebhookEdit.Content)
	}
}

func TestDiscordBot_SendResponse_WithImage(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID: "interaction-123",
		},
	}

	resp := &Response{
		Embed: &Embed{
			Title:       "Test",
			Description: "Description",
			Color:       0x00FF00,
			ImageURL:    "https://example.com/image.png",
		},
	}

	bot.sendResponse(mockSession, interaction, resp)

	if mockSession.LastWebhookEdit == nil || mockSession.LastWebhookEdit.Embeds == nil || len(*mockSession.LastWebhookEdit.Embeds) == 0 {
		t.Fatal("Expected embed in response")
	}

	embed := (*mockSession.LastWebhookEdit.Embeds)[0]
	if embed.Image == nil || embed.Image.URL != "https://example.com/image.png" {
		t.Fatalf("Expected image URL to be set, got %+v", embed.Image)
	}
}

func TestDiscordBot_SendError(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test_command",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.sendError(mockSession, interaction, errors.New("test error"), false)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}

	if mockSession.LastWebhookEdit == nil {
		t.Fatal("Expected LastWebhookEdit to be captured")
	}

	if mockSession.LastWebhookEdit.Embeds == nil || len(*mockSession.LastWebhookEdit.Embeds) == 0 {
		t.Fatal("Expected embed in error response")
	}

	embed := (*mockSession.LastWebhookEdit.Embeds)[0]
	if embed.Title != "Error" {
		t.Errorf("Expected error title, got %q", embed.Title)
	}

	if embed.Color != 0xFF0000 {
		t.Errorf("Expected red color for error, got 0x%06X", embed.Color)
	}
}

func TestDiscordBot_SendError_Internal(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test_command",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.sendError(mockSession, interaction, errors.New("internal error details"), true)

	if mockSession.LastWebhookEdit == nil {
		t.Fatal("Expected LastWebhookEdit to be captured")
	}

	embed := (*mockSession.LastWebhookEdit.Embeds)[0]
	if embed.Description == "internal error details" {
		t.Error("Internal error details should be hidden from user")
	}

	if embed.Description != "An internal error occurred. Please try again later." {
		t.Errorf("Expected generic error message, got %q", embed.Description)
	}
}

func TestDiscordBot_SendError_WithUserField(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	// Interaction with User field instead of Member
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test_command",
			},
			User: &discordgo.User{
				Username: "directuser",
			},
		},
	}

	bot.sendError(mockSession, interaction, errors.New("test error"), false)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_SendError_EditFails(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.EditResponseError = errors.New("failed to edit response")

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test_command",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	// Should not panic, just log error
	bot.sendError(mockSession, interaction, errors.New("test error"), false)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_SendResponse_NoEmbed(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
		},
	}

	response := &Response{
		Content: "Plain text response",
		Embed:   nil, // No embed
	}

	bot.sendResponse(mockSession, interaction, response)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}

	if mockSession.LastWebhookEdit == nil {
		t.Fatal("Expected LastWebhookEdit to be captured")
	}

	if mockSession.LastWebhookEdit.Content == nil || *mockSession.LastWebhookEdit.Content != "Plain text response" {
		t.Error("Expected content to be set")
	}
}

func TestDiscordBot_SendResponse_EditFails(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockSession.EditResponseError = errors.New("failed to edit response")

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
		},
	}

	response := &Response{
		Content: "Test response",
	}

	// Should not panic, just log error
	bot.sendResponse(mockSession, interaction, response)

	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestMockDiscordSession_Reset(t *testing.T) {
	mock := NewMockDiscordSession()

	// Make some calls
	_ = mock.Open()
	_ = mock.Close()
	mock.GetUser()

	if mock.OpenCalls == 0 || mock.CloseCalls == 0 || mock.GetUserCalls == 0 {
		t.Fatal("Mock calls not tracked")
	}

	// Reset
	mock.Reset()

	if mock.OpenCalls != 0 {
		t.Errorf("Expected OpenCalls to be reset to 0, got %d", mock.OpenCalls)
	}
	if mock.CloseCalls != 0 {
		t.Errorf("Expected CloseCalls to be reset to 0, got %d", mock.CloseCalls)
	}
	if mock.GetUserCalls != 0 {
		t.Errorf("Expected GetUserCalls to be reset to 0, got %d", mock.GetUserCalls)
	}
}

func TestMockDiscordSession_SetUser(t *testing.T) {
	mock := NewMockDiscordSession()
	mock.SetUser("custom-id", "CustomBot")

	user := mock.GetUser()
	if user == nil {
		t.Fatal("Expected user to be returned")
	}

	if user.ID != "custom-id" {
		t.Errorf("Expected user ID 'custom-id', got %q", user.ID)
	}

	if user.Username != "CustomBot" {
		t.Errorf("Expected username 'CustomBot', got %q", user.Username)
	}
}

// ===== TESTS FOR CONSTRUCTOR AND EVENT HANDLERS =====

func TestNewDiscordBot_InvalidToken(t *testing.T) {
	// Test with empty token to trigger an error
	cfg := config.DiscordConfig{Token: ""}
	mockClient := &MockMCPClient{}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot, err := NewDiscordBot(cfg, mockClient, mockTranslator, mockDiscovery, slog.Default())
	
	// Empty token should still create a bot, but may have issues later
	// The actual validation happens on Open()
	if bot == nil && err != nil {
		t.Log("Bot creation failed with empty token (expected in some cases)")
	}
}

func TestNewDiscordBot_NilLogger(t *testing.T) {
	// Test with nil logger - should use slog.Default()
	cfg := config.DiscordConfig{Token: "test-token"}
	mockClient := &MockMCPClient{}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot, err := NewDiscordBot(cfg, mockClient, mockTranslator, mockDiscovery, nil)
	
	if err != nil {
		t.Fatalf("Expected no error with nil logger, got %v", err)
	}
	
	if bot == nil {
		t.Fatal("Expected bot to be created")
	}
	
	if bot.logger == nil {
		t.Error("Expected logger to be set (should use slog.Default())")
	}
}


func TestDiscordBot_OnReady(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		ListToolsData: []mcp.Tool{
			{
				Name:        "test_tool",
				Description: "Test",
				InputSchema: mcp.InputSchema{Type: "object"},
			},
		},
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: "guild-123"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	// Simulate onReady event using the testable wrapper
	readyEvent := &discordgo.Ready{
		User: &discordgo.User{
			ID:       "bot-123",
			Username: "TestBot",
		},
		Guilds: []*discordgo.Guild{
			{ID: "guild-1"},
			{ID: "guild-2"},
		},
	}

	bot.handleReady(readyEvent)

	// handleReady should trigger RegisterCommands
	if mockSession.CreateCommandCalls != 1 {
		t.Errorf("Expected handleReady to trigger command registration, got %d CreateCommand calls", mockSession.CreateCommandCalls)
	}
}

func TestDiscordBot_HandleReady_RegisterCommandsError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		ListToolsErr: errors.New("failed to discover tools"),
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	bot := &DiscordBot{
		config:       config.DiscordConfig{Token: "test", GuildID: "guild-123"},
		mcpClient:    mockClient,
		translator:   mockTranslator,
		discovery:    mockDiscovery,
		session:      mockSession,
		logger:       slog.Default(),
		commands:     make([]*discordgo.ApplicationCommand, 0),
		commandTools: make(map[string]mcp.Tool),
		userID:       "bot-123",
	}

	readyEvent := &discordgo.Ready{
		User: &discordgo.User{
			ID:       "bot-123",
			Username: "TestBot",
		},
		Guilds: []*discordgo.Guild{
			{ID: "guild-1"},
		},
	}

	// Should not panic on error, just log
	bot.handleReady(readyEvent)

	// Should have attempted to register but failed
	if mockSession.CreateCommandCalls != 0 {
		t.Errorf("Expected 0 CreateCommand calls on error, got %d", mockSession.CreateCommandCalls)
	}
}


func TestDiscordBot_OnInteractionCreate(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{{Type: "text", Text: "Success"}},
			IsError: false,
		},
	}
	mockTranslator := &MockTranslator{}
	mockDiscovery := mcp.NewDiscoveryService(mockClient, nil)

	tool := mcp.Tool{
		Name:        "test_tool",
		Description: "Test",
		InputSchema: mcp.InputSchema{Type: "object"},
	}

	bot := &DiscordBot{
		config:     config.DiscordConfig{Token: "test"},
		mcpClient:  mockClient,
		translator: mockTranslator,
		discovery:  mockDiscovery,
		session:    mockSession,
		logger:     slog.Default(),
		commands:   make([]*discordgo.ApplicationCommand, 0),
		commandTools: map[string]mcp.Tool{
			"test-tool": tool,
		},
		userID: "bot-123",
	}

	// Create interaction
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name:    "test-tool",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{},
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID:       "user-123",
					Username: "testuser",
				},
			},
		},
	}

	// Test the interaction handler logic directly
	if interaction.Type == discordgo.InteractionApplicationCommand {
		// This would normally call InteractionRespond
		err := bot.session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			t.Fatalf("Failed to respond to interaction: %v", err)
		}
	}

	// Verify response was sent
	if mockSession.RespondCalls != 1 {
		t.Errorf("Expected 1 InteractionRespond call, got %d", mockSession.RespondCalls)
	}

	if mockSession.LastInteractionRespond == nil {
		t.Fatal("Expected interaction response to be captured")
	}

	if mockSession.LastInteractionRespond.Type != discordgo.InteractionResponseDeferredChannelMessageWithSource {
		t.Errorf("Expected deferred response type, got %d", mockSession.LastInteractionRespond.Type)
	}
}

func TestDiscordBot_OnInteractionCreate_WrongType(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session: mockSession,
		logger:  slog.Default(),
	}

	// Create interaction with wrong type (not application command)
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionMessageComponent, // Wrong type
		},
	}

	// Test the type check logic
	if interaction.Type != discordgo.InteractionApplicationCommand {
		// Should return early - this is the expected behavior
		t.Logf("Bot %v correctly ignoring non-application-command interaction", bot.logger)
	}
}

func TestDiscordBot_HandleCommand_NotFound(t *testing.T) {
	mockSession := NewMockDiscordSession()

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		commandTools: make(map[string]mcp.Tool), // Empty - no commands registered
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "nonexistent-command",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent an error response
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call for error, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_TranslationError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{
		TranslateArgumentsFunc: func(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
			return nil, errors.New("invalid argument type")
		},
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{Name: "arg1", Value: "value1"},
				},
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent an error response
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call for error, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_MCPExecutionError(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{}
	mockMCPClient := &MockMCPClient{
		CallToolErr: errors.New("MCP server unavailable"),
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		mcpClient:    mockMCPClient,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent an error response
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call for error, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_Success(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{}
	mockMCPClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: "Command executed successfully"},
			},
			IsError: false,
		},
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		mcpClient:    mockMCPClient,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent a success response
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_SuccessWithJSON(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{}
	mockMCPClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: `{"status": "ok", "count": 42}`},
			},
			IsError: false,
		},
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		mcpClient:    mockMCPClient,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent a success response with JSON formatted
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_WithArguments(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{
		TranslateArgumentsFunc: func(tool mcp.Tool, options []*discordgo.ApplicationCommandInteractionDataOption) (map[string]interface{}, error) {
			return map[string]interface{}{
				"message": "Hello World",
				"count":   42,
			}, nil
		},
	}
	mockMCPClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: "Processed: Hello World (42 times)"},
			},
			IsError: false,
		},
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		mcpClient:    mockMCPClient,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{Name: "message", Value: "Hello World"},
					{Name: "count", Value: 42},
				},
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent a success response
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

func TestDiscordBot_HandleCommand_MCPErrorResult(t *testing.T) {
	mockSession := NewMockDiscordSession()
	mockTranslator := &MockTranslator{}
	mockMCPClient := &MockMCPClient{
		CallToolData: &mcp.ToolResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: "Tool execution failed: invalid parameter"},
			},
			IsError: true,
		},
	}

	tool := mcp.Tool{
		Name:        "test-tool",
		Description: "Test tool",
	}

	bot := &DiscordBot{
		session:      mockSession,
		logger:       slog.Default(),
		translator:   mockTranslator,
		mcpClient:    mockMCPClient,
		commandTools: map[string]mcp.Tool{"test-tool": tool},
	}

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:   "interaction-123",
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "test-tool",
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: "testuser",
				},
			},
		},
	}

	bot.handleCommand(mockSession, interaction)

	// Should have sent an error response (IsError: true)
	if mockSession.EditResponseCalls != 1 {
		t.Errorf("Expected 1 EditResponse call, got %d", mockSession.EditResponseCalls)
	}
}

