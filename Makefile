# MCP-Discord Makefile
# Build automation for the MCP-Discord bot

# Binary name
BINARY_NAME=mcpdiscord
BINARY_PATH=./cmd/mcpdiscord

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags="-w -s -X main.version=$(VERSION)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build test test-integration coverage lint clean run install help

# Default target
all: test build

# Build the binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "$(GREEN)Build complete: $(BINARY_NAME)$(NC)"

# Build for multiple platforms
build-all:
	@echo "$(GREEN)Building for multiple platforms...$(NC)"
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(BINARY_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 $(BINARY_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(BINARY_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)
	@echo "$(GREEN)Multi-platform build complete$(NC)"

# Run tests
test:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GOTEST) -v -race -short ./...

# Run integration tests
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GOTEST) -v -race -tags=integration ./...

# Run all tests (unit + integration)
test-all: test test-integration

# Generate coverage report
coverage:
	@echo "$(GREEN)Generating coverage report...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Check coverage threshold (95%)
coverage-check: coverage
	@echo "$(YELLOW)Checking coverage threshold...$(NC)"
	@TOTAL=$$($(GOCMD) tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	THRESHOLD=95; \
	echo "Total coverage: $$TOTAL%"; \
	if [ "$$(echo "$$TOTAL < $$THRESHOLD" | bc)" -eq 1 ]; then \
		echo "$(RED)Coverage $$TOTAL% is below threshold $$THRESHOLD%$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)Coverage $$TOTAL% meets threshold $$THRESHOLD%$(NC)"; \
	fi

# Run linter
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)" && exit 1)
	golangci-lint run

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) ./...

# Run the application
run: build
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME) --config mcp-config.json

# Run with custom config
run-config:
	@echo "$(GREEN)Running $(BINARY_NAME) with custom config...$(NC)"
	./$(BINARY_NAME) --config $(CONFIG)

# Install binary to system
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"

# Uninstall binary from system
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstalled$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

# Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOGET) -v ./...
	$(GOMOD) download
	$(GOMOD) tidy

# Verify dependencies
deps-verify:
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	$(GOMOD) verify

# Update dependencies
deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Run all checks (format, lint, test, coverage)
check: fmt lint test coverage-check
	@echo "$(GREEN)All checks passed!$(NC)"

# Quick check without coverage
check-quick: fmt lint test
	@echo "$(GREEN)Quick checks passed!$(NC)"

# Development mode (build + run with file watching)
dev: build
	@echo "$(GREEN)Running in development mode...$(NC)"
	./$(BINARY_NAME) --config mcp-config.json

# Docker build
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) -f examples/deployment/Dockerfile .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "$(GREEN)Docker image built: $(BINARY_NAME):$(VERSION)$(NC)"

# Docker run
docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker-compose -f examples/deployment/docker-compose.yml up

# Docker stop
docker-stop:
	@echo "$(YELLOW)Stopping Docker container...$(NC)"
	docker-compose -f examples/deployment/docker-compose.yml down

# Show help
help:
	@echo "$(GREEN)MCP-Discord Bot - Available Commands:$(NC)"
	@echo ""
	@echo "  $(YELLOW)Building:$(NC)"
	@echo "    make build         - Build the binary"
	@echo "    make build-all     - Build for multiple platforms"
	@echo "    make install       - Install binary to /usr/local/bin"
	@echo "    make uninstall     - Remove installed binary"
	@echo ""
	@echo "  $(YELLOW)Testing:$(NC)"
	@echo "    make test          - Run unit tests"
	@echo "    make test-integration - Run integration tests"
	@echo "    make test-all      - Run all tests"
	@echo "    make coverage      - Generate coverage report"
	@echo "    make coverage-check - Check coverage meets 95% threshold"
	@echo ""
	@echo "  $(YELLOW)Code Quality:$(NC)"
	@echo "    make lint          - Run golangci-lint"
	@echo "    make fmt           - Format code with gofmt"
	@echo "    make check         - Run all checks (fmt, lint, test, coverage)"
	@echo "    make check-quick   - Run quick checks without coverage"
	@echo ""
	@echo "  $(YELLOW)Running:$(NC)"
	@echo "    make run           - Build and run the bot"
	@echo "    make run-config CONFIG=path - Run with custom config"
	@echo "    make dev           - Run in development mode"
	@echo ""
	@echo "  $(YELLOW)Dependencies:$(NC)"
	@echo "    make deps          - Download dependencies"
	@echo "    make deps-verify   - Verify dependencies"
	@echo "    make deps-update   - Update all dependencies"
	@echo ""
	@echo "  $(YELLOW)Docker:$(NC)"
	@echo "    make docker-build  - Build Docker image"
	@echo "    make docker-run    - Run with docker-compose"
	@echo "    make docker-stop   - Stop docker-compose"
	@echo ""
	@echo "  $(YELLOW)Maintenance:$(NC)"
	@echo "    make clean         - Remove build artifacts"
	@echo "    make help          - Show this help message"
	@echo ""
