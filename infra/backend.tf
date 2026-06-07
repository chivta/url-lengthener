terraform {
  backend "s3" {
    bucket       = "url-lengthener-terraform-state"
    key          = "terraform.tfstate"
    region       = "eu-west-1"
    use_lockfile = true
  }
}