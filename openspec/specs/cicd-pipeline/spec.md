# cicd-pipeline Specification

## Purpose
TBD - created by archiving change terraform-aws-deployment. Update Purpose after archive.
## Requirements
### Requirement: Automated build and push on push to main
The repository SHALL include a GitHub Actions workflow that builds and pushes the Docker image to ECR on every push to `main`.

#### Scenario: Workflow trigger
- **WHEN** a commit is pushed to the `main` branch
- **THEN** the `deploy.yml` workflow SHALL be triggered automatically
- **THEN** the workflow SHALL NOT run on pushes to other branches

#### Scenario: Image build and push
- **WHEN** the workflow runs
- **THEN** it SHALL authenticate to AWS using repository secrets
- **THEN** it SHALL log in to the ECR registry
- **THEN** it SHALL build the Docker image from the repo root `Dockerfile`
- **THEN** it SHALL tag the image with the Git commit SHA and `latest`
- **THEN** it SHALL push both tags to ECR

### Requirement: Automated ECS deployment after image push
The workflow SHALL trigger an ECS service update after pushing the new image.

#### Scenario: ECS task definition update
- **WHEN** the new image is pushed to ECR
- **THEN** the workflow SHALL render a new ECS task definition with the updated image URI
- **THEN** it SHALL register the new task definition revision with AWS

#### Scenario: ECS service rollout
- **WHEN** the new task definition is registered
- **THEN** the workflow SHALL update the ECS service to use the new task definition
- **THEN** ECS SHALL perform a rolling update, starting new tasks before stopping old ones
- **THEN** the workflow SHALL wait for the service to become stable before completing

### Requirement: Required GitHub Secrets
The CI/CD pipeline SHALL document and require specific GitHub Secrets to operate.

#### Scenario: AWS credential secrets
- **WHEN** the workflow configures AWS credentials
- **THEN** it SHALL read `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` from GitHub Secrets
- **THEN** it SHALL read `AWS_REGION` from GitHub Secrets

#### Scenario: Deployment target secrets
- **WHEN** the workflow pushes and deploys
- **THEN** it SHALL read `ECR_REPOSITORY` (repository name, not full URI) from GitHub Secrets
- **THEN** it SHALL read `ECS_CLUSTER` and `ECS_SERVICE` from GitHub Secrets

#### Scenario: Application secrets
- **WHEN** ECS runs the task
- **THEN** `DISCORD_BOT_TOKEN` SHALL be injected into the task environment from AWS Secrets Manager or ECS task definition environment (NOT stored in the image)

### Requirement: Terraform configuration validation
The CI/CD pipeline SHALL validate Terraform configuration files on every pull request and push to main or develop branches to catch infrastructure configuration errors early.

#### Scenario: Terraform validation job trigger
- **WHEN** code is pushed to main or develop branches OR a pull request is opened
- **THEN** a terraform-validate job SHALL run in the test workflow
- **AND** the job SHALL run independently of other test jobs

#### Scenario: Terraform syntax and consistency validation
- **WHEN** the terraform-validate job runs
- **THEN** it SHALL install Terraform CLI using hashicorp/setup-terraform action
- **AND** it SHALL run `terraform init` in the `/infra` directory
- **AND** it SHALL run `terraform validate` to check syntax and internal consistency
- **AND** it SHALL fail the workflow if validation finds errors

#### Scenario: Validation without AWS credentials
- **WHEN** terraform validation runs
- **THEN** it SHALL NOT require AWS credentials or secrets
- **AND** it SHALL only validate configuration syntax and consistency
- **AND** it SHALL NOT perform terraform plan or apply operations

### Requirement: Standardized Go version across toolchain
All build artifacts, CI workflows, and development configuration SHALL use Go 1.26 as the standard version to ensure consistency across environments.

#### Scenario: Go version in module definition
- **WHEN** the project is built or dependencies are resolved
- **THEN** go.mod SHALL declare `go 1.26` as the language version
- **AND** this SHALL be the minimum Go version required for development

#### Scenario: Go version in Docker builds
- **WHEN** the Docker image is built
- **THEN** Dockerfile SHALL use `golang:1.26-alpine` as the builder base image
- **AND** all Go binaries SHALL be compiled with Go 1.26 toolchain

#### Scenario: Go version in CI test workflows
- **WHEN** tests run in GitHub Actions
- **THEN** test.yml SHALL use Go version 1.26 for all test jobs
- **AND** matrix testing across multiple Go versions SHALL NOT be used
- **AND** coverage job SHALL explicitly use Go 1.26

#### Scenario: Consistent version across environments
- **WHEN** code is built in any environment (local, CI, Docker)
- **THEN** the same Go 1.26 version SHALL be used
- **AND** version inconsistencies SHALL be prevented by explicit version declarations

