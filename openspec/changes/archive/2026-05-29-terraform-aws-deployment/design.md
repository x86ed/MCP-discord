## Context

The MCP-Discord bot currently has no production infrastructure. An example Dockerfile and docker-compose exist under `examples/deployment/` but are not production-grade and have no cloud deployment story. The user requires a self-contained runtime image that includes an embedded MCP server binary (~50 MB) and an embedded SQLite DB file (~20 MB), plus co-located process startup.

The change introduces three new layers: a production Docker image, AWS infrastructure via Terraform, and a GitHub Actions CI/CD pipeline that wires them together.

## Goals / Non-Goals

**Goals:**
- Single `Dockerfile` at the repo root that produces a production image containing the bot, MCP server binary (~50 MB), and SQLite DB (~20 MB)
- Terraform module under `infra/` that provisions ECS Fargate + ECR + IAM + networking with explicit ephemeral storage sizing
- GitHub Actions workflow that builds, pushes to ECR, and rolls out ECS on every push to `main`
- Config loader extension to accept `MCP_CONFIG_JSON` env var so ECS task definitions can inject config without a mounted file

**Non-Goals:**
- Multi-region or multi-environment Terraform (single AWS region, single environment)
- Aurora/RDS — not needed for this self-contained container model
- External persistent volumes (EFS/EBS) for initial rollout
- Kubernetes / EKS
- Blue/green or canary deployments — rolling ECS update is sufficient

## Decisions

### 1. ECS Fargate over EC2
**Decision**: Use ECS Fargate (serverless containers) rather than EC2 instances.  
**Rationale**: No AMI management, no SSH access, pay-per-task. The bot is a single long-running process with modest resource needs (256 CPU / 512 MB RAM).  
**Alternative considered**: EC2 Auto Scaling Group — more operational overhead for no benefit at this scale.

### 2. Embedded SQLite in container over external DB
**Decision**: Ship a seeded SQLite database file (~20 MB) inside the container image and run with container-local writable storage.  
**Rationale**: Meets the explicit requirement that SQLite lives inside the container, avoids external dependencies, and simplifies initial deployment.  
**Alternative considered**: EFS-mounted SQLite and RDS. Both add infrastructure complexity not required by current acceptance criteria.  
**Trade-off**: Data is not durable across task replacement unless exported externally. This is acceptable for current scope.

### 3. MCP server as sidecar process inside the same container
**Decision**: The Dockerfile `CMD` launches both the bot binary and the MCP server process (configured via `mcp-config.json` / `MCP_CONFIG_JSON`). A process supervisor (s6-overlay or a small shell wrapper) manages both.  
**Rationale**: ECS Fargate supports multiple containers per task via task definition `containerDefinitions`, but stdio transport requires the MCP server to be in the same process namespace as the bot. A sidecar container would require switching to SSE/WebSocket transport.  
**Alternative considered**: Separate ECS sidecar container with SSE transport — would require changes to the MCP transport layer and adds networking complexity.  
**Trade-off**: If the MCP server crashes, the supervisor restarts it; if the bot crashes, the whole task restarts. This is acceptable for a single-server deployment.

### 4. Embedded MCP server binary asset
**Decision**: Package the MCP server binary artifact (~50 MB) directly into the runtime image at build time, with checksum verification during CI build.  
**Rationale**: Ensures the runtime environment is self-contained and predictable, and removes startup-time dependency downloads.  
**Alternative considered**: Downloading the binary at container startup. This increases cold-start time and introduces runtime network dependency.

### 5. `MCP_CONFIG_JSON` env var for config injection
**Decision**: Extend `internal/config/loader.go` to check `MCP_CONFIG_JSON` environment variable (raw JSON string) before falling back to file path.  
**Rationale**: ECS task definitions pass configuration via environment variables. Mounting a JSON config file in ECS requires either a secrets manager integration or a baked-in file — both are more complex than a single env var.  
**Alternative considered**: AWS Secrets Manager with a config file template — higher complexity, slower startup.

### 6. Terraform structure
**Decision**: Single `infra/` module (not a module registry pattern) with files split by concern: `ecr.tf`, `ecs.tf`, `iam.tf`, `networking.tf`, `variables.tf`, `outputs.tf`.  
**Rationale**: Simplest structure for a single-environment deployment. The user can promote to a module pattern later.

### 7. CI/CD with GitHub Actions
**Decision**: `.github/workflows/deploy.yml` triggers on push to `main`. Steps: checkout → configure AWS credentials → login to ECR → build & push image → update ECS service.  
**Rationale**: GitHub Actions is already available in this repo. Uses official AWS actions (`aws-actions/configure-aws-credentials`, `aws-actions/amazon-ecr-login`, `aws-actions/amazon-ecs-deploy-task-definition`).

## Risks / Trade-offs

- **Image bloat from embedded assets** → Runtime image grows by at least ~70 MB from MCP binary + SQLite DB. Mitigated by multi-stage build and slim runtime base.
- **Single Fargate task = single point of failure** → Acceptable for a Discord bot; ECS will restart on crash. For HA, set `desiredCount = 2` (requires SQLite → shared DB migration).
- **Secrets in GitHub Actions** → Discord token and AWS credentials stored as GitHub Secrets. Ensure repo has branch protection on `main`.
- **Terraform state** → `infra/` uses local state by default. Teams should migrate to S3 backend with DynamoDB locking before sharing infra management.
- **MCP server process management** → If the MCP server exits unexpectedly inside the container, the bot will fail tool calls. The shell wrapper must detect this and restart or exit (causing ECS to restart the task).
- **Container-local SQLite durability** → Task replacement resets writable layer state unless DB is re-baked or exported. Mitigated by explicitly documenting this behavior and adding backup/export hooks as a follow-up.

## Migration Plan

1. **Build & push initial image** (CI/CD first run on merge to `main`)
2. **Apply Terraform** (`terraform init && terraform apply` from `infra/`) — provisions ECR, ECS cluster, and networking
3. **Set GitHub Secrets** — `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `ECR_REPOSITORY`, `ECS_CLUSTER`, `ECS_SERVICE`, `DISCORD_BOT_TOKEN`
4. **Verify ECS task is running** via AWS Console or `aws ecs describe-services`
5. **Rollback**: `terraform destroy` removes all AWS resources; previous bot deployment was manual so no rollback needed

## Open Questions

- Should SQLite data survive task replacement in v1? (Current plan: no; later enhancement can move to EFS or export/restore job.)
- Does the user want a VPC with private subnets + NAT gateway (higher cost) or public subnets with security groups (simpler, lower cost)? (Default design: public subnets with restrictive security groups — outbound-only for the bot.)
