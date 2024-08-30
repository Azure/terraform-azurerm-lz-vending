# The subscription submodule creates a new subscription alias
# If we don't use this module, supply the `subscription_id` variable
# to be able to deploy resources to an existing subscription.
module "subscription" {
  source = "./modules/subscription"
  count  = (var.subscription_id != "" && var.subscription_update_existing) || var.subscription_alias_enabled || var.subscription_management_group_association_enabled ? 1 : 0

  subscription_alias_enabled                           = var.subscription_alias_enabled
  subscription_alias_name                              = var.subscription_alias_name
  subscription_billing_scope                           = var.subscription_billing_scope
  subscription_display_name                            = var.subscription_display_name
  subscription_id                                      = var.subscription_id
  subscription_management_group_association_enabled    = var.subscription_management_group_association_enabled
  subscription_management_group_id                     = var.subscription_management_group_id
  subscription_tags                                    = var.subscription_tags
  subscription_tags_ignored                            = var.subscription_tags_ignored
  subscription_use_azapi                               = var.subscription_use_azapi
  subscription_update_existing                         = var.subscription_update_existing
  subscription_workload                                = var.subscription_workload
  wait_for_subscription_before_subscription_operations = var.wait_for_subscription_before_subscription_operations
}
