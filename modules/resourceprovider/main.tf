# Register resource providers
resource "azapi_resource_action" "resource_provider_registration" {
  type        = "Microsoft.Resources/subscriptions@2021-04-01"
  resource_id = "/subscriptions/${var.subscription_id}"
  action      = "providers/${var.resource_provider}/register"
  method      = "POST"
}

resource "azapi_resource_action" "resource_provider_feature_registration" {
  for_each    = var.features
  type        = "${var.resource_provider}/features@2021-07-01"
  resource_id = "/subscriptions/${var.subscription_id}/providers/Microsoft.Features/providers/${each.value.resource_provider_name}/features/${each.value.feature_name}"
  action      = "register"
  method      = "POST"
}
