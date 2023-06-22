# The azurerm_subscription resource represents the subscription alias that is being created.
resource "azurerm_subscription" "this" {
  count             = var.subscription_alias_enabled ? 1 : 0
  subscription_name = var.subscription_display_name
  alias             = var.subscription_alias_name
  billing_scope_id  = var.subscription_billing_scope
  workload          = var.subscription_workload
  tags              = var.subscription_tags

  # This provisioner requires az cli to be installed and logged in.
  provisioner "local-exec" {
    when       = destroy
    command    = "az resource delete --ids /subscriptions/${self.subscription_id}/resourceGroups/NetworkWatcherRG"
    on_failure = continue
  }
}

# This resource ensures that we can manage the management group for the subscription
# throughout its lifecycle.
resource "azurerm_management_group_subscription_association" "this" {
  count               = var.subscription_management_group_association_enabled ? 1 : 0
  management_group_id = "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}"
  subscription_id     = "/subscriptions/${local.subscription_id}"
}

# Register resource providers
resource "azapi_resource_action" "resource_provider_registration" {
  for_each    = var.subscription_register_resource_providers_and_features
  type        = "Microsoft.Resources/subscriptions@2021-04-01"
  resource_id = "/subscriptions/${local.subscription_id}"
  action      = "providers/${each.key}/register"
  method      = "POST"
}

resource "azapi_resource_action" "resource_provider_feature_registration" {
  for_each    = local.resource_provider_feature_map
  type        = "${each.value.resource_provider_name}/features@2021-07-01"
  resource_id = "/subscriptions/${local.subscription_id}/providers/Microsoft.Features/providers/${each.value.resource_provider_name}/features/${each.value.feature_name}}"
  action      = "register"
  method      = "POST"
}
