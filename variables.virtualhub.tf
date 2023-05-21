variable "virtual_hub_enabled" {
  description = "Enables and disables the virtual hub submodule."
  type        = bool
  default     = false
}

variable "virtual_hubs" {
  type = map(object({
    vwan_hub_resource_id                  = optional(string, "")
    vhub_firewall_resource_id             = optional(string, "")
    intent_based_internet_traffic_enabled = optional(bool, false)
    intent_based_private_traffic_enabled  = optional(bool, false)
  }))
  default = {}
}
