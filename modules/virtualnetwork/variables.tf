variable "subscription_id" {
  type        = string
  description = <<DESCRIPTION
The subscription ID of the subscription to create the virtual network in.
DESCRIPTION
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "location" {
  type        = string
  description = <<DESCRIPTION
The default location of resources created by this module.
Virtual networks will be created in this location unless overridden by the `location` attribute.
DESCRIPTION
  default     = ""
}

variable "virtual_networks" {
  type = map(object({
    name                = string
    address_space       = list(string)
    resource_group_name = string

    location = optional(string, "")

    dns_servers             = optional(list(string), [])
    flow_timeout_in_minutes = optional(number, null)

    ddos_protection_enabled = optional(bool, false)
    ddos_protection_plan_id = optional(string, null)

    subnets = optional(map(object(
      {
        name             = string
        address_prefixes = list(string)
        nat_gateway = optional(object({
          id = string
        }))
        network_security_group = optional(object({
          id = string
        }))
        private_endpoint_network_policies             = optional(string, "Enabled")
        private_link_service_network_policies_enabled = optional(bool, true)
        route_table = optional(object({
          id = optional(string)
        }))
        default_outbound_access_enabled = optional(bool, false)
        service_endpoints               = optional(set(string))
        service_endpoint_policies = optional(map(object({
          id = string
        })))
        delegation = optional(list(
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

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(bool, false)
    hub_peering_direction           = optional(string, "both")
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

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
- `resource_group_name`: The name of the resource group to create the virtual network in. The default is that the resource group will be created by this module. [required]

### DNS servers

- `dns_servers`: A list of DNS servers to use for the virtual network, e.g. `["192.168.0.1", "10.0.0.1]`. If empty will use the Azure default DNS. [optional - default empty list]
DNS. [optional - default empty list]

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
  - `private_endpoint_network_policies_enabled` - (Optional) Enable or Disable network policies for the private endpoint on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true.
  - `private_link_service_network_policies_enabled` - (Optional) Enable or Disable network policies for the private link service on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true.
  - `route_table` - (Optional) An object with the following fields which are mutually exclusive, choose either an external route table or the generated route table:
    - `id` - The ID of the Route Table which should be associated with the Subnet. Changing this forces a new association to be created.
  - `default_outbound_access_enabled` - (Optional) Whether to allow internet access from the subnet. Defaults to `false`.
  - `service_endpoints` - (Optional) The list of Service endpoints to associate with the subnet.
  - `service_endpoint_policies` - (Optional) The list of Service Endpoint Policy objects with the resource id to associate with the subnet.
    - `id` - The ID of the endpoint policy that should be associated with the subnet.
  - `service_endpoint_policy_assignment_enabled` - (Optional) Should the Service Endpoint Policy be assigned to the subnet? Default `true`.
  - `delegation` - (Optional) An object with the following fields:
    - `name` - The name of the delegation.
    - `service_delegation` - An object with the following fields:
      - `name` - The name of the service delegation.

### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network.

- `hub_peering_enabled`: Whether to enable hub peering. [optional - default `false`]
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
- `vwan_security_configuration`: A map of security configuration values for VWAN hub connection - see below. [optional - default empty]
  - `secure_internet_traffic`: Whether to forward internet-bound traffic to the destination specified in the routing policy. [optional - default `false`]
  - `secure_private_traffic`: Whether to all internal traffic to the destination specified in the routing policy. Not compatible with `routing_intent_enabled`. [optional - default `false`]
  - `routing_intent_enabled`: Enable to use with a Virtual WAN hub with routing intent enabled. Routing intent on hub is configured outside this module. [optional - default `false`]

### Tags

- `tags`: A map of tags to apply to the virtual network. [optional - default empty]
DESCRIPTION

  # validate virtual_networks is no zero length
  validation {
    condition     = length(var.virtual_networks) > 0
    error_message = "The virtual_networks variable must not be empty."
  }

  # validate virtual network name
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^[\\w-_.]{2,64}$", v.name))
    ])
    error_message = "Virtual network name must consist of a-z, A-Z, 0-9, -, _, and . (period) and be between 2 and 64 characters in length."
  }

  # validate address space is not zero length
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      length(v.address_space) > 0
    ])
    error_message = "At least 1 address space must be specified."
  }

  # validate address space CIDR blocks are valid
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for cidr in v.address_space :
        can(regex("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\\/(3[0-2]|[1-2][0-9]|[0-9]))$", cidr))
      ]
    ]))
    error_message = "Address space entries must be specified in CIDR notation, e.g. 192.168.0.0/24."
  }

  # validate virtual network subnet names
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for subnet in v.subnets :
        can(regex("^[\\w-_.]{2,64}$", v.name))
      ]
    ]))
    error_message = "Virtual network subnet name must consist of a-z, A-Z, 0-9, -, _, and . (period) and be between 2 and 64 characters in length."
  }

  # validate subnet address prefixes is not zero length
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for subnet in v.subnets :
        length(subnet.address_prefixes) > 0
      ]
    ]))
    error_message = "At least 1 subnet address prefix must be specified."
  }

  # validate subnet nat gateway id is valid
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for subnet in v.subnets :
        can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/natGateways/[\\w-_.]{2,64}$", subnet.nat_gateway.id)) if try(subnet["nat_gateway"], null) != null
      ]
    ]))
    error_message = "Nat Gateway resource id must be valid, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/natGateways/testvnatgw."
  }

  # validate subnet network security group id is valid
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for subnet in v.subnets :
        can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/networkSecurityGroups/[\\w-_.]{2,64}$", subnet.network_security_group.id)) if try(subnet["network_security_group"], null) != null
      ]
    ]))
    error_message = "Network security group resource id must be valid, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/networkSecurityGroups/testvnsg."
  }

  # validate ddos protection plan resource id for networks with ddos protection enabled
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/ddosProtectionPlans/[\\w-_.]{2,64}$", v.ddos_protection_plan_id)) if v.ddos_protection_enabled
    ])
    error_message = "Hub network resource id must be an Azure ddos protection plan resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/ddosProtectionPlans/my-protection-plan."
  }

  # validate hub network resource id for networks with hub peering enabled
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", v.hub_network_resource_id)) if v.hub_peering_enabled
    ])
    error_message = "Hub network resource id must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet."
  }

  # validate vwan hub resource id for networks with vwan connection enabled
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}$", v.vwan_hub_resource_id)) if v.vwan_connection_enabled
    ])
    error_message = "The vWAN hub resource id must be an Azure vWAN hub network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub."
  }

  # validate vwan associated routetable resource id for networks with vwan connection enabled
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}/hubRouteTables/[\\w-_.]{1,80}$", v.vwan_associated_routetable_resource_id)) if v.vwan_connection_enabled
    ])
    error_message = "The vWAN associated routetable resource id must be an Azure vwan hub routetable resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable."
  }

  # validate vwan propagated routetable resource ids for networks with vwan connection enabled
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for i in v.vwan_propagated_routetables_resource_ids :
        can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}/hubRouteTables/[\\w-_.]{1,80}$", i))
      ] if v.vwan_connection_enabled
    ]))
    error_message = "The vWAN propagated routetables resource id must be an Azure vwan hub routetable resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable."
  }

  # Reserved for future functionality
  #
  # # validate other peering network resource id
  # validation {
  #   condition = alltrue(flatten([
  #     for k, v in var.virtual_networks :
  #     [
  #       for k2, v2 in v.other_peerings :
  #       can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", v2.remote_network_resource_id))
  #     ]
  #   ]))
  #   error_message = "Other peering remote network resource id must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet."
  # }

  # validate resource groups with creation enabled have unique names.
  validation {
    condition = can(
      {
        for i in toset([
          for k, v in var.virtual_networks : {
            name     = v.resource_group_name
            location = v.location
          } if v.resource_group_creation_enabled
        ]) : i.name => i.location
      }
    )
    error_message = "Resource group names with creation enabled must be unique. Virtual networks deployed into the same resource group must have only one enabled for resource group creation."
  }
}

variable "enable_telemetry" {
  type        = bool
  default     = true
  description = <<DESCRIPTION
This variable controls whether or not telemetry is enabled for the module.
For more information see https://aka.ms/avm/telemetryinfo.
If it is set to false, then no telemetry will be collected.
DESCRIPTION
  nullable    = false
}
