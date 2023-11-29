# The azurerm_subscription resource represents the subscription alias that is being created.
resource "azurerm_subscription" "this" {
  count             = var.subscription_alias_enabled && !var.subscription_use_azapi ? 1 : 0
  subscription_name = var.subscription_display_name
  alias             = var.subscription_alias_name
  billing_scope_id  = var.subscription_billing_scope
  workload          = var.subscription_workload
  tags              = var.subscription_tags
}

# This resource ensures that we can manage the management group for the subscription
# throughout its lifecycle.
resource "azurerm_management_group_subscription_association" "this" {
  count               = var.subscription_management_group_association_enabled ? 1 : 0
  management_group_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  subscription_id     = "/subscriptions/${local.subscription_id}"
}

resource "azapi_resource" "subscription" {
  count = var.subscription_alias_enabled && var.subscription_use_azapi ? 1 : 0

  type      = "Microsoft.Subscription/aliases@2021-10-01"
  name      = var.subscription_alias_name
  parent_id = "/"

  body = jsonencode({
    properties = {
      displayName  = var.subscription_display_name
      workload     = var.subscription_workload
      billingScope = var.subscription_billing_scope
      additionalProperties = {
        managementGroupId = var.subscription_management_group_association_enabled ? "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}" : null
        tags              = var.subscription_tags
      }
    }
  })
  response_export_values = ["properties.subscriptionId"]
  lifecycle {
    ignore_changes = [
      body
    ]
  }
}

resource "azapi_update_resource" "subscription_tags" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) || (var.subscription_id != "" && var.subscription_update_existing) ? 1 : 0

  type        = "Microsoft.Resources/tags@2022-09-01"
  resource_id = "/subscriptions/${local.subscription_id}/providers/Microsoft.Resources/tags/default"
  body = jsonencode({
    properties = {
      tags = var.subscription_tags
    }
  })
}

resource "azapi_resource_action" "subscription_rename" {
  count = (var.subscription_alias_enabled && var.subscription_use_azapi) || (var.subscription_id != "" && var.subscription_update_existing) ? 1 : 0

  type        = "Microsoft.Resources/subscriptions@2021-10-01"
  resource_id = "/subscriptions/${local.subscription_id}"
  method      = "POST"
  action      = "providers/Microsoft.Subscription/rename"
  body = jsonencode({
    subscriptionName = var.subscription_display_name
  })
}
