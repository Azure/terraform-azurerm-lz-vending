output "subscription_id" {
  value = local.subscription_id_alias
  description = "The subscription_id is the id of the newly created subscription."
}

output "subscription_resource_id" {
  value = "/subscriptions/${local.subscription_id_alias}"
  description = "The subscription_resource_id output is the Azure resource id for the newly created subscription."
}
