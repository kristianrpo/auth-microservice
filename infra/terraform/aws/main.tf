terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "mycompany-terraform-state"
    key    = "auth-microservice/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

# Data source - Consume shared infrastructure from infrastructure-shared repo
data "terraform_remote_state" "shared" {
  backend = "s3"
  config = {
    bucket = var.tf_backend_bucket
    key    = "shared/terraform.tfstate"
    region = var.aws_region
  }
}

# RDS PostgreSQL - Specific for Auth Microservice
resource "aws_db_instance" "auth_postgres" {
  identifier        = "auth-postgres-${var.environment}"
  engine            = "postgres"
  engine_version    = "16.1"
  instance_class    = var.db_instance_class
  allocated_storage = 20
  storage_encrypted = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.auth_rds.id]
  db_subnet_group_name   = aws_db_subnet_group.auth_db.name

  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"

  skip_final_snapshot = var.environment != "production"

  tags = {
    Name        = "auth-postgres-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

# ElastiCache Redis - Specific for Auth Microservice (Token Cache)
resource "aws_elasticache_cluster" "auth_redis" {
  cluster_id           = "auth-redis-${var.environment}"
  engine               = "redis"
  engine_version       = "7.0"
  node_type            = var.redis_node_type
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379

  security_group_ids = [aws_security_group.auth_redis.id]
  subnet_group_name  = aws_elasticache_subnet_group.auth_cache.name

  tags = {
    Name        = "auth-redis-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

# Security Group for RDS
resource "aws_security_group" "auth_rds" {
  name        = "auth-rds-sg-${var.environment}"
  description = "Security group for Auth service RDS"
  vpc_id      = data.terraform_remote_state.shared.outputs.vpc_id

  ingress {
    description     = "PostgreSQL from EKS nodes"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [data.terraform_remote_state.shared.outputs.eks_node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "auth-rds-sg-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

# Security Group for Redis
resource "aws_security_group" "auth_redis" {
  name        = "auth-redis-sg-${var.environment}"
  description = "Security group for Auth service Redis"
  vpc_id      = data.terraform_remote_state.shared.outputs.vpc_id

  ingress {
    description     = "Redis from EKS nodes"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [data.terraform_remote_state.shared.outputs.eks_node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "auth-redis-sg-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

# DB Subnet Group - Uses shared private subnets
resource "aws_db_subnet_group" "auth_db" {
  name       = "auth-db-subnet-group-${var.environment}"
  subnet_ids = data.terraform_remote_state.shared.outputs.private_subnet_ids

  tags = {
    Name        = "auth-db-subnet-group-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

# ElastiCache Subnet Group - Uses shared private subnets
resource "aws_elasticache_subnet_group" "auth_cache" {
  name       = "auth-cache-subnet-group-${var.environment}"
  subnet_ids = data.terraform_remote_state.shared.outputs.private_subnet_ids

  tags = {
    Name        = "auth-cache-subnet-group-${var.environment}"
    Environment = var.environment
    Service     = "auth-microservice"
  }
}

