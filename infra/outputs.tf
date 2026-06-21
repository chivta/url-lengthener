output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "lbc_iam_role_arn" {
  value = module.eks.lbc_iam_role_arn
}

output "rds_endpoint" {
  value     = module.rds.endpoint
  sensitive = true
}

output "ecr_api_repository_url" {
  value = one(module.ecr[*].api_repository_url)
}

output "ecr_frontend_repository_url" {
  value = one(module.ecr[*].frontend_repository_url)
}
