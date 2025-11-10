# Test wrapper for the subscription module
module "subscription" {
  source = "../../modules/subscription"

  # Subscription alias configuration
  subscription_alias_enabled = var.subscription_alias_enabled
  subscription_alias_name    = var.subscription_alias_name
  subscription_billing_scope = var.subscription_billing_scope
  subscription_display_name  = var.subscription_display_name
  subscription_workload      = var.subscription_workload

  # Subscription ID for existing subscriptions
  subscription_id = var.subscription_id

  # Management Group configuration
  subscription_management_group_id                  = var.subscription_management_group_id
  subscription_management_group_association_enabled = var.subscription_management_group_association_enabled

  # Tags
  subscription_tags = var.subscription_tags

  # Update existing subscription
  subscription_update_existing = var.subscription_update_existing
}
