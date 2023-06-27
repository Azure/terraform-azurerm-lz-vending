# Network Watcher Resource Group module
module "resourcegroups" {
  source          = "./modules/resourcegroups"
  count           = var.network_watcher_resource_group_enabled || var.resource_group_creation_enabled ? 1 : 0
  subscription_id = local.subscription_id
  resource_groups_to_create = merge(
    var.network_watcher_resource_group_enabled ? {
      NetworkWatcherRG = {
        name     = "NetworkWatcherRG"
        location = var.location
      }
    } : {},
    var.resource_groups_to_create
  )
}
