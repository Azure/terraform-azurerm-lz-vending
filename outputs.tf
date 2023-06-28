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

output "virtual_network_resource_group_ids" {
  value       = local.virtual_network_resource_group_ids
  description = "A map of resource group ids, keyed by the var.virtual_networks input map. Only populated if the virtualnetwork submodule is enabled."
}

output "management_group_subscription_association_id" {
  value       = var.subscription_management_group_association_enabled ? module.subscription[0].management_group_subscription_association_id : null
  description = <<DESCRIPTION
The management_group_subscription_association_id output is the ID of the management group subscription association.
Value will be null if `var.subscription_management_group_association_enabled` is false.
DESCRIPTION
}

output "umi_client_id" {
  description = <<DESCRIPTION
The client id of the user managed identity.
Value will be null if `var.umi_enabled` is false.
DESCRIPTION
  value       = one(module.usermanagedidentity).client_id
}

output "umi_tenant_id" {
  description = <<DESCRIPTION
The tenant id of the user managed identity.
Value will be null if `var.umi_enabled` is false.
DESCRIPTION
  value       = one(module.usermanagedidentity).tenant_id
}

output "umi_principal_id" {
  description = <<DESCRIPTION
The principal id of the user managed identity, sometimes known as the object id.
Value will be null if `var.umi_enabled` is false.
DESCRIPTION
  value       = one(module.usermanagedidentity).principal_id
}

output "umi_id" {
  description = <<DESCRIPTION
The Azure resource id of the user managed identity.
Value will be null if `var.umi_enabled` is false.
DESCRIPTION
  value       = one(module.usermanagedidentity).umi_id
}
