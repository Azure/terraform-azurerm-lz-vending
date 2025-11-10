# Test wrapper for the usermanagedidentity module
module "usermanagedidentity" {
  source = "../../modules/usermanagedidentity"

  name                                   = var.name
  location                               = var.location
  parent_id                              = var.parent_id
  tags                                   = var.tags
  federated_credentials_github           = var.federated_credentials_github
  federated_credentials_terraform_cloud  = var.federated_credentials_terraform_cloud
  federated_credentials_advanced         = var.federated_credentials_advanced
}
