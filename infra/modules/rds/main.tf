resource "aws_db_subnet_group" "this" {
  name       = "${var.project}-${var.environment}"
  subnet_ids = var.private_subnet_ids
}

resource "aws_security_group" "rds" {
  name   = "${var.project}-${var.environment}-rds"
  vpc_id = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.eks_node_sg_id]
  }

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_db_instance" "this" {
  identifier        = "${var.project}-${var.environment}"
  engine            = "postgres"
  engine_version    = "16"
  instance_class    = var.db_instance_class
  allocated_storage = 20
  storage_type      = "gp3"
  storage_encrypted = true

  db_name  = "urllengthener"
  username = "urllengthener"
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  skip_final_snapshot     = true
  backup_retention_period = 7
  deletion_protection     = false

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}
