## 1. Config Loader Enhancement

- [x] 1.1 Add `MCP_CONFIG_JSON` env var check to `internal/config/loader.go` — read and parse JSON from env before falling back to file path
- [x] 1.2 Add log message at info level when config is loaded from `MCP_CONFIG_JSON`
- [x] 1.3 Add unit tests to `internal/config/loader_test.go` covering `MCP_CONFIG_JSON` precedence and fallback behavior

## 2. Production Dockerfile

- [x] 2.1 Create `Dockerfile` at repo root with multi-stage build (Go builder → alpine runtime)
- [x] 2.2 Package MCP server binary artifact (~50 MB) into runtime image at `/app/mcp-server` and ensure execute permissions
- [x] 2.3 Package seeded SQLite DB artifact (~20 MB) into image at `/app/data/bot.db` and ensure writable ownership for non-root user
- [x] 2.4 Write `entrypoint.sh` script that starts MCP server process (if configured) and the bot binary, exiting non-zero if either process dies
- [x] 2.5 Add Docker `HEALTHCHECK` instruction using `pgrep mcpdiscord`
- [x] 2.6 Add build-time validation for embedded MCP binary and SQLite DB presence/expected size range
- [x] 2.7 Verify image builds locally: `docker build -t mcpdiscord .`

## 3. Terraform Infrastructure

- [x] 3.1 Create `infra/` directory and `infra/versions.tf` with required provider versions (`aws >= 5.0`, `terraform >= 1.5`)
- [x] 3.2 Create `infra/variables.tf` with all configurable inputs (region, VPC CIDR, ECR repo name, ECS cluster/service names, task CPU/memory, desired count, Discord bot token as sensitive)
- [x] 3.3 Create `infra/networking.tf` — VPC, public subnets (2 AZs), internet gateway, route tables, ECS security group (egress 443)
- [x] 3.4 Create `infra/ecr.tf` — ECR private repository with image scanning and 10-image lifecycle policy
- [x] 3.5 Create `infra/iam.tf` — ECS task execution role (ECR pull + CloudWatch Logs) and task role (least privilege)
- [x] 3.6 Create `infra/ecs.tf` — ECS cluster, task definition (Fargate, env vars injected, ephemeral storage configured), ECS service with rolling update
- [x] 3.7 Create `infra/outputs.tf` — output ECR repository URL, ECS cluster ARN, ECS service ARN
- [x] 3.8 Validate Terraform: `terraform -chdir=infra init && terraform -chdir=infra validate`

## 4. CI/CD Pipeline

- [x] 4.1 Create `.github/workflows/deploy.yml` triggered on push to `main`
- [x] 4.2 Add step: configure AWS credentials using `aws-actions/configure-aws-credentials` with `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` secrets
- [x] 4.3 Add step: log in to ECR using `aws-actions/amazon-ecr-login`
- [x] 4.4 Add step: build and push Docker image tagged with Git commit SHA and `latest` to ECR
- [x] 4.5 Add step: render updated ECS task definition with new image URI using `aws-actions/amazon-ecs-render-task-definition`
- [x] 4.6 Add step: deploy updated task definition and wait for service stability using `aws-actions/amazon-ecs-deploy-task-definition`
- [x] 4.7 Document required GitHub Secrets in `docs/deployment.md` (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `ECR_REPOSITORY`, `ECS_CLUSTER`, `ECS_SERVICE`)

## 5. Documentation Updates

- [x] 5.1 Update `docs/deployment.md` with AWS deployment section: prerequisites, Terraform apply steps, GitHub Secrets setup, verification steps
- [x] 5.2 Add `infra/README.md` documenting Terraform variables and outputs
