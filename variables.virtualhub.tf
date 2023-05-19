variable "intent_based_routing_enabled" {
  description = "Enables and disables the virtual network submodule."
  type        = bool
  default     = false
}

variable "virtual_hubs" {
  type = map(object({
    virtual_hub_id                = string
    intent_based_routing_next_hop_firewall = string
  }))
  default     = {}
}
