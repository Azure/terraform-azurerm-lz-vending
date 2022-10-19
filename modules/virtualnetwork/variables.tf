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

variable "virtual_networks" {
  type = map(object({
    name                = string
    address_space       = list(string)
    location            = string
    resource_group_name = string

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(string, false)
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    other_peerings = optional(map(object({
      remote_network_resource_id            = string
      name_inbound                          = optional(string, "")
      name_outbound                         = optional(string, "")
      outbound_only                         = optional(bool, false)
      allow_forwarded_traffic_inbound       = optional(bool, true)
      allow_forwarded_traffic_outbound      = optional(bool, true)
      allow_gateway_transit_inbound         = optional(bool, false)
      allow_gateway_transit_outbound        = optional(bool, false)
      allow_virtual_network_access_inbound  = optional(bool, true)
      allow_virtual_network_access_outbound = optional(bool, true)
      use_remote_gateways_inbound           = optional(bool, false)
      use_remote_gateways_outbound          = optional(bool, false)
    })), {})

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

    tags = optional(map(string), {})
  }))
  description = <<DESCRIPTION
TODO
DESCRIPTION

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

  # validate hub network resource id
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", v.hub_network_resource_id))
    ])
    error_message = "Hub network resource id must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet."
  }

  # validate vwan hub resource id
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}$", v.vwan_hub_resource_id))
    ])
    error_message = "vWAN hub resource id must be an Azure vWAN hub network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub."
  }

  # validate vwan associated routetable resource id
  validation {
    condition = alltrue([
      for k, v in var.virtual_networks :
      can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}/hubRouteTables/[\\w-_.]{1,80}$", v.vwan_associated_routetable_resource_id))
    ])
    error_message = "vWAN associated routetable resource id must be an Azure vwan hub routetable resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable."
  }

  # validate vwan propagated routetable resource ids
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for i in v.vwan_propagated_routetables_resource_ids :
        can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w]{1,89}[^\\s.]/providers/Microsoft.Network/virtualHubs/[\\w-_.]{1,80}/hubRouteTables/[\\w-_.]{1,80}$", i))
      ]
    ]))
    error_message = "vWAN propagated routetables resource id must be an Azure vwan hub routetable resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable."
  }

  # validate other peering network resource id
  validation {
    condition = alltrue(flatten([
      for k, v in var.virtual_networks :
      [
        for k2, v2 in v.other_peerings :
        can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w-._]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", v2.remote_network_resource_id))
      ]
    ]))
    error_message = "Remote network resource id must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet."
  }

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

  default = {}
}
