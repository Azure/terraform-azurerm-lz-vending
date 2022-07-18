
variable "subscription_id" {
  type = string
}

variable "virtual_network_name" {
  type = string
}

variable "virtual_network_resource_group_name" {
  type = string
}

variable "virtual_network_location" {
  type = string
}

variable "virtual_network_address_space" {
  type = list(string)
}

variable "virtual_network_use_remote_gateways" {
  type = bool
}

variable "virtual_network_peering_enabled" {
  type = bool
}

variable "virtual_network_resource_lock_enabled" {
  type = bool
}
