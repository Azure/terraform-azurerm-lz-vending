# variable "subscription_id" {
#   type        = string
#   description = <<DESCRIPTION
# The subscription ID of the subscription to create the virtual network in.
# DESCRIPTION
#   validation {
#     condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
#     error_message = "Must a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
#   }
# }

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
