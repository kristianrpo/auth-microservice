# Example Terraform variables file
# Copy this file to terraform.tfvars and fill in your values

aws_region         = "us-east-1"
environment        = "dev"
tf_backend_bucket  = "mycompany-terraform-state"

# Database Configuration
db_instance_class = "db.t3.micro"
db_name          = "authdb"
db_username      = "authuser"
db_password      = "CHANGE_THIS_PASSWORD"

# Redis Configuration
redis_node_type = "cache.t3.micro"

