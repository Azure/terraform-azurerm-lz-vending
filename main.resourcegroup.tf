# Resource groups module
module "resourcegroup" {
  source = "./modules/resourcegroup"
  for_each = merge(
    var.network_watcher_resource_group_enabled && var.resource_group_creation_enabled ? {
      NetworkWatcherRG = {
        name     = "NetworkWatcherRG"
        location = var.location
      },
    } : {},
    var.resource_group_creation_enabled ? var.resource_groups : {}
  )
  subscription_id     = local.subscription_id
  location            = each.value.location
  resource_group_name = each.value.name
  tags                = each.value.tags
}
