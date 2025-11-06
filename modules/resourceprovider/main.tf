# Register resource providers
resource "azapi_resource_action" "resource_provider_registration" {
  action      = "providers/${var.resource_provider}/register"
  method      = "POST"
  resource_id = "/subscriptions/${var.subscription_id}"
  type        = "Microsoft.Resources/subscriptions@2021-04-01"
}

resource "azapi_resource_action" "resource_provider_feature_registration" {
  for_each = var.features

  action      = "register"
  method      = "POST"
  resource_id = "/subscriptions/${var.subscription_id}/providers/Microsoft.Features/providers/${var.resource_provider}/features/${each.value}"
  type        = "${var.resource_provider}/features@2021-07-01"
}
