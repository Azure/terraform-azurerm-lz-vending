variable "location" {
  type        = string
  description = "The Azure region to deploy resources into. E.g. `eastus`"
  nullable    = false
}

variable "resource_group_name" {
  type        = string
  description = "The name of the resource group E.g. `rg-test`"
  nullable    = false

  validation {
    condition     = trimspace(var.resource_group_name) != ""
    error_message = "The resource_group_name must not be empty."
  }
}

variable "subscription_id" {
  type        = string
  description = "The ID of the subscription to deploy resources into. E.g. `00000000-0000-0000-0000-000000000000`"

  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must a subscription id in the format of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "lock_enabled" {
  type        = bool
  default     = false
  description = "Whether to enable resource group lock for the resource group"
  nullable    = false
}

variable "lock_name" {
  type        = string
  default     = null
  description = "The name of the resource group lock for the resource group, if `null` will be set to `lock-<resource_group_name>`"

}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Map of tags to be applied to the resource group"
}
