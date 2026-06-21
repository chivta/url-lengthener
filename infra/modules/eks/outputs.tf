output "cluster_name" { value = module.eks.cluster_name }
output "node_security_group_id" { value = module.eks.node_security_group_id }
output "lbc_iam_role_arn" { value = aws_iam_role.lbc.arn }
