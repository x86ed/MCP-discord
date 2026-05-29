## Why

The MCP-Discord bot currently runs as a local process with no production infrastructure. To operate reliably in production, it needs containerized packaging, persistent storage, and automated deployment to AWS — eliminating manual setup and enabling teams to ship updates continuously.

## What Changes

- **Dockerfile** (production-grade): Multi-stage build producing a minimal image that embeds the bot binary, an MCP server binary (~50 MB), and a seeded SQLite database (~20 MB)
- **In-container SQLite storage**: The bot uses a SQLite file inside the container filesystem at startup (default `/app/data/bot.db`) with no external database dependency
- **MCP server sidecar support**: The container entrypoint can launch a co-located MCP server process (stdio transport) configured via environment variables and the existing `mcp-config.json` contract
- **Terraform module** (`infra/`): Provisions AWS ECS Fargate task + service, ECR repository, IAM roles, VPC networking, and task ephemeral storage sized for embedded assets
- **CI/CD pipeline** (GitHub Actions): On push to `main`, builds and pushes the Docker image to ECR and triggers an ECS service update

## Capabilities

### New Capabilities

- `aws-container-deployment`: Infrastructure-as-code (Terraform) that provisions ECS Fargate, ECR, and networking to run a self-contained image on AWS
- `production-dockerfile`: A production-hardened Dockerfile that bundles an MCP server binary (~50 MB) and a SQLite database (~20 MB) inside the image, and launches both bot and MCP server in one container
- `cicd-pipeline`: GitHub Actions workflow that builds, tags, pushes the image to ECR, and deploys to ECS on every push to `main`

### Modified Capabilities

- `mcp-configuration`: The existing MCP config loader needs to support reading config from an environment variable (`MCP_CONFIG_JSON`) in addition to a file path, so that ECS task definitions can inject config without a mounted file

## Impact

- New top-level `infra/` directory containing Terraform modules
- New `.github/workflows/deploy.yml` CI/CD pipeline
- `Dockerfile` at repo root (replaces/supersedes `examples/deployment/Dockerfile`)
- New image asset packaging flow for bundled MCP server binary and seeded SQLite DB
- `internal/config/loader.go` modified to support env-var config injection
- Requires AWS credentials in GitHub Actions secrets: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `ECR_REPOSITORY`, `ECS_CLUSTER`, `ECS_SERVICE`
- Requires Discord secrets: `DISCORD_BOT_TOKEN` (and optionally `DISCORD_GUILD_ID`) in ECS task environment
