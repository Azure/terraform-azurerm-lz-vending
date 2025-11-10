# Test wrapper for the resourceprovider module
module "resourceprovider" {
  source = "../../modules/resourceprovider"

  resource_provider = var.resource_provider
  features          = var.features
  subscription_id   = var.subscription_id
}
