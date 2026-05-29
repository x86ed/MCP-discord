## Why

MCP tool parameter descriptions are not being propagated to Discord slash command options, resulting in users seeing command arguments with no hints or context about what values to provide. This creates poor UX as users cannot discover what each parameter does without external documentation.

## What Changes

- Propagate MCP tool parameter descriptions to Discord slash command option descriptions
- Ensure parameter hints include type information for complex types (arrays, objects)
- Handle missing or empty parameter descriptions with sensible defaults

## Capabilities

### New Capabilities
<!-- None - this is a fix to existing capability -->

### Modified Capabilities

- `slash-command-registration`: Add requirement to map MCP parameter descriptions to Discord option descriptions when registering commands

## Impact

- Affected code: `internal/translator/translator.go` (parameter to option mapping)
- Affected code: Discord command registration logic
- User-visible: Discord slash commands will now show helpful descriptions for each parameter
- No breaking changes - purely additive improvement to command registration
