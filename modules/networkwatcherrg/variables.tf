variable "location" {
  type        = string
  description = "The Azure region to deploy resources into."
}

variable "subscription_id" {
  type        = string
  description = "The ID of the subscription to deploy resources into. E.g. `00000000-0000-0000-0000-000000000000`"
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a subscription id in the format of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "tags" {
  type        = map(string)
  description = "A mapping of tags to assign to the resource."
  default     = {}
}

variable "network_watcher_rg_name" {
  type        = string
  description = "The name of the resource group in which to create the network watcher. This only needs changing for parallel testing purposes."
  default     = "NetworkWatcherRG"
}
