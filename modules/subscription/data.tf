data "azapi_resource_list" "subscription_management_group_association" {
  count = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? 1 : 0

  type                   = "Microsoft.Management/managementGroups/subscriptions@2020-05-01"
  parent_id              = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  response_export_values = ["*"]
}

data "azapi_resource_list" "subscriptions" {
  count = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? 1 : 0

  type                   = "Microsoft.Resources/subscriptions@2022-12-01"
  parent_id              = "/"
  response_export_values = ["*"]
}
