## ADDED Requirements

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
