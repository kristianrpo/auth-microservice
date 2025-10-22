# ============================================================================
# Data source: Consume shared infrastructure from remote state (igual que tienes)
# ============================================================================
data "terraform_remote_state" "shared" {
  backend = "s3"
  config = {
    bucket = var.tf_backend_bucket
    key    = var.shared_state_key
    region = var.aws_region
  }
}

locals {
  name               = "${var.project}-${var.environment}"
  cluster_name       = data.terraform_remote_state.shared.outputs.cluster_name
  rabbitmq_url       = data.terraform_remote_state.shared.outputs.rabbitmq_amqp_url

  vpc_id             = data.terraform_remote_state.shared.outputs.vpc_id
  private_subnet_ids = data.terraform_remote_state.shared.outputs.private_subnet_ids

  oidc_provider_arn  = data.terraform_remote_state.shared.outputs.oidc_provider_arn
}

# Derivamos el CIDR de la VPC SIN tocar el shared
data "aws_vpc" "this" {
  id = local.vpc_id
}

# ============================================================================
# Recursos del micro: RDS + Redis + Secrets + IAM (barato para dev)
# ============================================================================

resource "random_id" "suffix" {
  byte_length = 2
}

resource "random_password" "db_password" {
  length  = 20
  special = true
}

resource "random_password" "redis_auth_token" {
  length  = 24
  special = true
}

# ----------------------------------------------------------------------------
# Secret principal de configuración de la aplicación (ya lo tenías)
# ----------------------------------------------------------------------------
resource "aws_secretsmanager_secret" "app" {
  name        = "${local.name}/application-${random_id.suffix.hex}"
  description = "Application configuration for ${local.name}"
}

# ----------------------------------------------------------------------------
# Security Groups (limitados al CIDR de la VPC)
# ----------------------------------------------------------------------------
resource "aws_security_group" "rds" {
  name        = "${local.name}-rds-sg"
  description = "SG for RDS PostgreSQL (${local.name})"
  vpc_id      = local.vpc_id
  tags        = { Name = "${local.name}-rds-sg" }
}

resource "aws_vpc_security_group_ingress_rule" "rds_ingress" {
  security_group_id = aws_security_group.rds.id
  cidr_ipv4         = data.aws_vpc.this.cidr_block
  ip_protocol       = "tcp"
  from_port         = 5432
  to_port           = 5432
}

resource "aws_vpc_security_group_egress_rule" "rds_egress" {
  security_group_id = aws_security_group.rds.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_security_group" "redis" {
  name        = "${local.name}-redis-sg"
  description = "SG for ElastiCache Redis (${local.name})"
  vpc_id      = local.vpc_id
  tags        = { Name = "${local.name}-redis-sg" }
}

resource "aws_vpc_security_group_ingress_rule" "redis_ingress" {
  security_group_id = aws_security_group.redis.id
  cidr_ipv4         = data.aws_vpc.this.cidr_block
  ip_protocol       = "tcp"
  from_port         = 6379
  to_port           = 6379
}

resource "aws_vpc_security_group_egress_rule" "redis_egress" {
  security_group_id = aws_security_group.redis.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

# ----------------------------------------------------------------------------
# RDS PostgreSQL (barato dev)
# ----------------------------------------------------------------------------
resource "aws_db_subnet_group" "rds" {
  name       = "${local.name}-rds-subnets"
  subnet_ids = local.private_subnet_ids
  tags       = { Name = "${local.name}-rds-subnets" }
}

resource "aws_db_parameter_group" "rds" {
  name   = "${local.name}-rds-params"
  family = "postgres16"

  parameter {
    name  = "log_min_duration_statement"
    value = "2000"
  }
}

resource "aws_secretsmanager_secret" "rds_credentials" {
  name        = "${local.name}/rds/postgresql-${random_id.suffix.hex}"
  description = "RDS PostgreSQL credentials for ${local.name}"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_initial" {
  secret_id     = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = "appuser"
    password = random_password.db_password.result
    engine   = "postgres"
    host     = null
    port     = 5432
    dbname   = "appdb"
  })
}

resource "aws_db_instance" "postgres" {
  identifier                 = "${local.name}-pg"
  engine                     = "postgres"
  engine_version             = "16.3"
  instance_class             = "db.t4g.micro"   # económico
  allocated_storage          = 20               # económico
  db_name                    = "appdb"
  username                   = "appuser"
  password                   = random_password.db_password.result
  db_subnet_group_name       = aws_db_subnet_group.rds.name
  vpc_security_group_ids     = [aws_security_group.rds.id]
  parameter_group_name       = aws_db_parameter_group.rds.name

  storage_encrypted          = true
  backup_retention_period    = 1
  deletion_protection        = false
  auto_minor_version_upgrade = true
  multi_az                   = false
  publicly_accessible        = false
  skip_final_snapshot        = true

  tags = { Name = "${local.name}-postgres" }
}

# Actualiza secret con el host al tener endpoint
resource "aws_secretsmanager_secret_version" "rds_credentials_with_host" {
  secret_id = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = "appuser"
    password = random_password.db_password.result
    engine   = "postgres"
    host     = aws_db_instance.postgres.address
    port     = 5432
    dbname   = "appdb"
  })
  depends_on = [aws_db_instance.postgres]
}

