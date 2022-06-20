variable "virtual_network_name" {
  type = string
  description = <<DESCRIPTION
    The name of the virtual network.
  DESCRIPTION
  validation {
    condition = can(regex("^[\\w-_.]{2,64}$", var.virtual_network_name))
    error_message = "The string must consist of a-z, A-Z, 0-9, -, _, and . (period) and be between 2 and 64 characters in length."
  }
}

variable "virtual_network_address_space" {
  type = list(string)
  description = <<DESCRIPTION
    The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/16","172.16.0.0/12"]`.
  DESCRIPTION
}

variable "hub_network_resource_id" {
  type = string
  description = <<DESCRIPTION
    The resource ID of the virtual network in the hub.

    E.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet
  DESCRIPTION
  default = ""
  validation {
    condition = can(regex("^$|^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}/resourceGroups/[\\w]{1,89}[^\\s.]/providers/Microsoft.Network/virtualNetworks/[\\w-_.]{2,64}$", var.hub_network_resource_id))
    error_message = "Value must be an Azure virtual network resource id, e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet"
  }
}

variable "subscription_id" {
  type = string
  description = <<DESCRIPTION
    The subscription ID of the subscription to create the virtual network in.
  DESCRIPTION
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "virtual_network_resource_group_name" {
  type = string
  description = <<DESCRIPTION
    The name of the resource group to create the virtual network in.
  DESCRIPTION
  validation {
    condition = can(regex("^[\\w]{1,89}[^\\s.]$", var.virtual_network_resource_group_name))
    error_message = "Value must be between 1 and 90 characters in length and start with a letter or number, and end with a letter or number."
  }
}

variable "virtual_network_location" {
  type = string
  description = <<DESCRIPTION
    The location of the virtual network.
  DESCRIPTION
}

# variable "virtual_network_subnets" {
#   type = map(object({
#     address_prefix = string
#     }))
#   description = <<DESCRIPTION
#     The subnets of the virtual network, supplied as multiple objects.

#     e.g.

#     ```terraform
#     virtual_network_subnets = {
#       subnet0 = {
#         address_prefix = "10.0.0.0/24"
#       },
#       subnet1 = {
#         address_prefix = "10.0.1.0/24"
#     } }
#     ```
#   DESCRIPTION
# }
