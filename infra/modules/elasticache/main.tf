resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.project}-${var.environment}"
  subnet_ids = var.private_subnet_ids

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_security_group" "redis" {
  name   = "${var.project}-${var.environment}-redis"
  vpc_id = var.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [var.eks_node_sg_id]
  }

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_elasticache_replication_group" "this" {
  replication_group_id    = "${var.project}-${var.environment}"
  description             = "Redis for ${var.project} ${var.environment}"
  node_type               = "cache.t3.micro"
  engine_version          = "7.0"
  num_node_groups         = 1
  replicas_per_node_group = 0
  port                    = 6379

  subnet_group_name  = aws_elasticache_subnet_group.this.name
  security_group_ids = [aws_security_group.redis.id]

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}