# ----------------------------------------------------------------------------
# ElastiCache Redis (barato dev)
# ----------------------------------------------------------------------------
resource "aws_elasticache_subnet_group" "redis" {
  name       = "${local.name}-redis-subnets"
  subnet_ids = local.private_subnet_ids
}

resource "aws_elasticache_parameter_group" "redis" {
  name   = "${local.name}-redis-params"
  family = "redis7"

  parameter {
    name  = "timeout"
    value = "0"
  }
}

resource "aws_secretsmanager_secret" "redis_auth" {
  name        = "${local.name}/redis/auth-${random_id.suffix.hex}"
  description = "Redis AUTH token for ${local.name}"
}

resource "aws_secretsmanager_secret_version" "redis_auth" {
  secret_id     = aws_secretsmanager_secret.redis_auth.id
  secret_string = jsonencode({ auth_token = random_password.redis_auth_token.result })
}


resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "${local.name}-redis"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis6.x"
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.redis.id]
}
# ----------------------------------------------------------------------------
# Secret con cadenas de conexión (para la app y/o External Secrets)
# ----------------------------------------------------------------------------
resource "aws_secretsmanager_secret" "app_connections" {
  name        = "${local.name}/connections-${random_id.suffix.hex}"
  description = "App connections for ${local.name}"
}

resource "aws_secretsmanager_secret_version" "app_connections" {
  secret_id     = aws_secretsmanager_secret.app_connections.id
  secret_string = jsonencode({
    DATABASE_URL = "postgres://appuser:${random_password.db_password.result}@${aws_db_instance.postgres.address}:5432/appdb"
    REDIS_URL    = "rediss://:${random_password.redis_auth_token.result}@${aws_elasticache_cluster.redis.cache_nodes[0].address}:6379"
    RABBITMQ_URL = local.rabbitmq_url
  })
}

# ----------------------------------------------------------------------------
# IAM para el servicio (IRSA) - acceso a estos secretos
# ----------------------------------------------------------------------------
data "aws_iam_policy_document" "service_policy" {
  statement {
    actions = [
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret"
    ]
    resources = [
      aws_secretsmanager_secret.app.arn,
      aws_secretsmanager_secret.rds_credentials.arn,
      aws_secretsmanager_secret.redis_auth.arn,
      aws_secretsmanager_secret.app_connections.arn
    ]
  }
}

resource "aws_iam_policy" "service" {
  name_prefix = "${local.name}-policy-"
  policy      = data.aws_iam_policy_document.service_policy.json
}

# Usa el módulo IRSA igual que en tu otro micro (ajusta namespace:sa si difiere)
module "irsa" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "~> 5.39"

  role_name = "${local.name}-service-irsa"
  oidc_providers = {
    main = {
      provider_arn               = local.oidc_provider_arn
      namespace_service_accounts = ["default:app-sa"]
    }
  }
  role_policy_arns = { service = aws_iam_policy.service.arn }
}

# ----------------------------------------------------------------------------
# IAM Policy para External Secrets -> acceso a los secretos de este micro
# (se adjunta al rol compartido del ESO que viene del shared)
# ----------------------------------------------------------------------------
data "aws_iam_policy_document" "external_secrets" {
  statement {
    actions   = ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"]
    resources = [
      aws_secretsmanager_secret.app.arn,
      aws_secretsmanager_secret.rds_credentials.arn,
      aws_secretsmanager_secret.redis_auth.arn,
      aws_secretsmanager_secret.app_connections.arn
    ]
  }
}

resource "aws_iam_policy" "external_secrets" {
  name_prefix = "${local.name}-external-secrets-policy-"
  policy      = data.aws_iam_policy_document.external_secrets.json
}

resource "aws_iam_role_policy_attachment" "eso_attach" {
  role       = data.terraform_remote_state.shared.outputs.eso_irsa_role_name
  policy_arn = aws_iam_policy.external_secrets.arn
}

# ============================================================================
# Outputs - recursos de este micro + conveniencia desde shared
# ============================================================================
output "rds_endpoint"                  { value = aws_db_instance.postgres.address }
output "redis_primary_endpoint"        { value = aws_elasticache_cluster.redis.cache_nodes[0].address }

output "rds_secret_arn"                { value = aws_secretsmanager_secret.rds_credentials.arn }
output "redis_auth_secret_arn"         { value = aws_secretsmanager_secret.redis_auth.arn }
output "app_connections_secret_arn"    { value = aws_secretsmanager_secret.app_connections.arn }
output "app_secret_arn"                { value = aws_secretsmanager_secret.app.arn }

output "rds_security_group_id"         { value = aws_security_group.rds.id }
output "redis_security_group_id"       { value = aws_security_group.redis.id }

output "irsa_role_arn"                 { value = module.irsa.iam_role_arn }

# Outputs de shared (conveniencia)
output "cluster_name"                  { value = local.cluster_name }
output "cluster_endpoint"              { value = data.terraform_remote_state.shared.outputs.cluster_endpoint }
output "cluster_ca_certificate"        { value = data.terraform_remote_state.shared.outputs.cluster_ca_certificate }
output "eso_irsa_role_arn"             { value = data.terraform_remote_state.shared.outputs.external_secrets_irsa_role_arn }
output "aws_lb_controller_role_arn"    { value = data.terraform_remote_state.shared.outputs.aws_load_balancer_controller_irsa_role_arn }
