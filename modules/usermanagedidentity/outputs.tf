output "client_id" {
  description = "The client id of the user managed identity"
  value       = local.umi_output.properties.clientId
}

output "tenant_id" {
  description = "The tenant id of the user managed identity"
  value       = local.umi_output.properties.tenantId
}

output "principal_id" {
  description = "The object id of the user managed identity"
  value       = local.umi_output.properties.principalId
}
