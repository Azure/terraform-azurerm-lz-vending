resource "azapi_resource" "network_watcher_rg" {
  parent_id = var.subscription_resource_id
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = "NetworkWatcherRG"
  location  = var.location
  body      = jsonencode({})
  tags      = var.tags
}
