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
    hub_peering_use_remote_gateways = optional(string, false)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    other_peerings = optional(map(object({
      remote_network_resource_id   = string
      peering_name                 = optional(string, "")
      outbound_only                = optional(bool, false)
      allow_virtual_network_access = optional(bool, true)
      allow_forwarded_traffic      = optional(bool, true)
      allow_gateway_transit        = optional(bool, false)
      use_remote_gateways          = optional(bool, false)
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
      for k, v in var.virtual_networks : can(regex("^[\\w-_.]{2,64}$", v.name))
    ])
    error_message = "The string must consist of a-z, A-Z, 0-9, -, _, and . (period) and be between 2 and 64 characters in length."
  }
  default = {}
}
