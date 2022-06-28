output "subscription_id" {
  value       = local.subscription_id
  description = "The subscription_id is the Azure subscription id that resources have been deployed into."
}

output "subscription_resource_id" {
  value       = local.subscription_resource_id
  description = "The subscription_resource_id is the Azure subscription resource id that resources have been deployed into"
}
