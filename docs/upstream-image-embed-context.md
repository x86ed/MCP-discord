# Upstream Context: Optional Image Embed in Discord Responses

## Summary

This project added optional image support for Discord embed responses so a successful tool result can include a visual alongside text.

The change is intentionally minimal:

- Added one optional field on the bot embed model: `ImageURL`.
- Mapped `ImageURL` to Discord's embed image field at send time.
- Added integration coverage to verify the webhook payload includes the image URL when provided.

## Why This Was Needed

Some tool outputs are easier to consume when a screenshot, chart, or generated image is shown directly in the Discord response. Text-only embeds were already supported, but there was no way to pass an image URL through the response type.

## Implementation Details

### 1. Response Model Extension

File: `internal/bot/bot.go`

- Extended `Embed` with:
  - `ImageURL string`

Behavior:

- Empty `ImageURL`: no image is sent (existing behavior preserved).
- Non-empty `ImageURL`: image is included in the embed.

### 2. Discord Mapping

File: `internal/bot/session.go`

In `sendResponse(...)`, embed construction now includes:

- `embed.Image = &discordgo.MessageEmbedImage{URL: resp.Embed.ImageURL}` when `ImageURL != ""`.

This keeps image handling localized to transport formatting and avoids changing command execution flow.

### 3. Test Coverage

File: `internal/bot/integration_test.go`

Added:

- `TestDiscordBot_SendResponse_WithImage`

This test verifies:

- `InteractionResponseEdit` receives an embed.
- Embed image is present.
- Image URL matches the value supplied in `Response.Embed.ImageURL`.

## Compatibility and Risk

- Backward compatible: existing responses without image continue to behave exactly the same.
- No protocol breakage: this is an additive field on an internal response struct.
- No change to command registration, argument translation, or MCP client APIs.

## Usage Pattern

Set the image URL when building the response embed:

```go
return &Response{
    Embed: &Embed{
        Title:       "✓ my_tool",
        Description: "Operation complete",
        Color:       0x00FF00,
        ImageURL:    "https://example.com/output.png",
    },
}
```

## Notes for Upstream Discussion

- This implementation assumes image URLs are externally reachable by Discord.
- If upstream wants first-class MCP image block support, a follow-up can parse non-text content blocks and map supported image MIME/content forms into either:
  - Hosted URL embeds, or
  - Attachment uploads (larger change, requires file upload path handling).

## Validation Run

Targeted tests passed for response behavior after this change, including the new image case.
