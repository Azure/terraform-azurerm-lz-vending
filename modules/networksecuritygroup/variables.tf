# This is required for most resource modules
variable "parent_id" {
  type        = string
  description = "The ID of the parent resource to which this user-assigned managed identity."

  validation {
    condition     = length(var.parent_id) > 0
    error_message = "The parent_id must not be empty."
  }
  validation {
    condition     = can(regex("^/subscriptions/[a-fA-F0-9-]+/resourceGroups/[a-zA-Z0-9-_.()]+$", var.parent_id))
    error_message = "The parent_id must be a valid Azure Resource Group ID."
  }
}

variable "location" {
  type        = string
  description = "(Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created."
  nullable    = false
}

variable "name" {
  type        = string
  description = "(Required) Specifies the name of the network security group. Changing this forces a new resource to be created."
  nullable    = false

  validation {
    condition     = can(regex("^[[:alnum:]]([[:alnum:]_.-]{0,78}?[[:alnum:]_])?$", var.name))
    error_message = <<EOT
    The name must be between 1 and 80 characters long and can only contain alphanumerics, underscores, periods, and hyphens.
    It must start with an alphanumeric and end with an alphanumeric or underscore.
    EOT
  }
}

# Optional variables
variable "security_rules" {
  type = map(object({
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
  }))
  description = <<DESCRIPTION
These are the security rule configuration properties.
- `access` - (Required) Specifies whether network traffic is allowed or denied. Possible values are `Allow` and `Deny`.
- `description` - (Optional) A description for this rule. Restricted to 140 characters.
- `destination_address_prefix` - (Optional) CIDR or destination IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. Besides, it also supports all available Service Tags like ‘Sql.WestEurope‘, ‘Storage.EastUS‘, etc.
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
  default     = {}
  nullable    = false
}

variable "tags" {
  type        = map(string)
  default     = null
  description = "(Optional) A mapping of tags to assign to the resource."
}
