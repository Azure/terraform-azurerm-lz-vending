# The subscription submodule creates a new subscription alias
# If we don't use this module, supply the `subscription_id` variable
# to be able to deploy resources to an existing subscription.
module "subscription" {
  source                                 = "./modules/subscription"
  count                                  = var.subscription_alias_enabled ? 1 : 0

  # Required variables
  subscription_alias_billing_scope       = var.subscription_alias_billing_scope
  subscription_alias_workload            = var.subscription_alias_workload

  subscription_alias_display_name        = var.subscription_alias_display_name
  subscription_alias_management_group_id = var.subscription_alias_management_group_id
}
