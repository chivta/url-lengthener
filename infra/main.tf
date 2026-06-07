locals {
  environment = terraform.workspace
}

module "vpc" {
  source = "./modules/vpc"

  project     = var.project
  environment = local.environment
  aws_region  = var.aws_region
}

module "eks" {
  source = "./modules/eks"

  project            = var.project
  environment        = local.environment
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  node_instance_type = var.eks_node_instance_type
}

module "rds" {
  source = "./modules/rds"

  project            = var.project
  environment        = local.environment
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  eks_node_sg_id     = module.eks.node_security_group_id
  db_instance_class  = var.db_instance_class
  db_password        = var.db_password
}

module "elasticache" {
  source = "./modules/elasticache"

  project            = var.project
  environment        = local.environment
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  eks_node_sg_id     = module.eks.node_security_group_id
}

module "ecr" {
  source = "./modules/ecr"
  count  = terraform.workspace == "prod" ? 1 : 0

  project     = var.project
  environment = local.environment
}

module "s3" {
  source = "./modules/s3"

  project     = var.project
  environment = local.environment
}
