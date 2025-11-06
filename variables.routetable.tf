variable "route_table_enabled" {
  type        = bool
  default     = false
  description = "Whether to create route tables and routes in the target subscription. Requires `var.route_tables`."
}

variable "route_tables" {
  type = map(object({
    name                          = string
    location                      = string
    resource_group_key            = optional(string)
    resource_group_name_existing  = optional(string)
    bgp_route_propagation_enabled = optional(bool, true)
    tags                          = optional(map(string))

    routes = optional(map(object({
      name                   = string
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = optional(string)
    })), {})
  }))
  default     = {}
  description = <<DESCRIPTION
A map defining route tables and their associated routes to be created:

- `name` (required): The name of the route table.
- `location` (required): The location of the resource group.
- `resource_group_key`: The resource group key from the resource groups map to create the user assigned identity in. [optional]
- `resource_group_name_existing`: The name of an existing resource group to create the user assigned identity in. [optional]

**One of `resource_group_key` or `resource_group_name_existing` must be specified.**

- `bgp_route_propagation_enabled` (optional): Boolean that controls whether routes learned by BGP are propagated to the route table. Default is `true`.
- `tags` (optional): A map of key-value pairs for tags associated with the route table.
- `routes` (optional): A map defining routes for the route table. Each route object has the following properties:
- `name` (required): The name of the route.
- `address_prefix` (required): The address prefix for the route.
- `next_hop_type` (required): The next hop type, must be one of: 'Internet', 'None', 'VirtualAppliance', 'VirtualNetworkGateway', 'VnetLocal'.
- `next_hop_in_ip_address` (optional): The next hop IP address for the route. Required if next hop type is 'VirtualAppliance'.
DESCRIPTION
  nullable    = false
}
