data "azapi_resource_list" "subscription_management_group_association" {
  count = var.subscription_management_group_association_enabled ? 1 : 0

  parent_id              = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  type                   = "Microsoft.Management/managementGroups/subscriptions@2020-05-01"
  response_export_values = ["*"]
}

data "azapi_resource_list" "subscriptions" {
  count = var.subscription_management_group_association_enabled ? 1 : 0

  parent_id              = "/"
  type                   = "Microsoft.Resources/subscriptions@2022-12-01"
  response_export_values = ["*"]
}
