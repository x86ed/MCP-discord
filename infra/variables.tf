variable "aws_region" {
  type        = string
  description = "AWS region for deployment"
  default     = "us-east-1"
}

variable "name_prefix" {
  type        = string
  description = "Prefix for AWS resource names"
  default     = "mcpdiscord"
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR block for VPC"
  default     = "10.42.0.0/16"
}

variable "ecr_repository_name" {
  type        = string
  description = "ECR repository name"
  default     = "mcpdiscord"
}

variable "ecs_cluster_name" {
  type        = string
  description = "ECS cluster name"
  default     = "mcpdiscord-cluster"
}

variable "ecs_service_name" {
  type        = string
  description = "ECS service name"
  default     = "mcpdiscord-service"
}

variable "container_name" {
  type        = string
  description = "Container name in task definition"
  default     = "mcpdiscord"
}

variable "image_tag" {
  type        = string
  description = "Container image tag"
  default     = "latest"
}

variable "task_cpu" {
  type        = number
  description = "Fargate task CPU units"
  default     = 256
}

variable "task_memory" {
  type        = number
  description = "Fargate task memory in MiB"
  default     = 512
}

variable "desired_count" {
  type        = number
  description = "Desired ECS service task count"
  default     = 1
}

variable "ephemeral_storage_gib" {
  type        = number
  description = "ECS task ephemeral storage in GiB"
  default     = 21
}

variable "discord_bot_token" {
  type        = string
  description = "Discord bot token"
  sensitive   = true
}

variable "discord_guild_id" {
  type        = string
  description = "Optional Discord guild ID"
  default     = ""
}

variable "mcp_config_json" {
  type        = string
  description = "Optional inline MCP config JSON"
  default     = ""
}
