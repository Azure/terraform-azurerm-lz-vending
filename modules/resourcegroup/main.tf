resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  body      = {}
  location  = var.location
  name      = var.resource_group_name
  parent_id = "/subscriptions/${var.subscription_id}"
  tags      = var.tags
}

resource "azapi_resource" "rg_lock" {
  count = var.lock_enabled ? 1 : 0

  type = "Microsoft.Authorization/locks@2020-05-01"
  body = {
    properties = {
      level = "CanNotDelete"
    }
  }
  name      = coalesce(var.lock_name, "lock-${azapi_resource.rg.name}")
  parent_id = azapi_resource.rg.id

  depends_on = [
    azapi_resource.rg,
  ]
}
