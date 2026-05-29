# Terraform AWS Deployment

This directory provisions AWS infrastructure for running MCP Discord on ECS Fargate.

## What It Creates

- ECR repository for container images
- ECS cluster, task definition, and service
- VPC with 2 public subnets and internet gateway
- Security group with outbound HTTPS (443)
- IAM roles for ECS task execution and runtime
- CloudWatch log group

## Prerequisites

- Terraform >= 1.5
- AWS credentials with permissions for ECS, ECR, IAM, VPC, and CloudWatch

## Usage

```bash
cd infra
terraform init
terraform plan \
  -var="discord_bot_token=your-token-here" \
  -var="aws_region=us-east-1"
terraform apply \
  -var="discord_bot_token=your-token-here" \
  -var="aws_region=us-east-1"
```

## Key Variables

- `aws_region`: AWS region
- `vpc_cidr`: VPC CIDR range
- `ecr_repository_name`: ECR repository name
- `ecs_cluster_name`: ECS cluster name
- `ecs_service_name`: ECS service name
- `task_cpu`: Fargate CPU units
- `task_memory`: Fargate memory in MiB
- `desired_count`: Number of running tasks
- `ephemeral_storage_gib`: Task ephemeral storage in GiB (default: 21)
- `discord_bot_token` (sensitive): Discord bot token
- `mcp_config_json`: Optional inline MCP config JSON

## Outputs

- `ecr_repository_url`
- `ecs_cluster_arn`
- `ecs_service_arn`
