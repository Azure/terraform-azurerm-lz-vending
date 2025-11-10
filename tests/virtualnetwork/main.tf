# Test wrapper for virtual network module

module "virtualnetwork" {
  source = "../../modules/virtualnetwork"

  subscription_id  = var.subscription_id
  location         = var.location
  enable_telemetry = var.enable_telemetry
  virtual_networks = var.virtual_networks
}

variable "subscription_id" {
  type = string
}

variable "location" {
  type    = string
  default = "westeurope"
}

variable "enable_telemetry" {
  type    = bool
  default = false
}

variable "virtual_networks" {
  type = map(object({
    name                                          = string
    address_space                                 = list(string)
    location                                      = optional(string)
    resource_group_name                           = string
    dns_servers                                   = optional(list(string), [])
    ddos_protection_enabled                       = optional(bool, false)
    ddos_protection_plan_id                       = optional(string)
    flow_timeout_in_minutes                       = optional(number, 4)
    hub_peering_enabled                           = optional(bool, false)
    hub_peering_direction                         = optional(string, "both")
    hub_peering_name_tohub                        = optional(string)
    hub_peering_name_fromhub                      = optional(string)
    hub_network_resource_id                       = optional(string)
    hub_peering_options_tohub                     = optional(map(bool), {})
    hub_peering_options_fromhub                   = optional(map(bool), {})
    mesh_peering_enabled                          = optional(bool, false)
    mesh_peering_allow_forwarded_traffic          = optional(bool, false)
    subnets                                       = optional(map(any), {})
    tags                                          = optional(map(string), {})
    vwan_hub_resource_id                          = optional(string)
    vwan_connection_enabled                       = optional(bool, false)
    vwan_associated_routetable_resource_id        = optional(string)
    vwan_propagated_routetables_labels            = optional(list(string), [])
    vwan_propagated_routetables_resource_ids      = optional(list(string), [])
    vwan_security_configuration                   = optional(map(any), {})
  }))
}

output "virtual_network_resource_ids" {
  value = module.virtualnetwork.virtual_network_resource_ids
}
