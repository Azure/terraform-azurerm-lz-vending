# The subsubscription_alias resource represents the subscription alias that is being created.
# This is used when subscription_id is not supplied, therefore we are
# creating a new subscription.
resource "azapi_resource" "subscription_alias" {
  count                   = var.subscription_alias_enabled ? 1 : 0
  type                    = "Microsoft.Subscription/aliases@2021-10-01"
  parent_id               = "/"
  name                    = var.subscription_alias_name
  ignore_missing_property = true
  response_export_values = [
    "properties.subscriptionId",
  ]
  body = jsonencode({
    properties = {
      displayName  = coalesce(var.subscription_display_name, var.subscription_alias_name)
      billingScope = var.subscription_billing_scope
      workload     = var.subscription_workload
      # Disabled as we use the azurerm provider to do this instead
      # additionalProperties = {
      #   managementGroupId = var.subscription_alias_management_group_id == "" ? null : local.subscription_alias_management_group_resource_id
      # }
    }
  })
}

# This resource ensures that we can manage the management group for the subscription
# throughout its lifecycle.
resource "azurerm_management_group_subscription_association" "this" {
  count               = var.subscription_management_group_association_enabled ? 1 : 0
  management_group_id = local.subscription_alias_management_group_resource_id
  subscription_id     = local.subscription_id_alias
}

# Creating an alias for an existing subscription is not currently supported.
# Need use case data to justify the effort in testing support.
#
# # The subscription_alias_existing resource represents the subscription
# # alias that is being created for an existing subscription
# resource "azapi_resource" "subscription_alias_existing" {
#   count                   = var.subscription_alias_enabled && var.subscription_id  != "" ? 1 : 0
#   type                    = "Microsoft.Subscription/aliases@2021-10-01"
#   parent_id               = "/"
#   name                    = var.subscription_alias_name
#   ignore_missing_property = true
#   response_export_values = [
#     "properties.subscriptionId",
#   ]
#   body = jsonencode({
#     properties = {
#       subscriptionId = var.subscription_id
#     }
#   })
# }
