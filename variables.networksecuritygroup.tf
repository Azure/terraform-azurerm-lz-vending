variable "network_security_group_enabled" {
  type        = bool
  description = "Whether to create network security groups and security rules in the target subscription. Requires `var.network_security_groups`."
  default     = false
}

variable "network_security_groups" {
  type = map(object({
    name                       = string
    location                   = optional(string)
    resource_group_resource_id = string
    tags                       = optional(map(string))

    security_rules = optional(map(object({
      access                                     = string
      description                                = optional(string)
      destination_address_prefix                 = optional(string)
      destination_address_prefixes               = optional(set(string))
      destination_application_security_group_ids = optional(set(string))
      destination_port_range                     = optional(string)
      destination_port_ranges                    = optional(set(string))
      direction                                  = string
      name                                       = string
      priority                                   = number
      protocol                                   = string
      source_address_prefix                      = optional(string)
      source_address_prefixes                    = optional(set(string))
      source_application_security_group_ids      = optional(set(string))
      source_port_range                          = optional(string)
      source_port_ranges                         = optional(set(string))
    })))
  }))
  description = <<DESCRIPTION
A map of the network security groups to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the network security group. Changing this forces a new resource to be created. [required]
- `resource_group_resource_id`: The resource id of the resource group to create the network security group in. Moving forward, the modules within this accelerator will adopt the standard of requiring the input be a resource id rather than a resource group name. Changing this forces a new resource to be created. [required]

### Location

- `location`: The supported Azure location where the resource exists. Changing this forces a new resource to be created.

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]


### Security Rules

- `security_rules` - (Optional) A map of security rules to create within the network network security group. The value is an object with the following fields: 
  - `access` - (Required) Specifies whether network traffic is allowed or denied. Possible values are `Allow` and `Deny`.
  - `description` - (Optional) A description for this rule. Restricted to 140 characters.
  - `destination_address_prefix` - (Optional) CIDR or destination IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. Besides, it also supports all available Service Tags like ‘Sql.WestEurope‘, ‘Storage.EastUS‘, etc. You can list the available service tags with the CLI: ```shell az network list-service-tags --location westcentralus```. For further information please see [Azure CLI
  - `destination_address_prefixes` - (Optional) List of destination address prefixes. Tags may not be used. This is required if `destination_address_prefix` is not specified.
  - `destination_application_security_group_ids` - (Optional) A List of destination Application Security Group IDs
  - `destination_port_range` - (Optional) Destination Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `destination_port_ranges` is not specified.
  - `destination_port_ranges` - (Optional) List of destination ports or port ranges. This is required if `destination_port_range` is not specified.
  - `direction` - (Required) The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are `Inbound` and `Outbound`.
  - `name` - (Required) The name of the security rule. This needs to be unique across all Rules in the Network Security Group. Changing this forces a new resource to be created.
  - `priority` - (Required) Specifies the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.
  - `protocol` - (Required) Network protocol this rule applies to. Possible values include `Tcp`, `Udp`, `Icmp`, `Esp`, `Ah` or `*` (which matches all).
  - `source_address_prefix` - (Optional) CIDR or source IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. This is required if `source_address_prefixes` is not specified.
  - `source_address_prefixes` - (Optional) List of source address prefixes. Tags may not be used. This is required if `source_address_prefix` is not specified.
  - `source_application_security_group_ids` - (Optional) A List of source Application Security Group IDs
  - `source_port_range` - (Optional) Source Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `source_port_ranges` is not specified.
  - `source_port_ranges` - (Optional) List of source ports or port ranges. This is required if `source_port_range` is not specified.

DESCRIPTION
  nullable    = false
  default     = {}
}
