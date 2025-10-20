output "rds_endpoint" {
  description = "Auth service RDS PostgreSQL endpoint"
  value       = aws_db_instance.auth_postgres.endpoint
  sensitive   = true
}

output "rds_address" {
  description = "Auth service RDS PostgreSQL address"
  value       = aws_db_instance.auth_postgres.address
  sensitive   = true
}

output "redis_endpoint" {
  description = "Auth service Redis endpoint for token cache"
  value       = aws_elasticache_cluster.auth_redis.cache_nodes[0].address
  sensitive   = true
}

output "redis_port" {
  description = "Auth service Redis port"
  value       = aws_elasticache_cluster.auth_redis.cache_nodes[0].port
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.auth_postgres.db_name
}

# Outputs from shared infrastructure (for reference)
output "eks_cluster_name" {
  description = "EKS cluster name from shared infrastructure"
  value       = data.terraform_remote_state.shared.outputs.cluster_name
}

output "vpc_id" {
  description = "VPC ID from shared infrastructure"
  value       = data.terraform_remote_state.shared.outputs.vpc_id
}

