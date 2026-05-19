## ADDED Requirements

### Requirement: Test files alongside source code
The project SHALL place test files in the same package as the code being tested following Go conventions.

#### Scenario: Test file naming
- **WHEN** creating tests for a source file
- **THEN** test file SHALL be named `<filename>_test.go`
- **THEN** test file SHALL be in the same directory as source file

#### Scenario: Test package naming
- **WHEN** writing tests
- **THEN** tests MAY use same package name for white-box testing
- **THEN** tests MAY use `<package>_test` suffix for black-box testing

### Requirement: Unit test support
The project SHALL support unit testing with proper isolation and mocking.

#### Scenario: Standard testing library
- **WHEN** running unit tests
- **THEN** tests SHALL use Go's standard `testing` package
- **THEN** test functions SHALL follow `func TestXxx(t *testing.T)` convention

#### Scenario: Mock interfaces
- **WHEN** testing components with external dependencies
- **THEN** interfaces SHALL be defined for mockable dependencies
- **THEN** mock implementations SHALL be available for testing

#### Scenario: Table-driven tests
- **WHEN** testing multiple scenarios
- **THEN** tests SHALL use table-driven test patterns where appropriate
- **THEN** each test case SHALL have descriptive name

### Requirement: Integration test support
The project SHALL support integration tests for end-to-end scenarios.

#### Scenario: Integration test separation
- **WHEN** running integration tests
- **THEN** integration tests SHALL use build tags `// +build integration`
- **THEN** integration tests SHALL be skippable with standard `go test` command

#### Scenario: Test MCP server
- **WHEN** integration testing MCP interactions
- **THEN** a mock/test MCP server SHALL be available
- **THEN** tests SHALL be able to verify command translation

### Requirement: Test utilities and helpers
The project SHALL provide testing utilities for common test scenarios.

#### Scenario: Configuration test helpers
- **WHEN** testing configuration parsing
- **THEN** helper functions SHALL create test configurations
- **THEN** helpers SHALL support inline JSON or file-based configs

#### Scenario: Discord interaction test helpers
- **WHEN** testing Discord command handling
- **THEN** helper functions SHALL create mock Discord interactions
- **THEN** helpers SHALL verify interaction responses

### Requirement: Code coverage measurement
The project SHALL support measuring and reporting test coverage.

#### Scenario: Coverage reporting
- **WHEN** running tests with coverage
- **THEN** `go test -cover` SHALL work for all packages
- **THEN** coverage reports SHALL be generatable in multiple formats

#### Scenario: Coverage targets
- **WHEN** evaluating test completeness
- **THEN** critical packages SHALL aim for >80% coverage
- **THEN** translation and configuration parsing SHALL have comprehensive test coverage

### Requirement: Code coverage enforcement
The project SHALL enforce minimum code coverage thresholds in CI.

#### Scenario: Coverage threshold in GitHub Actions
- **WHEN** tests run in GitHub Actions
- **THEN** CI SHALL measure code coverage across all packages
- **THEN** CI SHALL fail if overall coverage is below 95%

#### Scenario: Coverage reporting in CI
- **WHEN** coverage is measured in CI
- **THEN** coverage percentage SHALL be reported in the workflow output
- **THEN** coverage report SHALL show per-package breakdown

#### Scenario: Pull request coverage checks
- **WHEN** a pull request is created
- **THEN** GitHub Actions SHALL verify coverage meets the 95% threshold
- **THEN** PR status check SHALL block merge if coverage is below threshold

### Requirement: Continuous Integration ready
The project SHALL include CI configuration for automated testing.

#### Scenario: Test automation
- **WHEN** code is pushed to repository
- **THEN** tests SHALL run automatically via CI
- **THEN** CI configuration SHALL run unit tests by default
- **THEN** CI configuration SHALL optionally run integration tests

#### Scenario: Multiple Go versions
- **WHEN** running CI tests
- **THEN** tests SHALL run on multiple Go versions
- **THEN** minimum supported Go version SHALL be documented
