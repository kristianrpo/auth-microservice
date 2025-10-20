variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment (dev, staging, production)"
  type        = string
  default     = "dev"
}

variable "tf_backend_bucket" {
  description = "S3 bucket for Terraform remote state (shared infrastructure)"
  type        = string
  default     = "mycompany-terraform-state"
}

variable "db_instance_class" {
  description = "RDS instance class for Auth service"
  type        = string
  default     = "db.t3.micro"
}

variable "db_name" {
  description = "Database name for Auth service"
  type        = string
  default     = "authdb"
}

variable "db_username" {
  description = "Database username for Auth service"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password for Auth service"
  type        = string
  sensitive   = true
}

variable "redis_node_type" {
  description = "Redis node type for token cache"
  type        = string
  default     = "cache.t3.micro"
}

