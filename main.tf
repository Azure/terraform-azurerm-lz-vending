# The subsubscription_alias resource represents the subscription alias that is being created.
# This is used when subscription_id is not supplied, therefore we are
# creating a new subscription.
resource "azapi_resource" "subscription_alias" {
  count                   = var.subscription_alias_enabled && var.subscription_id == "" ? 1 : 0
  type                    = "Microsoft.Subscription/aliases@2021-10-01"
  parent_id               = "/"
  name                    = var.subscription_alias_name
  ignore_missing_property = true
  response_export_values = [
    "properties.subscriptionId",
  ]
  body = jsonencode({
    properties = {
      displayName  = var.subscription_alias_display_name
      billingScope = var.subscription_alias_billing_scope
      workload     = var.subscription_alias_workload
      additionalProperties = {
        managementGroupId = var.subscription_alias_management_group_id == "" ? null : local.subscription_alias_management_group_resource_id
      }
    }
  })
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
