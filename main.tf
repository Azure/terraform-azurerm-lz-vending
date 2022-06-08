# The sub resource represents the subscription alias that is being created.
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
      displayName  = var.subscription_alias_display_name
      billingScope = var.subscription_alias_billing_scope
      workload     = var.subscription_alias_workload
    }
  })
}
