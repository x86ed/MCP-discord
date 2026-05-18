## Context

This is the initial scaffolding for MCP-Discord, a Go application that acts as a 1-to-1 bridge between a Model Context Protocol (MCP) server and Discord. The project currently has only `go.mod` and basic repository files. We need to establish the foundational structure, configuration system, documentation, and testing infrastructure before implementing the core bot functionality.

The configuration system is inspired by Claude Desktop and VS Code's MCP integration patterns, where users specify how to launch and connect to MCP servers via JSON configuration files.

## Goals / Non-Goals

**Goals:**
- Establish a maintainable Go project structure following standard conventions
- Create a flexible JSON-based configuration system for MCP server connections
- Set up comprehensive documentation for users and developers
- Implement testing infrastructure supporting unit and integration tests
- Make the project easy to understand, extend, and deploy

**Non-Goals:**
- Implementing the actual bot logic (that comes in later changes)
- Supporting multiple MCP servers simultaneously (1-to-1 wrapper only)
- Creating a web UI for configuration management
- Building deployment infrastructure (Docker/K8s configs can come later)

## Decisions

### Directory Structure: Standard Go Project Layout

**Decision:** Use the standard Go project layout with `cmd/`, `internal/`, and `pkg/` directories.

**Rationale:** 
- This is the widely accepted convention in the Go community
- Clear separation between application entry points (`cmd/`), private code (`internal/`), and reusable libraries (`pkg/`)
- Makes the project immediately familiar to Go developers
- Tools and IDEs work well with this structure

**Structure:**
```
cmd/mcpdiscord/      # Main application
internal/            # Private application code
  bot/              # Discord bot logic
  mcp/              # MCP client implementation
  translator/       # Schema translation between MCP and Discord
  config/           # Configuration parsing and validation
pkg/                # Public libraries (if any emerge)
examples/           # Example configurations
docs/               # Documentation
```

**Alternatives Considered:**
- Flat structure: Rejected - doesn't scale, hard to navigate
- Feature-based folders: Rejected - Go convention is package-based

### Configuration Format: JSON with Environment Variable Support

**Decision:** Use JSON configuration files with support for environment variable interpolation (`${VAR_NAME}` syntax).

**Rationale:**
- JSON is what Claude Desktop and VS Code use for MCP configuration - familiar to users
- Easy to parse with Go's standard library
- Supports structured data (objects, arrays) naturally
- Environment variable interpolation keeps secrets out of config files
- Schema validation is straightforward with JSON

**Example Configuration:**
```json
{
  "discord": {
    "token": "${DISCORD_BOT_TOKEN}",
    "guildId": "optional-for-dev"
  },
  "mcp": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-weather"],
    "env": {
      "API_KEY": "${WEATHER_API_KEY}"
    },
    "transport": "stdio"
  }
}
```

**Alternatives Considered:**
- YAML: Rejected - more complex parsing, no significant benefit
- TOML: Rejected - less familiar, nested structures more verbose
- Environment variables only: Rejected - doesn't support complex MCP server configurations

### Configuration Loading Strategy

**Decision:** Support multiple configuration sources with precedence: CLI flag > Environment variable > Default path.

**Rationale:**
- Flexibility for different deployment scenarios
- CLI flag: Explicit control for development/testing
- Environment variable: Container/cloud deployments
- Default path: Convenience for local development

**Implementation:**
```go
// Priority order:
1. --config flag
2. MCP_CONFIG_PATH environment variable
3. ./mcp-config.json (default)
```

### Testing Approach: Go Standard Testing with Build Tags

**Decision:** Use Go's standard `testing` package with build tags to separate unit and integration tests.

**Rationale:**
- No external test framework dependencies for basic tests
- Table-driven tests are idiomatic in Go
- Build tags (`// +build integration`) allow selective test execution
- Integration tests can be skipped in CI if MCP servers aren't available

**Test Organization:**
- Unit tests: Same package, `_test.go` suffix, no build tags
- Integration tests: Build tag `integration`, requires external MCP server
- Mock interfaces for all external dependencies (Discord API, MCP client)

**Alternatives Considered:**
- Testify/suite: Rejected for unit tests (unnecessary), may add later for assertions
- Separate test directories: Rejected - goes against Go conventions

### Code Coverage Enforcement: 95% Threshold in CI

**Decision:** Enforce a minimum 95% code coverage threshold in GitHub Actions that blocks PRs if not met.

**Rationale:**
- High coverage ensures comprehensive testing of critical path logic
- Automated enforcement prevents coverage regression
- Failing builds on low coverage makes quality gates explicit
- 95% is achievable with proper unit tests and allows for reasonable edge cases

**Implementation:**
- Use `go test -coverprofile=coverage.out ./...` to generate coverage report
- Parse coverage output to calculate overall percentage
- Fail CI workflow if coverage < 95%
- Display per-package coverage breakdown in workflow logs

**Alternatives Considered:**
- 80% threshold: Rejected - too lenient for a new project, better to start strict
- No enforcement: Rejected - coverage tends to decline without automated checks
- Coverage per-package instead of overall: Deferred - start with overall, can add later
- External coverage service (Codecov): Deferred - stdlib tools sufficient initially

### Documentation Structure: README + docs/ Directory

**Decision:** Comprehensive README with detailed docs in separate `docs/` directory.

**Rationale:**
- README provides quick start and overview
- Separate docs for detailed topics (architecture, deployment, contribution)
- Markdown for easy maintenance and GitHub rendering
- ASCII diagrams for architecture visualization (no external dependencies)

**Documentation Files:**
```
README.md              # Quick start, overview, basic config
docs/
  architecture.md      # System design, component diagram, data flow
  configuration.md     # Detailed config reference
  deployment.md        # Production deployment guide
  development.md       # Development setup, testing, contributing
```

### Makefile for Build Tasks

**Decision:** Provide a Makefile with common development tasks.

**Rationale:**
- Standard in Go projects for convenience
- Documents common commands
- Makes onboarding easier for new developers

**Tasks:**
```makefile
build             # Compile binary
test              # Run unit tests
test-integration  # Run integration tests
coverage          # Generate and display coverage report
run               # Run locally with example config
clean             # Remove build artifacts
lint              # Run linters (golangci-lint)
```

## Risks / Trade-offs

**Risk: Configuration schema changes break existing configs**
- Mitigation: Version the config schema, implement migration helpers if schema evolves

**Risk: Go standard testing may not be sufficient for complex scenarios**
- Mitigation: Start simple, add testify or other libraries if needed later

**Trade-off: Standard layout adds directories for small project**
- Accepted: Better to start with good structure than refactor later

**Risk: Documentation becomes stale**
- Mitigation: Keep docs close to code, review in PRs, use OpenSpec to track changes

**Trade-off: Environment variable interpolation adds parsing complexity**
- Accepted: Significant UX benefit for secrets management outweighs implementation complexity

**Trade-off: 95% coverage threshold is strict and may slow down development**
- Accepted: High quality bar from the start prevents technical debt; threshold can be adjusted if it becomes blocker
- Mitigation: Good test utilities and examples make achieving high coverage easier
