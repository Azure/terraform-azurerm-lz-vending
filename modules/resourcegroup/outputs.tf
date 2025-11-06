output "resource_group_name" {
  description = "The created resource group name."
  value       = azapi_resource.rg.name
}

output "resource_group_resource_id" {
  description = "The created resource group resource ID."
  value       = azapi_resource.rg.id
}
