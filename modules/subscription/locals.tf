locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(azurerm_subscription.this[0].subscription_id, jsondecode(azapi_resource.subscription[0].output).properties.subscriptionId, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = coalesce(local.subscription_id_alias, var.subscription_id)
}

locals {
  # Check for drift between subscription and target management group.
  desired_subscription_management_group_association_exists = var.subscription_management_group_association_enabled && var.subscription_use_azapi ? contains(jsondecode(data.azapi_resource_list.subscription_management_group_association[0].output).value[*].id, "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}/subscriptions/${jsondecode(azapi_resource.subscription[0].output).properties.subscriptionId}") : false
}
