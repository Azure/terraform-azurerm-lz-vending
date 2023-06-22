# Network Watcher Resource Group module
module "networkwatcherrg" {
  source = "./modules/networkwatcherrg"
  count  = var.network_watcher_resource_group_enabled ? 1 : 0

  subscription_id = local.subscription_id
  location        = var.location
}
