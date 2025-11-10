# Test wrapper for the resourcegroup module
module "resourcegroup" {
  source = "../../modules/resourcegroup"

  resource_group_name = var.resource_group_name
  location            = var.location
  subscription_id     = var.subscription_id
  tags                = var.tags
  lock_enabled        = var.lock_enabled
}
