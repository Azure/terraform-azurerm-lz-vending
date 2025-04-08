output "role_assignment_name" {
  description = "The Azure name (uuid) of the created role assignment."
  value       = azapi_resource.this.name
}

output "role_assignment_resource_id" {
  description = "The Azure resource id of the created role assignment."
  value       = azapi_resource.this.id
}
