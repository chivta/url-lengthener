resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true
  upgrade_install  = true
  repository       = "https://charts.jetstack.io"
  chart            = "cert-manager"
  version          = "1.20.2"

  set = [
    {
      name  = "crds.enabled"
      value = "true"
    },
  ]
}
