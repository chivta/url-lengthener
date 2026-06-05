output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "rds_endpoint" {
  value     = module.rds.endpoint
  sensitive = true
}

output "ecr_api_repository_url" {
  value = module.ecr.api_repository_url
}

output "ecr_frontend_repository_url" {
  value = module.ecr.frontend_repository_url
}
