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
