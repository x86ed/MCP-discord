# aws-container-deployment Specification

## Purpose
TBD - created by archiving change terraform-aws-deployment. Update Purpose after archive.
## Requirements
### Requirement: ECS Fargate cluster and service
The Terraform module SHALL provision an AWS ECS Fargate cluster and service to run the bot container.

#### Scenario: Cluster creation
- **WHEN** `terraform apply` is run
- **THEN** an ECS cluster SHALL be created with the configured name
- **THEN** the cluster SHALL use FARGATE capacity provider

#### Scenario: Service desired count
- **WHEN** the ECS service is created
- **THEN** it SHALL run a configurable number of tasks (default: 1)
- **THEN** ECS SHALL automatically restart tasks that exit or fail health checks

#### Scenario: Task definition resources
- **WHEN** the ECS task definition is registered
- **THEN** it SHALL allocate at least 256 CPU units and 512 MB memory (configurable)
- **THEN** the task SHALL run as a non-root user

### Requirement: ECR repository for container images
The Terraform module SHALL provision an AWS ECR private repository to store Docker images.

#### Scenario: Repository creation
- **WHEN** `terraform apply` is run
- **THEN** an ECR repository SHALL be created with image scanning on push enabled
- **THEN** a lifecycle policy SHALL retain only the last 10 images

#### Scenario: Image pull from ECS
- **WHEN** ECS launches a task
- **THEN** it SHALL pull the image from the ECR repository using the task execution IAM role

### Requirement: ECS ephemeral storage sizing for embedded assets
The Terraform module SHALL configure ECS task ephemeral storage to support a self-contained image with embedded MCP binary and SQLite DB files.

#### Scenario: Ephemeral storage configuration
- **WHEN** `terraform apply` is run
- **THEN** the ECS task definition SHALL explicitly set ephemeral storage size in GiB
- **THEN** default ephemeral storage SHALL be at least 21 GiB

#### Scenario: Embedded asset runtime support
- **WHEN** the ECS task starts
- **THEN** the container SHALL support an embedded MCP binary of approximately 50 MB and a SQLite DB of approximately 20 MB
- **THEN** the container SHALL be able to read and write the SQLite DB at its configured in-container path

### Requirement: IAM roles and policies
The Terraform module SHALL create least-privilege IAM roles for ECS task execution and task runtime.

#### Scenario: Task execution role
- **WHEN** ECS launches a task
- **THEN** the execution role SHALL allow ECR image pull and CloudWatch Logs write
- **THEN** no additional permissions SHALL be granted beyond these

#### Scenario: Task role
- **WHEN** the application code runs inside the task
- **THEN** the task role SHALL only include permissions required by runtime integrations configured for the bot
- **THEN** no broad AWS permissions SHALL be granted to the application

### Requirement: VPC networking for ECS tasks
The Terraform module SHALL configure VPC networking so ECS tasks can reach Discord and MCP server endpoints.

#### Scenario: Security group outbound access
- **WHEN** the ECS task needs to connect to Discord's API and any external MCP server
- **THEN** the security group SHALL allow all outbound traffic on port 443
- **THEN** no inbound traffic SHALL be permitted to the task from the internet

### Requirement: Terraform variables and outputs
The Terraform module SHALL expose configurable variables and useful outputs.

#### Scenario: Required variables
- **WHEN** running `terraform apply`
- **THEN** the user SHALL be able to configure: AWS region, VPC CIDR, ECR repo name, ECS cluster name, ECS service name, task CPU/memory, desired task count, Discord bot token (as sensitive variable), and task ephemeral storage size in GiB

#### Scenario: Outputs
- **WHEN** `terraform apply` completes
- **THEN** it SHALL output: ECR repository URL, ECS cluster ARN, ECS service ARN

