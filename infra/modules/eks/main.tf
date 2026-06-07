module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name    = "${var.project}-${var.environment}"
  kubernetes_version = "1.33"

  vpc_id                         = var.vpc_id
  subnet_ids                     = var.private_subnet_ids
  endpoint_public_access = true
  enable_cluster_creator_admin_permissions = true

  tags = {
    Project     = var.project
    Environment = var.environment
  }

  addons = {
    vpc-cni = {
      most_recent    = true
      before_compute = true
    }
    kube-proxy = {
      most_recent = true
    }
    coredns = {
      most_recent = true
    }
  }

  eks_managed_node_groups = {
    default = {
      instance_types = [var.node_instance_type]
      min_size       = 1
      max_size       = 3
      desired_size   = 2
    }
  }
}
