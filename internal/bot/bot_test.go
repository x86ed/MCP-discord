package bot

import (
	"strings"
	"testing"

	"mcpdiscord/internal/mcp"
)

func TestResponseTruncation(t *testing.T) {
	// Test that formatResult properly truncates long content
	longContent := strings.Repeat("a", 5000)

	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: longContent,
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if len(resp.Embed.Description) > 4096 {
		t.Errorf("Description length %d exceeds Discord limit of 4096", len(resp.Embed.Description))
	}

	if !strings.Contains(resp.Embed.Description, "truncated") {
		t.Error("Expected truncation marker in description")
	}
}

func TestFormatResultSuccess(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: "Success message",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.Color != 0x00FF00 {
		t.Errorf("Expected green color (0x00FF00), got 0x%06X", resp.Embed.Color)
	}

	if !strings.Contains(resp.Embed.Title, "✓") {
		t.Error("Expected success checkmark in title")
	}

	if !strings.Contains(resp.Embed.Description, "Success message") {
		t.Error("Expected result content in description")
	}
}

func TestFormatResultError(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: "Error message",
			},
		},
		IsError: true,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.Color != 0xFF0000 {
		t.Errorf("Expected red color (0xFF0000), got 0x%06X", resp.Embed.Color)
	}

	if !strings.Contains(resp.Embed.Title, "✗") {
		t.Error("Expected error X in title")
	}
}

func TestFormatResultJSONDetection(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: `{"key": "value"}`,
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	// Should wrap JSON in code blocks
	if !strings.Contains(resp.Embed.Description, "```json") {
		t.Error("Expected JSON code block formatting")
	}
}

func TestFormatResultMultipleBlocks(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: "Block 1",
			},
			{
				Type: "text",
				Text: "Block 2",
			},
			{
				Type: "text",
				Text: "Block 3",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	// Should combine all blocks
	if !strings.Contains(resp.Embed.Description, "Block 1") ||
		!strings.Contains(resp.Embed.Description, "Block 2") ||
		!strings.Contains(resp.Embed.Description, "Block 3") {
		t.Error("Expected all content blocks in description")
	}
}

func TestFormatResultMarkdownImageToEmbedImage(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: "Here is the chart:\n\n![chart](https://example.com/chart.png)",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.ImageURL != "https://example.com/chart.png" {
		t.Fatalf("Expected image URL to be extracted, got %q", resp.Embed.ImageURL)
	}

	if strings.Contains(resp.Embed.Description, "![chart]") {
		t.Fatalf("Expected markdown image syntax to be removed from description, got %q", resp.Embed.Description)
	}
}

func TestFormatResultMarkdownImageJSONUnchanged(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "text",
				Text: `{"note":"![chart](https://example.com/chart.png)"}`,
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.ImageURL != "" {
		t.Fatalf("Expected JSON payload to skip markdown extraction, got %q", resp.Embed.ImageURL)
	}

	if !strings.Contains(resp.Embed.Description, "```json") {
		t.Fatalf("Expected JSON formatting to be preserved, got %q", resp.Embed.Description)
	}
}

func TestFormatResultResourceLinkImageToEmbedImage(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "resource",
				ResourceLink: &mcp.ResourceLink{
					URL:      "https://example.com/from-resource-link.png",
					MimeType: "image/png",
				},
			},
			{
				Type: "text",
				Text: "Image provided by resourceLink",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.ImageURL != "https://example.com/from-resource-link.png" {
		t.Fatalf("Expected image URL from resourceLink, got %q", resp.Embed.ImageURL)
	}
}

func TestFormatResultResorceLinkImageToEmbedImage(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type: "resource",
				ResorceLink: &mcp.ResourceLink{
					URI:      "https://example.com/from-resorce-link.jpg",
					MimeType: "image/jpeg",
				},
			},
			{
				Type: "text",
				Text: "Image provided by resorceLink",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.ImageURL != "https://example.com/from-resorce-link.jpg" {
		t.Fatalf("Expected image URL from resorceLink, got %q", resp.Embed.ImageURL)
	}
}

func TestFormatResultInlineResourceLinkImageToEmbedImage(t *testing.T) {
	result := &mcp.ToolResult{
		Content: []mcp.ContentBlock{
			{
				Type:     "resource_link",
				MimeType: "image/webp",
				URI:      "https://skydex.info/img/airports/KLAX.webp",
				Name:     "LAX",
			},
			{
				Type: "text",
				Text: "Airport image",
			},
		},
		IsError: false,
	}

	bot := &DiscordBot{}
	resp := bot.formatResult("test_tool", result)

	if resp.Embed == nil {
		t.Fatal("Expected embed to be created")
	}

	if resp.Embed.ImageURL != "https://skydex.info/img/airports/KLAX.webp" {
		t.Fatalf("Expected image URL from inline resource_link, got %q", resp.Embed.ImageURL)
	}
}
