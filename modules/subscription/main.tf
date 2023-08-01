# The azurerm_subscription resource represents the subscription alias that is being created.
resource "azurerm_subscription" "this" {
  count             = var.subscription_alias_enabled && !(var.subscription_use_azapi) ? 1 : 0
  subscription_name = var.subscription_display_name
  alias             = var.subscription_alias_name
  billing_scope_id  = var.subscription_billing_scope
  workload          = var.subscription_workload
  tags              = var.subscription_tags
}

# Optionally make use of AzAPI to create the subscription to allow creation without access
# to the default management group.
resource "azapi_resource" "this" {
  count = var.subscription_alias_enabled && var.subscription_use_azapi ? 1 : 0

  type      = "Microsoft.Subscription/aliases@2021-10-01"
  name      = var.subscription_alias_name
  parent_id = "/"

  body = jsonencode({
    properties = {
      displayName    = var.subscription_display_name
      workload       = var.subscription_workload
      billingScope   = var.subscription_billing_scope
      subscriptionId = null
      additionalProperties = {
        managementGroupId = var.subscription_management_group_id
        tags              = var.subscription_tags
      }
    }
  })
}

# This resource ensures that we can manage the management group for the subscription
# throughout its lifecycle.
resource "azurerm_management_group_subscription_association" "this" {
  count               = var.subscription_management_group_association_enabled ? 1 : 0
  management_group_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  subscription_id     = "/subscriptions/${local.subscription_id}"
}


