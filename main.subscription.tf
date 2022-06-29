# The subscription submodule creates a new subscription alias
# If we don't use this module, supply the `subscription_id` variable
# to be able to deploy resources to an existing subscription.
module "subscription" {
  source = "./modules/subscription"
  count  = var.subscription_alias_enabled || var.subscription_management_group_association_enabled ? 1 : 0

  subscription_alias_enabled                        = var.subscription_alias_enabled
  subscription_alias_name                           = var.subscription_alias_name
  subscription_billing_scope                        = var.subscription_billing_scope
  subscription_display_name                         = var.subscription_display_name
  subscription_management_group_association_enabled = var.subscription_management_group_association_enabled
  subscription_management_group_id                  = var.subscription_management_group_id
  subscription_workload                             = var.subscription_workload
  subscription_tags                                 = var.subscription_tags
}
