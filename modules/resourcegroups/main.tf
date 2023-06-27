resource "azapi_resource" "rg" {
  for_each  = var.resource_groups_to_create
  parent_id = "/subscriptions/${var.subscription_id}"
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = each.value.name
  location  = each.value.location
  body      = jsonencode({})
  tags      = each.value.tags
}
