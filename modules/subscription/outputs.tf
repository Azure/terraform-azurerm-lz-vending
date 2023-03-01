output "subscription_id" {
  value       = local.subscription_id
  description = <<DESCRIPTION
The subscription_id is the id of the newly created subscription, or that of the supplied var.subscription_id.
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.
DESCRIPTION
}

output "subscription_resource_id" {
  value       = local.subscription_id != null ? "/subscriptions/${local.subscription_id}" : null
  description = <<DESCRIPTION
The subscription_resource_id output is the Azure resource id for the newly created subscription.
Value will be null if `var.subscription_id` is blank and `var.subscription_alias_enabled` is false.
DESCRIPTION
}

output "management_group_subscription_association_id" {
  value       = var.subscription_management_group_association_enabled ? azurerm_management_group_subscription_association.this[0].id : null
  description = <<DESCRIPTION
The management_group_subscription_association_id output is the ID of the management group subscription association.
Value will be null if `var.subscription_management_group_association_enabled` is false.
DESCRIPTION
}
