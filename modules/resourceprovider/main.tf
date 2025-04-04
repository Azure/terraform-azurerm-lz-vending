# Register resource providers
resource "azapi_resource_action" "resource_provider_registration" {
  resource_id = "/subscriptions/${var.subscription_id}"
  type        = "Microsoft.Resources/subscriptions@2021-04-01"
  action      = "providers/${var.resource_provider}/register"
  method      = "POST"
}

resource "azapi_resource_action" "resource_provider_feature_registration" {
  for_each = var.features

  resource_id = "/subscriptions/${var.subscription_id}/providers/Microsoft.Features/providers/${var.resource_provider}/features/${each.value}"
  type        = "${var.resource_provider}/features@2021-07-01"
  action      = "register"
  method      = "POST"
}
