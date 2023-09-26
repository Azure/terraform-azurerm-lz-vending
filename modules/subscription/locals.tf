locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(azurerm_subscription.this[0].subscription_id, jsondecode(azapi_resource.subscription[0].output).properties.subscriptionId, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = coalesce(local.subscription_id_alias, var.subscription_id)
}
