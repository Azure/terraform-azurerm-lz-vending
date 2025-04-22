resource "azapi_resource" "network_security_group" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  name      = var.name
  parent_id = var.resource_group_resource_id
  location  = var.location
  body = {
    properties = {
      securityRules = [
        for rule in var.security_rules : {
          name = rule.name
          properties = {
            access                               = rule.access
            description                          = rule.description
            destinationAddressPrefix             = rule.destination_address_prefix
            destinationAddressPrefixes           = rule.destination_address_prefixes
            destinationApplicationSecurityGroups = rule.destination_application_security_group_ids != null ? [for asg in rule.destination_application_security_group_ids : { id = asg }] : null
            destinationPortRange                 = rule.destination_port_range
            destinationPortRanges                = rule.destination_port_ranges
            direction                            = rule.direction
            priority                             = rule.priority
            protocol                             = rule.protocol
            sourceAddressPrefix                  = rule.source_address_prefix
            sourceAddressPrefixes                = rule.source_address_prefixes
            sourceApplicationSecurityGroups      = rule.source_application_security_group_ids != null ? [for asg in rule.source_application_security_group_ids : { id = asg }] : null
            sourcePortRange                      = rule.source_port_range
            sourcePortRanges                     = rule.source_port_ranges
          }
        }
      ]
    }
  }
}
