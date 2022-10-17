output "virtual_network_resource_id" {
  description = "The created virtual network resource ID"
  value       = {
    for k, v in azapi_resource.vnet : k => v.id
  }
}
