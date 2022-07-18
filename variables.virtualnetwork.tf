# Note:
# Variable validation rules are disabled in the parent module and defaults supplied
# to support the case when the user does not want the virtual network to be deployed.

variable "virtual_network_enabled" {
  type        = bool
  description = <<DESCRIPTION
Enables and disables the virtual network submodule.
DESCRIPTION
  default     = false
}

variable "virtual_network_name" {
  type        = string
  description = <<DESCRIPTION
The name of the virtual network.
DESCRIPTION
  default     = ""
}

variable "virtual_network_address_space" {
  type        = list(string)
  description = <<DESCRIPTION
The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/8","172.16.0.0/12"]`.
DESCRIPTION
  default     = []
}

variable "virtual_network_peering_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether to enable peering with the supplied hub virtual network.
Enables a hub & spoke networking topology.

If enabled the `hub_network_resource_id` must also be suppled.
DESCRIPTION
  default     = false
}

variable "hub_network_resource_id" {
  type        = string
  description = <<DESCRIPTION
The resource ID of the virtual network in the hub to which the created virtual network will be peered.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`

Leave blank to create the virtual network without peering.
DESCRIPTION
  default     = ""
}

variable "virtual_network_vwan_connection_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether to enable connection with supplied vwan hub.
Enables a vwan networking topology.

If enabled the `vwan_hub_resource_id` must also be supplied.
DESCRIPTION
  default     = false
}

variable "vwan_hub_resource_id" {
  type        = string
  description = <<DESCRIPTION
The resource ID of the vwan hub to which the virtual network will be connected.
E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`

Leave blank to create a virtual network without a vwan hub connection.
DESCRIPTION
  default     = ""
}

variable "virtual_network_resource_group_name" {
  type        = string
  description = <<DESCRIPTION
The name of the resource group to create the virtual network in.
DESCRIPTION
  default     = ""
}

variable "virtual_network_location" {
  type        = string
  description = <<DESCRIPTION
The location of the virtual network.

Use this to override the default location defined by `var.location`.
Leave blank to use the default location.
DESCRIPTION
  default     = ""
}

variable "virtual_network_use_remote_gateways" {
  type        = bool
  description = <<DESCRIPTION
Enables the use of remote gateways for the virtual network.

Applies to hub and spoke (vnet peerings).
DESCRIPTION
  default     = true
}

variable "virtual_network_vwan_routetable_resource_id" {
  type        = string
  description = <<DESCRIPTION
The resource ID of the virtual network route table to use for the virtual network.

Leave blank to use the `defaultRouteTable`.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable`
DESCRIPTION
  default     = ""
}

variable "virtual_network_vwan_propagated_routetables_resource_ids" {
  type        = list(string)
  description = <<DESCRIPTION
The list of route table resource ids to advertise routes to.

Leave blank to use the `defaultRouteTable`.
DESCRIPTION
  default     = []
}

variable "virtual_network_vwan_propagated_routetables_labels" {
  type        = list(string)
  description = <<DESCRIPTION
The list of virtual WAN labels to advertise the routes to.

Leave blank to use the `default` label.
DESCRIPTION
  default     = []
}

variable "virtual_network_resource_lock_enabled" {
  type        = bool
  description = <<DESCRIPTION
Enables the deployment of resource locks to the virtual network's resource group.
Currently only `CanNotDelete` locks are supported.
DESCRIPTION
  default     = true
}
