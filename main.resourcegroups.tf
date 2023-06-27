# Network Watcher Resource Group module
module "networkwatcherrg" {
  source = "./modules/networkwatcherrg"
  count  = var.network_watcher_resource_group_enabled || var.resource_group_creation_enabled ? 1 : 0
  resource_groups = merge(
    var.network_watcher_resource_group_enabled ? {
      name     = "NetworkWatcherRG"
      location = var.location
      tags     = {}
    } : {},
    var.resource_groups_to_create
  )
}
