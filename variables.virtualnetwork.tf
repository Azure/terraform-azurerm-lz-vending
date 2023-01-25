variable "virtual_network_enabled" {
  description = "Enables and disables the virtual network submodule."
  type        = bool
  default     = false
}

variable "virtual_networks" {
  type = map(object({
    name                = string
    address_space       = list(string)
    resource_group_name = string

    location = optional(string, "")

    dns_servers = optional(list(string), [])

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(bool, false)
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    # Reserved for future capability
    #
    # other_peerings = optional(map(object({
    #   remote_network_resource_id            = string
    #   name_inbound                          = optional(string, "")
    #   name_outbound                         = optional(string, "")
    #   outbound_only                         = optional(bool, false)
    #   allow_forwarded_traffic_inbound       = optional(bool, true)
    #   allow_forwarded_traffic_outbound      = optional(bool, true)
    #   allow_gateway_transit_inbound         = optional(bool, false)
    #   allow_gateway_transit_outbound        = optional(bool, false)
    #   allow_virtual_network_access_inbound  = optional(bool, true)
    #   allow_virtual_network_access_outbound = optional(bool, true)
    #   use_remote_gateways_inbound           = optional(bool, false)
    #   use_remote_gateways_outbound          = optional(bool, false)
    # })), {})

    resource_group_creation_enabled = optional(bool, true)
    resource_group_lock_enabled     = optional(bool, true)
    resource_group_lock_name        = optional(string, "")
    resource_group_tags             = optional(map(string), {})

    vwan_associated_routetable_resource_id   = optional(string, "")
    vwan_connection_enabled                  = optional(bool, false)
    vwan_connection_name                     = optional(string, "")
    vwan_hub_resource_id                     = optional(string, "")
    vwan_propagated_routetables_labels       = optional(list(string), [])
    vwan_propagated_routetables_resource_ids = optional(list(string), [])
    vwan_security_configuration = object({
      secure_internet_traffic = optional(bool, false)
      secure_private_traffic  = optional(bool, false)
      next_hop                = optional(string, "")
    })

    tags = optional(map(string), {})
  }))
  description = <<DESCRIPTION
A map of the virtual networks to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the virtual network. [required]
- `address_space`: The address space of the virtual network as a list of strings in CIDR format, e.g. `["192.168.0.0/24", "10.0.0.0/24"]`. [required]
- `resource_group_name`: The name of the resource group to create the virtual network in. [required]

### DNS servers

- `dns_servers`: A list of DNS servers to use for the virtual network, e.g. `["192.168.0.1", "10.0.0.1"]`. If empty will use the Azure default DNS. [optional - default empty list]

### Location

- `location`: The location of the virtual network (and resource group if creation is enabled). [optional, will use `var.location` if not specified or empty string]

> Note at least one of `location` or `var.location` must be specified.
> If both are empty then the module will fail.

### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network.

- `hub_peering_enabled`: Whether to enable hub peering. [optional]
- `hub_network_resource_id`: The resource ID of the hub network to peer with. [optional - but required if hub_peering_enabled is `true`]
- `hub_peering_name_tohub`: The name of the peering to the hub network. [optional - leave empty to use calculated name]
- `hub_peering_name_fromhub`: The name of the peering from the hub network. [optional - leave empty to use calculated name]
- `hub_peering_use_remote_gateways`: Whether to use remote gateways for the hub peering. [optional - default true]

### Mesh peering values

Mesh peering is the capability to create a bi-directional peerings between all supplied virtual networks in `var.virtual_networks`.
Peerings will only be created between virtual networks with the `mesh_peering_enabled` value set to `true`.

- `mesh_peering_enabled`: Whether to enable mesh peering for this virtual network. Must be enabled on more than one virtual network for any peerings to be created. [optional]
- `mesh_peering_allow_forwarded_traffic`: Whether to allow forwarded traffic for the mesh peering. [optional - default false]

### Resource group values

The default is that a resource group will be created for each resource_group_name specified in the `var.virtual_networks` map.
It is possible to use a pre-existing resource group by setting `resource_group_creation_enabled` to `false`.
We recommend using resource groups aligned to the region of the virtual network,
however if you want multiple virtual networks in more than one location to share a resource group,
only one of the virtual networks should have `resource_group_creation_enabled` set to `true`.

- `resource_group_creation_enabled`: Whether to create a resource group for the virtual network. [optional - default `true`]
- `resource_group_lock_enabled`: Whether to create a `CanNotDelete` resource lock on the resource group. [optional - default `true`]
- `resource_group_lock_name`: The name of the resource lock. [optional - leave empty to use calculated name]
- `resource_group_tags`: A map of tags to apply to the resource group, e.g. `{ mytag = "myvalue", mytag2 = "myvalue2" }`. [optional - default empty]

### Virtual WAN values

- `vwan_associated_routetable_resource_id`: The resource ID of the route table to associate with the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_connection_enabled`: Whether to create a connection to a Virtual WAN. [optional - default false]
- `vwan_connection_name`: The name of the connection to the Virtual WAN. [optional - leave empty to use calculated name]
- `vwan_hub_resource_id`: The resource ID of the hub to connect to. [optional - but required if vwan_connection_enabled is `true`]
- `vwan_propagated_routetables_labels`: A list of labels of route tables to propagate to the virtual network. [optional - leave empty to use `["default"]`]
- `vwan_propagated_routetables_resource_ids`: A list of resource IDs of route tables to propagate to the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]
DESCRIPTION
  default     = {}
}
