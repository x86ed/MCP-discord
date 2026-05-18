# Development Guide

Guide for contributing to and developing the MCP-Discord bot.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Standards](#code-standards)
- [Debugging](#debugging)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

- **Go 1.21 or higher** ([Install Go](https://go.dev/doc/install))
- **Git** for version control
- **Discord account** and bot token for testing
- **Node.js** (optional, for testing with npm-based MCP servers)
- **Make** (optional, for using Makefile commands)

### Fork and Clone

```bash
# Fork on GitHub first, then:
git clone https://github.com/YOUR_USERNAME/MCP-discord.git
cd MCP-discord

# Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/MCP-discord.git
```

## Development Setup

### 1. Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. Create Development Config

```bash
# Copy example config
cp examples/with-env/mcp-config.json mcp-config.json

# Create .env file
cat > .env << 'EOF'
DISCORD_BOT_TOKEN=your_development_bot_token
DISCORD_GUILD_ID=your_test_server_id
EOF

# Load environment variables (with direnv)
echo "dotenv" > .envrc
direnv allow

# Or manually export
set -a; source .env; set +a
```

### 3. Create Test Discord Server

1. Create a new Discord server for testing
2. Enable Developer Mode (User Settings > Advanced > Developer Mode)
3. Right-click server name → Copy Server ID
4. Use this ID as `DISCORD_GUILD_ID` for instant command updates

### 4. Verify Setup

```bash
# Run tests
go test ./...

# Build application
go build ./cmd/mcpdiscord

# Run with development config
./mcpdiscord --config mcp-config.json
```

## Project Structure

```
MCP-discord/
├── cmd/
│   └── mcpdiscord/          # Application entry point
│       └── main.go
├── internal/                 # Private application code
│   ├── bot/                 # Discord bot implementation
│   │   └── bot.go
│   ├── config/              # Configuration management
│   │   ├── types.go         # Config structs
│   │   ├── loader.go        # File loading and env interpolation
│   │   ├── validator.go     # Validation logic
│   │   └── *_test.go        # Unit tests
│   ├── mcp/                 # MCP client
│   │   └── client.go
│   ├── translator/          # MCP ↔ Discord translation
│   │   └── translator.go
│   └── testutil/            # Test helpers
│       ├── config.go
│       ├── mockmcp.go
│       └── mockdiscord.go
├── examples/                 # Example configurations
│   ├── basic/
│   ├── with-env/
│   └── deployment/
├── docs/                     # Documentation
├── .github/
│   └── workflows/           # CI/CD workflows
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build automation
└── README.md

Package Responsibilities:
- `cmd/mcpdiscord`: Entry point, CLI, initialization
- `internal/config`: Configuration loading and validation
- `internal/bot`: Discord API interaction
- `internal/mcp`: MCP protocol client
- `internal/translator`: Schema and data translation
- `internal/testutil`: Shared test utilities
```

## Development Workflow

### Standard Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make Changes**
   - Write code following [Code Standards](#code-standards)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Locally**
   ```bash
   # Run tests
   make test

   # Check coverage
   make coverage

   # Run linter
   make lint

   # Test integration
   make test-integration
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "Add feature: description of feature"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/my-new-feature
   # Create pull request on GitHub
   ```

### Commit Message Convention

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `chore`: Build process or tooling changes

**Examples:**
```
feat(translator): add support for nested object schemas

Implements recursive translation of nested JSON Schema objects
to Discord command options using JSON string fallback.

Closes #42
```

```
fix(config): handle empty string in environment variable interpolation

Previously, ${VAR} with empty string would be treated as unresolved.
Now correctly resolves to empty string.
```

## Testing

### Running Tests

```bash
# All unit tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/config/...

# Verbose output
go test -v ./...

# Integration tests only
go test -tags=integration ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using Makefile

```bash
make test              # Run unit tests
make test-integration  # Run integration tests
make coverage          # Generate coverage report
```

### Writing Tests

**Unit Test Example:**
```go
// internal/config/loader_test.go
func TestLoad_ValidConfig(t *testing.T) {
    // Arrange
    cfg := testutil.TestConfig(t)
    path := testutil.WriteTestConfig(t, cfg)

    // Act
    result, err := Load(path)

    // Assert
    if err != nil {
        t.Fatalf("Load() failed: %v", err)
    }
    if result.Discord.Token != cfg.Discord.Token {
        t.Errorf("Token mismatch: got %q, want %q",
            result.Discord.Token, cfg.Discord.Token)
    }
}
```

**Integration Test Example:**
```go
// +build integration

// internal/config/integration_test.go
func TestIntegration_FullFlow(t *testing.T) {
    // Set up environment
    os.Setenv("TEST_TOKEN", "test-value")
    defer os.Unsetenv("TEST_TOKEN")

    // Test full configuration flow
    cfg, err := LoadFromSource("", "", "testdata/config.json")
    if err != nil {
        t.Fatalf("LoadFromSource() failed: %v", err)
    }

    // Validate
    if err := cfg.Validate(); err != nil {
        t.Fatalf("Validate() failed: %v", err)
    }
}
```

### Test Coverage Requirements

- **Minimum**: 95% coverage (enforced by CI)
- **Goal**: 100% for critical packages (config, translator)
- **Focus**: Edge cases, error paths, validation logic

### Using Test Utilities

```go
import "github.com/yourusername/MCP-discord/internal/testutil"

// Create test config
cfg := testutil.TestConfig(t)

// Create test config with custom env
cfg := testutil.TestConfigWithEnv(t, map[string]string{
    "API_KEY": "test-key",
})

// Write config to temp file
path := testutil.WriteTestConfig(t, cfg)

// Create mock MCP server
mockMCP := testutil.NewMockMCPServer(t)
mockMCP.AddSimpleTool("test_tool", "Test tool description")

// Create mock Discord interaction
mockInteraction := testutil.NewMockDiscordInteraction(t, "test_command")
mockInteraction.AddOption("param1", "value1")
```

## Code Standards

### Go Code Style

Follow standard Go conventions:
- Use `gofmt` for formatting
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `golangci-lint` for linting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Or use Makefile
make lint
```

### Documentation

**Package Documentation:**
```go
// Package config provides configuration loading and validation
// for the MCP-Discord bot.
//
// Configuration can be loaded from JSON files with support for
// environment variable interpolation using ${VAR_NAME} syntax.
package config
```

**Function Documentation:**
```go
// Load reads and parses a configuration file from the given path.
// It returns an error if the file doesn't exist, contains invalid
// JSON, or fails validation.
//
// Environment variables in the format ${VAR_NAME} are automatically
// interpolated with their values from the environment.
func Load(path string) (*Config, error) {
    // ...
}
```

**Exported vs Unexported:**
- Export only what needs to be public
- Document all exported types and functions
- Keep implementation details private

### Error Handling

**Descriptive Errors:**
```go
if cfg.Discord.Token == "" {
    return fmt.Errorf("discord.token is required but was empty")
}
```

**Error Wrapping:**
```go
data, err := os.ReadFile(path)
if err != nil {
    return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
}
```

**Multiple Errors:**
```go
var errors []string
if cfg.Discord.Token == "" {
    errors = append(errors, "discord.token is required")
}
if cfg.MCP.Command == "" {
    errors = append(errors, "mcp.command is required")
}
if len(errors) > 0 {
    return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
}
```

### Logging

Use structured logging with `slog`:

```go
import "log/slog"

// Good
slog.Info("loaded configuration",
    "source", configPath,
    "transport", cfg.MCP.Transport)

// Bad (avoid unstructured logging)
log.Printf("Loaded config from %s with transport %s", configPath, cfg.MCP.Transport)
```

**Log Levels:**
- `Debug`: Detailed information for troubleshooting
- `Info`: General informational messages
- `Warn`: Warning messages (recoverable issues)
- `Error`: Error messages (failures)

## Debugging

### Local Debugging

**With VSCode:**

Create `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug MCP-Discord",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/mcpdiscord",
            "args": ["--config", "mcp-config.json"],
            "env": {
                "DISCORD_BOT_TOKEN": "your-token"
            }
        }
    ]
}
```

**With Delve:**
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/mcpdiscord -- --config mcp-config.json
```

### Debug Logging

Enable debug logging:
```go
// In main.go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug, // Set to debug
}))
slog.SetDefault(logger)
```

### Testing with Real Discord

1. Create test Discord server
2. Set `guildId` in config for instant command updates
3. Use different bot token than production
4. Invite bot with proper scopes

### Common Issues

**Commands not updating:**
- Use `guildId` for instant updates during development
- Global commands take up to 1 hour

**MCP server not starting:**
- Run MCP command manually to test
- Check `command` and `args` in config
- View MCP server logs

**Environment variables not working:**
- Export before running
- Check for typos (case sensitive)
- Use `printenv` to verify

## Contributing

### Pull Request Process

1. **Fork the repository**

2. **Create feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

3. **Make changes**
   - Follow code standards
   - Add tests (maintain 95%+ coverage)
   - Update documentation
   - Run linter and tests locally

4. **Commit with descriptive messages**
   ```bash
   git commit -m "feat(component): description"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/my-feature
   ```

6. **Create Pull Request**
   - Fill out PR template
   - Link related issues
   - Request review

7. **Address review feedback**
   - Make requested changes
   - Push updates to same branch
   - PR will update automatically

### Code Review Guidelines

**For Authors:**
- Keep PRs focused and small
- Include tests and documentation
- Respond to feedback promptly
- Be open to suggestions

**For Reviewers:**
- Be constructive and respectful
- Test changes locally if possible
- Check code coverage
- Verify documentation is updated

### Getting Help

- 🐛 [Report a Bug](https://github.com/yourusername/MCP-discord/issues/new?labels=bug)
- 💡 [Request a Feature](https://github.com/yourusername/MCP-discord/issues/new?labels=enhancement)
- 💬 [Discussions](https://github.com/yourusername/MCP-discord/discussions)
- 📖 [Documentation](../README.md)

## Development Tools

### Useful Commands

```bash
# Quick build and run
make run

# Clean build artifacts
make clean

# View coverage in browser
make coverage

# Run all checks (tests + lint)
make check

# Install as system command
make install
```

### Recommended VSCode Extensions

- Go (golang.go)
- Go Test Explorer (ethan-reesor.go-test-explorer)
- Error Lens (usernamehw.errorlens)
- GitLens (eamodio.gitlens)

### Git Hooks

Create `.git/hooks/pre-commit`:
```bash
#!/bin/bash
# Run tests and linter before commit

echo "Running tests..."
make test || exit 1

echo "Running linter..."
make lint || exit 1

echo "Pre-commit checks passed!"
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

Happy coding! 🚀
