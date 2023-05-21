# Virtual hub submodule, disabled by default
# Will add vhub specific configuration, intent based routing
module "virtualhub" {
  source = "./modules/virtualhub"
  count  = var.virtual_hub_enabled ? 1 : 0

  subscription_id = local.subscription_id
  virtual_hubs    = var.virtual_hubs
}
