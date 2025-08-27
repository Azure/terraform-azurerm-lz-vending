variable "virtual_network_enabled" {
  description = "Enables and disables the virtual network submodule."
  type        = bool
  default     = false
}

variable "virtual_networks" {
  type = map(object({
    name                         = string
    address_space                = list(string)
    resource_group_key           = optional(string)
    resource_group_name_existing = optional(string)

    location = optional(string)

    dns_servers             = optional(list(string), [])
    flow_timeout_in_minutes = optional(number)

    ddos_protection_enabled = optional(bool, false)
    ddos_protection_plan_id = optional(string)

    subnets = optional(map(object(
      {
        name             = string
        address_prefixes = list(string)
        nat_gateway = optional(object({
          id = string
        }))
        network_security_group = optional(object({
          id            = optional(string)
          key_reference = optional(string)
        }))
        private_endpoint_network_policies             = optional(string, "Enabled")
        private_link_service_network_policies_enabled = optional(bool, true)
        route_table = optional(object({
          id            = optional(string)
          key_reference = optional(string)
        }))
        default_outbound_access_enabled = optional(bool, false)
        service_endpoints               = optional(set(string))
        service_endpoint_policies = optional(map(object({
          id = string
        })))
        delegations = optional(list(
          object(
            {
              name = string
              service_delegation = object({
                name = string
              })
            }
          )
        ))
      }
    )), {})

    hub_network_resource_id = optional(string)
    hub_peering_enabled     = optional(bool, false)
    hub_peering_direction   = optional(string, "both")
    hub_peering_name_tohub  = optional(string)
    hub_peering_options_tohub = optional(object({
      allow_forwarded_traffic       = optional(bool, true)
      allow_gateway_transit         = optional(bool, false)
      allow_virtual_network_access  = optional(bool, true)
      do_not_verify_remote_gateways = optional(bool, false)
      enable_only_ipv6_peering      = optional(bool, false)
      local_peered_address_spaces   = optional(list(string), [])
      local_peered_subnets          = optional(list(string), [])
      peer_complete_vnets           = optional(bool, true)
      remote_peered_address_spaces  = optional(list(string), [])
      remote_peered_subnets         = optional(list(string), [])
      use_remote_gateways           = optional(bool, true)
    }), {})
    hub_peering_name_fromhub = optional(string)
    hub_peering_options_fromhub = optional(object({
      allow_forwarded_traffic       = optional(bool, true)
      allow_gateway_transit         = optional(bool, true)
      allow_virtual_network_access  = optional(bool, true)
      do_not_verify_remote_gateways = optional(bool, false)
      enable_only_ipv6_peering      = optional(bool, false)
      local_peered_address_spaces   = optional(list(string), [])
      local_peered_subnets          = optional(list(string), [])
      peer_complete_vnets           = optional(bool, true)
      remote_peered_address_spaces  = optional(list(string), [])
      remote_peered_subnets         = optional(list(string), [])
      use_remote_gateways           = optional(bool, false)
    }), {})

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    vwan_associated_routetable_resource_id   = optional(string)
    vwan_connection_enabled                  = optional(bool, false)
    vwan_connection_name                     = optional(string)
    vwan_hub_resource_id                     = optional(string)
    vwan_propagated_routetables_labels       = optional(list(string), [])
    vwan_propagated_routetables_resource_ids = optional(list(string), [])
    vwan_security_configuration = optional(object({
      secure_internet_traffic = optional(bool, false)
      secure_private_traffic  = optional(bool, false)
      routing_intent_enabled  = optional(bool, false)
    }), {})

    tags = optional(map(string), {})
  }))
  description = <<DESCRIPTION
A map of the virtual networks to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the virtual network. [required]
- `address_space`: The address space of the virtual network as a list of strings in CIDR format, e.g. `["192.168.0.0/24", "10.0.0.0/24"]`. [required]
- `resource_group_key`: The resource group key from the resource groups map to create the virtual network in. [optional]
- `resource_group_name_existing`: The name of an existing resource group to use for the virtual network. [optional]

**One of `resource_group_key` or `resource_group_name_existing` must be specified.**

### DNS servers

- `dns_servers`: A list of DNS servers to use for the virtual network, e.g. `["192.168.0.1", "10.0.0.1"]`. If empty will use the Azure default DNS. [optional - default empty list]

### DDOS protection plan

- `ddos_protection_enabled`: Whether to enable ddos protection. [optional]
- `ddos_protection_plan_id`: The resource ID of the protection plan to attach the vnet. [optional - but required if ddos_protection_enabled is `true`]

### Location

- `location`: The location of the virtual network (and resource group if creation is enabled). [optional, will use `var.location` if not specified or empty string]

> Note at least one of `location` or `var.location` must be specified.
> If both are empty then the module will fail.

#### Subnets

- `subnets` - (Optional) A map of subnets to create in the virtual network. The value is an object with the following fields:
  - `name` - The name of the subnet.
  - `address_prefixes` - The IPv4 address prefixes to use for the subnet in CIDR format.
  - `nat_gateway` - (Optional) An object with the following fields:
    - `id` - The ID of the NAT Gateway which should be associated with the Subnet. Changing this forces a new resource to be created.
  - `network_security_group` - (Optional) An object with the following fields:
    - `id` - The ID of the Network Security Group which should be associated with the Subnet. Changing this forces a new association to be created.
    - `key_reference` - The name of the var.network_security_group map key that should be associated with the subnet once it has been provisioned. If you are passing in an `id` value, this will not be used.
  - `private_endpoint_network_policies_enabled` - (Optional) Enable or Disable network policies for the private endpoint on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true.
  - `private_link_service_network_policies_enabled` - (Optional) Enable or Disable network policies for the private link service on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true.
  - `route_table` - (Optional) An object with the following fields which are mutually exclusive, choose either an external route table or the generated route table:
    - `id` - The ID of the Route Table which should be associated with the Subnet. Changing this forces a new association to be created.
    - `key_reference` - The name of the var.route_tables map key that should be associated with the subnet once it has been provisioned. If you are passing in an `id` value, this will not be used.
  - `default_outbound_access_enabled` - (Optional) Whether to allow internet access from the subnet. Defaults to `false`.
  - `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet.
  - `service_endpoint_policies` - (Optional) The list of Service Endpoint Policy objects with the resource id to associate with the subnet.
    - `id` - The ID of the endpoint policy that should be associated with the subnet.
  - `service_endpoint_policy_assignment_enabled` - (Optional) Should the Service Endpoint Policy be assigned to the subnet? Default `true`.
  - `delegation` - (Optional) An object with the following fields:
    - `name` - The name of the delegation.
    - `service_delegation` - An object with the following fields:
      - `name` - The name of the service delegation.
      - `actions` - A list of actions that should be delegated, the list is specific to the service being delegated.


### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network:

- `hub_peering_enabled`: Whether to enable hub peering. [optional]
- `hub_peering_direction`: The direction of the peering. [optional - allowed values are: `tohub`, `fromhub` or `both` - default `both`]
- `hub_network_resource_id`: The resource ID of the hub network to peer with. [optional - but required if hub_peering_enabled is `true`]
- `hub_peering_name_tohub`: The name of the peering to the hub network. [optional - leave empty to use calculated name]
- `hub_peering_name_fromhub`: The name of the peering from the hub network. [optional - leave empty to use calculated name]

#### Hub network peering options

The following values configure the options for the hub network peering. These are configurable in each direction:

- `allow_forwarded_traffic`: Whether to allow forwarded traffic for the peering. [optional - default `true`]
- `allow_gateway_transit`: Whether to allow gateway transit for the peering. [optional - default `false` (outbound) or `true` (inbound)]
- `allow_virtual_network_access`: Whether to allow virtual network access for the peering. [optional - default `true`]
- `do_not_verify_remote_gateways`: Whether to not verify remote gateways for the peering. [optional - default `false`]
- `enable_only_ipv6_peering`: Whether to enable only IPv6 peering. [optional - default `false`]
- `local_peered_address_spaces`: A list of local address spaces to peer with. [optional - default empty and only used if `peer_complete_vnets` is `false`]
- `local_peered_subnets`: A list of local subnets to peer with. [optional - default empty and only used if `peer_complete_vnets` is `false`]
- `peer_complete_vnets`: Whether to peer complete virtual networks. [optional - default `true`]
- `remote_peered_address_spaces`: A list of remote address spaces to peer with. [optional - default empty and only used if `peer_complete_vnets` is `false`]
- `remote_peered_subnets`: A list of remote subnets to peer with. [optional - default empty and only used if `peer_complete_vnets` is `false`]
- `use_remote_gateways`: Whether to use remote gateways for the peering. [optional - default `true` (outbound) or `false` (inbound)]

### Mesh peering values

Mesh peering is the capability to create a bi-directional peerings between all supplied virtual networks in `var.virtual_networks`.
Peerings will only be created between virtual networks with the `mesh_peering_enabled` value set to `true`.

- `mesh_peering_enabled`: Whether to enable mesh peering for this virtual network. Must be enabled on more than one virtual network for any peerings to be created. [optional]
- `mesh_peering_allow_forwarded_traffic`: Whether to allow forwarded traffic for the mesh peering. [optional - default false]

### Virtual WAN values

- `vwan_associated_routetable_resource_id`: The resource ID of the route table to associate with the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_connection_enabled`: Whether to create a connection to a Virtual WAN. [optional - default false]
- `vwan_connection_name`: The name of the connection to the Virtual WAN. [optional - leave empty to use calculated name]
- `vwan_hub_resource_id`: The resource ID of the hub to connect to. [optional - but required if vwan_connection_enabled is `true`]
- `vwan_propagated_routetables_labels`: A list of labels of route tables to propagate to the virtual network. [optional - leave empty to use `["default"]`]
- `vwan_propagated_routetables_resource_ids`: A list of resource IDs of route tables to propagate to the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_security_configuration`: A map of security configuration values for VWAN hub connection - see below. [optional - default empty]
  - `secure_internet_traffic`: Whether to forward internet-bound traffic to the destination specified in the routing policy. [optional - default `false`]
  - `secure_private_traffic`: Whether to all internal traffic to the destination specified in the routing policy. Not compatible with `routing_intent_enabled`. [optional - default `false`]
  - `routing_intent_enabled`: Enable to use with a Virtual WAN hub with routing intent enabled. Routing intent on hub is configured outside this module. [optional - default `false`]

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]
DESCRIPTION
  nullable    = false
  default     = {}

  validation {
    condition     = alltrue([for v in var.virtual_networks : try(coalesce(v.resource_group_key, v.resource_group_name_existing), null) != null])
    error_message = "Each virtual network must specify either 'resource_group_key' or 'resource_group_name_existing'."
  }
}
