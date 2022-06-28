locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(jsondecode(azapi_resource.subscription_alias[0].output).properties.subscriptionId, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = try(coalesce(local.subscription_id_alias, var.subscription_id), null)

  # management_group_resource_id_prefix is the prefix of the management group resource id.
  subscription_alias_management_group_resource_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
}
