# production-dockerfile Specification

## Purpose
TBD - created by archiving change terraform-aws-deployment. Update Purpose after archive.
## Requirements
### Requirement: Production Dockerfile at repo root
The repository SHALL include a production-grade `Dockerfile` at the repo root that builds and packages the bot for deployment.

#### Scenario: Multi-stage build
- **WHEN** `docker build` is run
- **THEN** the build SHALL use a multi-stage process: Go builder stage followed by a minimal runtime stage
- **THEN** the final image SHALL NOT include Go toolchain or source code

#### Scenario: Embedded runtime assets
- **WHEN** the Docker image is built
- **THEN** the image SHALL include an MCP server executable artifact of approximately 50 MB
- **THEN** the image SHALL include a SQLite database file of approximately 20 MB at `/app/data/bot.db`
- **THEN** the image SHALL fail build-time validation if either asset is missing

#### Scenario: Non-root user
- **WHEN** the container starts
- **THEN** the process SHALL run as a non-root user (`mcpdiscord`, uid 1000)
- **THEN** the `/data` directory SHALL be owned by this user for SQLite writes

### Requirement: MCP server process co-location
The container SHALL support launching an MCP server process alongside the bot within the same container.

#### Scenario: Entrypoint script launches both processes
- **WHEN** the container starts with an MCP server command configured
- **THEN** the entrypoint script SHALL start the MCP server process
- **THEN** the entrypoint script SHALL start the bot binary
- **THEN** if the MCP server process exits, the entrypoint script SHALL exit with a non-zero code so ECS restarts the task

#### Scenario: MCP server configured via environment
- **WHEN** the container is started with `MCP_CONFIG_JSON` environment variable set
- **THEN** the bot SHALL use that JSON as its configuration without requiring a mounted config file

### Requirement: SQLite data directory
The container SHALL expose a `/data` directory for SQLite database storage.

#### Scenario: Data directory creation
- **WHEN** the Docker image is built
- **THEN** an `/app/data` directory SHALL exist in the image with correct ownership
- **THEN** the bot SHALL default to storing its SQLite database at `/app/data/bot.db`

#### Scenario: In-container writable database
- **WHEN** the bot modifies `/app/data/bot.db` at runtime
- **THEN** writes SHALL occur on the container writable layer without requiring an external volume
- **THEN** the container SHALL have sufficient free filesystem capacity for the embedded SQLite file and expected growth

### Requirement: Health check
The container SHALL define a Docker health check.

#### Scenario: Process-based health check
- **WHEN** Docker or ECS evaluates container health
- **THEN** the health check SHALL verify the bot process (`mcpdiscord`) is running
- **THEN** unhealthy containers SHALL be restarted by the orchestrator

