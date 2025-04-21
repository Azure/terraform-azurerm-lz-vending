output "network_security_group_resource_id" {
  description = "The created network security group resource ID."
  value = {
    network_security_group = azapi_resource.network_security_group.id
  }
}
