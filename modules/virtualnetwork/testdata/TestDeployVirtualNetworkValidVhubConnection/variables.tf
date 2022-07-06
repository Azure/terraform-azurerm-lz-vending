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

variable "virtual_network_vwan_connection_enabled" {
  type = bool
}
