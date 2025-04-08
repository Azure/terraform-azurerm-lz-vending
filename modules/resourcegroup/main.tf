resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  body      = {}
  location  = var.location
  name      = var.resource_group_name
  parent_id = "/subscriptions/${var.subscription_id}"
  tags      = var.tags
}
