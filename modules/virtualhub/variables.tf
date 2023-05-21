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
