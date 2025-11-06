output "management_group_subscription_association_id" {
  description = <<DESCRIPTION
The management_group_subscription_association_id output is the ID of the management group subscription association.
Value will be null if `var.subscription_management_group_association_enabled` is false.
DESCRIPTION
  value       = var.subscription_management_group_association_enabled ? try(azapi_resource_action.subscription_association[0].id, null) : null
}

output "subscription_id" {
  description = <<DESCRIPTION
The subscription_id is the id of the newly created subscription, or that of the supplied var.subscription_id.
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.
DESCRIPTION
  value       = local.subscription_id
}

output "subscription_resource_id" {
  description = <<DESCRIPTION
The subscription_resource_id output is the Azure resource id for the newly created subscription.
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.
DESCRIPTION
  value       = local.subscription_id != null ? "/subscriptions/${local.subscription_id}" : null
}
