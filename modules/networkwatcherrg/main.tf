resource "azapi_resource" "network_watcher_rg" {
  parent_id = "/subscriptions/${var.subscription_id}"
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = var.networkwatcher_rg_name
  location  = var.location
  body      = jsonencode({})
  tags      = var.tags
}
