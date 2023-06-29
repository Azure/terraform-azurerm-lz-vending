# Resource groups module
module "resourcegroup" {
  source              = "./modules/resourcegroup"
  for_each            = var.resource_group_creation_enabled ? var.resource_groups : {}
  subscription_id     = local.subscription_id
  location            = each.value.location
  resource_group_name = each.value.name
  tags                = each.value.tags
}

# Resource groups module for network watcher
module "resourcegroup_networkwatcherrg" {
  source              = "./modules/resourcegroup"
  count               = var.network_watcher_resource_group_enabled ? 1 : 0
  subscription_id     = local.subscription_id
  location            = var.location
  resource_group_name = "NetworkWatcherRG"
  tags                = {}
}

moved {
  from = module.networkwatcherrg.azapi_resource.network_watcher_rg
  to   = module.resourcegroup_networkwatcherrg.azapi_resource.rg
}
