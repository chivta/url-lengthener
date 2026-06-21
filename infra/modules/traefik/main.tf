resource "helm_release" "traefik" {
  name             = "traefik"
  namespace        = "traefik"
  create_namespace = true
  upgrade_install  = true
  repository       = "https://traefik.github.io/charts"
  chart            = "traefik"
  version          = "41.0.0"

  # Single subnet -> single AZ -> single public IPv4 address.
  # A Classic ELB (the in-tree CCM default) allocates one public IP per AZ it spans;
  # the AWS Load Balancer Controller NLB here is pinned to one subnet to avoid that cost.
  set = [
    {
      name  = "service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-type"
      value = "external"
    },
    {
      name  = "service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-nlb-target-type"
      value = "ip"
    },
    {
      name  = "service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-scheme"
      value = "internet-facing"
    },
    {
      name  = "service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-subnets"
      value = var.public_subnet_ids[0]
    },
  ]
}
