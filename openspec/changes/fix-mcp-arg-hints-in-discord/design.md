## Context

The MCP-Discord bot translates MCP tool definitions into Discord slash commands. The current implementation reads parameter descriptions from the MCP tool schema's `inputSchema.properties[name].description` field and passes them to Discord, but there's no visibility when descriptions are missing or empty. This makes it difficult to diagnose whether parameter hints are missing due to:
- MCP server not providing descriptions in the tool schema
- Discord bot not correctly reading or passing descriptions

The root cause could be either server-side (MCP not exposing descriptions) or client-side (bot not handling descriptions correctly).

## Goals / Non-Goals

**Goals:**
- Make description propagation explicit and debuggable
- Add logging to identify when MCP servers don't provide parameter descriptions
- Ensure users see helpful hints for all parameters (type hints for complex types, descriptive text for simple types)
- Provide clear diagnostic information to help users fix their MCP server configurations

**Non-Goals:**
- Modify the MCP protocol or server implementations
- Create a UI for editing descriptions in Discord
- Support custom description overrides in bot configuration
- Change how type hints work for arrays/objects (existing "(Comma-separated list)" format is sufficient)

## Decisions

### Decision 1: Add warning logs for missing descriptions

**Rationale:** The current implementation silently falls back to "No description" when descriptions are missing. Adding logging will help users identify when their MCP server isn't providing proper parameter metadata, making it easier to fix the root cause rather than working around it in the bot.

**Alternative considered:** Show a notice in Discord when commands have missing descriptions. Rejected because it would clutter the user experience and the right place to fix this is in the MCP server configuration.

### Decision 2: Keep existing fallback behavior

**Rationale:** The current fallback to "No description" and type hints like "(Comma-separated list)" is functional. We're making the behavior explicit in the spec rather than changing it, ensuring the implementation matches expected behavior.

**Alternative considered:** Use parameter name as description when missing. Rejected because it doesn't provide additional value and "No description" is clearer about the issue.

### Decision 3: Treat empty string descriptions as missing

**Rationale:** An empty string provides no information to users, so it should be treated the same as a missing description field. This makes the behavior consistent and prevents confusing empty descriptions in Discord.

**Alternative considered:** Distinguish between null/undefined and empty string. Rejected because the practical difference for users is zero.

## Risks / Trade-offs

**Risk:** Increased log volume if many MCP servers have missing descriptions  
→ **Mitigation:** Use warning level (not error), log once per tool registration rather than per command invocation

**Risk:** Users may think the bot is broken when seeing "No description"  
→ **Mitigation:** The warning logs will help developers/operators identify and fix MCP server configuration issues

**Trade-off:** Adding logging adds slight complexity  
→ **Accepted:** The diagnostic value outweighs the minimal code complexity
