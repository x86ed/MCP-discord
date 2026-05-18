## 1. Project Directory Structure

- [x] 1.1 Create `cmd/mcpdiscord/` directory for main application entry point
- [x] 1.2 Create `internal/bot/` directory for Discord bot logic
- [x] 1.3 Create `internal/mcp/` directory for MCP client implementation
- [x] 1.4 Create `internal/translator/` directory for schema translation
- [x] 1.5 Create `internal/config/` directory for configuration handling
- [x] 1.6 Create `pkg/` directory for public packages (if needed)
- [x] 1.7 Create `examples/` directory for example configurations
- [x] 1.8 Create `docs/` directory for detailed documentation

## 2. Configuration System

- [x] 2.1 Define configuration struct types in `internal/config/types.go`
- [x] 2.2 Implement JSON parsing with environment variable interpolation in `internal/config/loader.go`
- [x] 2.3 Add configuration validation logic in `internal/config/validator.go`
- [x] 2.4 Implement multi-source config loading (flag > env > default) in `internal/config/loader.go`
- [x] 2.5 Add descriptive error messages for missing or invalid configuration
- [x] 2.6 Write unit tests for configuration parsing in `internal/config/loader_test.go`
- [x] 2.7 Write unit tests for environment variable interpolation
- [x] 2.8 Write unit tests for configuration validation

## 3. Example Configurations

- [x] 3.1 Create `examples/basic/mcp-config.json` with weather server example
- [x] 3.2 Create `examples/basic/README.md` explaining the basic example
- [x] 3.3 Create `examples/with-env/mcp-config.json` demonstrating environment variable usage
- [x] 3.4 Create `examples/with-env/.env.example` showing required environment variables
- [x] 3.5 Add comments/documentation to all example config files

## 4. Main Application Entry Point

- [x] 4.1 Create `cmd/mcpdiscord/main.go` with basic structure
- [x] 4.2 Add CLI flag parsing for config file path
- [x] 4.3 Add environment variable support for config path
- [x] 4.4 Implement configuration loading on startup
- [x] 4.5 Add structured logging initialization
- [x] 4.6 Add graceful shutdown signal handling
- [x] 4.7 Add version information flag (`--version`)

## 5. Testing Infrastructure

- [x] 5.1 Create test utility helpers in `internal/testutil/` package
- [x] 5.2 Add configuration test helpers for creating test configs
- [x] 5.3 Add mock MCP server implementation for testing in `internal/testutil/mockmcp.go`
- [x] 5.4 Add Discord interaction mock helpers in `internal/testutil/mockdiscord.go`
- [x] 5.5 Create example integration test with build tag in `internal/config/integration_test.go`
- [x] 5.6 Set up GitHub Actions workflow for CI in `.github/workflows/test.yml`
- [x] 5.7 Configure CI to run unit tests on multiple Go versions
- [x] 5.8 Configure CI to generate coverage report with `go test -coverprofile`
- [x] 5.9 Add coverage percentage calculation to CI workflow
- [x] 5.10 Configure CI to enforce 95% minimum coverage threshold
- [x] 5.11 Configure CI to fail the build if coverage is below 95%
- [x] 5.12 Add per-package coverage breakdown reporting in CI output

## 6. Documentation

- [x] 6.1 Update `README.md` with project overview and quick start guide
- [x] 6.2 Add installation instructions to README
- [x] 6.3 Add basic configuration example to README
- [x] 6.4 Add usage examples to README
- [x] 6.5 Create `docs/architecture.md` with component diagrams
- [x] 6.6 Create `docs/configuration.md` with detailed config reference
- [x] 6.7 Create `docs/deployment.md` with production deployment guide
- [x] 6.8 Create `docs/development.md` with development setup and contributing guide
- [x] 6.9 Add systemd service file example in `examples/deployment/mcpdiscord.service`
- [x] 6.10 Add Docker example configuration in `examples/deployment/`

## 7. Build Tooling

- [x] 7.1 Create `Makefile` with build task
- [x] 7.2 Add test task to Makefile (unit tests only)
- [x] 7.3 Add test-integration task to Makefile (with build tag)
- [x] 7.4 Add run task to Makefile (runs with example config)
- [x] 7.5 Add clean task to Makefile
- [x] 7.6 Add lint task to Makefile (golangci-lint)
- [x] 7.7 Add coverage task to Makefile
- [x] 7.8 Update `.gitignore` with Go build artifacts and binary names

## 8. Package Stubs

- [x] 8.1 Create package stub file `internal/bot/bot.go` with package documentation
- [x] 8.2 Create package stub file `internal/mcp/client.go` with package documentation
- [x] 8.3 Create package stub file `internal/translator/translator.go` with package documentation
- [x] 8.4 Add interface definitions for key components (bot, mcp client, translator)
- [x] 8.5 Add package-level documentation comments explaining each package's purpose

## 9. Dependencies

- [x] 9.1 Add discordgo dependency to `go.mod`
- [x] 9.2 Add any JSON parsing utilities if needed (beyond stdlib)
- [x] 9.3 Add structured logging library (e.g., slog from stdlib or zerolog)
- [x] 9.4 Add golangci-lint configuration file `.golangci.yml`
- [x] 9.5 Run `go mod tidy` to clean up dependencies

## 10. Verification

- [x] 10.1 Verify project builds successfully with `go build ./...`
- [x] 10.2 Verify all tests pass with `go test ./...`
- [x] 10.3 Verify linter passes with `make lint`
- [x] 10.4 Verify documentation is complete and renders correctly
- [x] 10.5 Verify example configurations are valid and well-documented
- [x] 10.6 Run the application with example config to verify startup and config loading
