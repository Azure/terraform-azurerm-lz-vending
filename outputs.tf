output "subscription_id" {
  value       = local.subscription_id
  description = "The subscription_id is the Azure subscription id that resources have been deployed into."
}

output "subscription_resource_id" {
  value       = local.subscription_resource_id
  description = "The subscription_resource_id is the Azure subscription resource id that resources have been deployed into"
}

output "virtual_network_resource_ids" {
  value       = local.virtual_network_resource_ids
  description = "A map of virtual network resource ids, keyed by the var.virtual_networks input map. Only populated if the virtualnetwork submodule is enabled."
}

output "resource_group_ids" {
  value       = local.resource_group_ids
  description = "A map of resource group ids, keyed by the var.virtual_networks input map. Only populated if the virtualnetwork submodule is enabled."
}

output "management_group_subscription_association_id" {
  value       = var.subscription_management_group_association_enabled ? module.subscription[0].management_group_subscription_association_id : null
  description = <<DESCRIPTION
The management_group_subscription_association_id output is the ID of the management group subscription association.
Value will be null if `var.subscription_management_group_association_enabled` is false.
DESCRIPTION
}
