# Resource groups module
module "resourcegroup" {
<<<<<<< HEAD
  source              = "./modules/resourcegroup"
  for_each            = var.resource_group_creation_enabled ? var.resource_groups : {}
  subscription_id     = local.subscription_id
  location            = coalesce(each.value.location, var.location)
  resource_group_name = each.value.name
  subscription_id     = local.subscription_id
  lock_enabled        = each.value.lock_enabled
  lock_name           = each.value.lock_name
  tags                = each.value.tags
}

# v3.3.0 introduced networkwatcherrg support,
# this was then moved into a more general resourcegroups module in later versions
=======
  source   = "./modules/resourcegroup"
  for_each = var.resource_group_creation_enabled ? var.resource_groups : {}

  location            = each.value.location
  resource_group_name = each.value.name
  subscription_id     = local.subscription_id
  lock_enabled        = each.value.lock_enabled
  lock_name           = each.value.lock_name
  tags                = each.value.tags
}

# v3.3.0 introduced networkwatcherrg support,
# this was then moved into a more general resourcegroups module in later versions
>>>>>>> e2d7552 (fix line endings)
moved {
  from = module.networkwatcherrg[0].azapi_resource.network_watcher_rg
  to   = module.resourcegroup_networkwatcherrg[0].azapi_resource.rg
}

# v5.0.0 consolidated the resource groups into a single module call
moved {
  from = module.resourcegroup_networkwatcherrg[0].azapi_resource.rg
  to   = module.resourcegroup["NetworkWatcherRG"].azapi_resource.rg
}
