variable "subscription_alias_enabled" {
  type = bool
}

variable "subscription_billing_scope" {
  type = string
}

variable "subscription_display_name" {
  type = string
}

variable "subscription_alias_name" {
  type = string
}

variable "subscription_workload" {
  type = string
}

variable "virtual_network_enabled" {
  type = string
}

variable "virtual_network_address_space" {
  type = list(string)
}

variable "virtual_network_resource_group_name" {
  type = string
}

variable "virtual_network_peering_enabled" {
  type = bool
}

variable "virtual_network_use_remote_gateways" {
  type = bool
}

variable "virtual_network_name" {
  type = string
}

variable "location" {
  type = string
}

variable "role_assignment_enabled" {
  type = bool
}
