variable "resource_provider" {
  type        = string
  description = <<DESCRIPTION
The resource provider namespace, e.g. `Microsoft.Compute`.
DESCRIPTION
  nullable    = false
}

variable "subscription_id" {
  type        = string
  description = <<DESCRIPTION
The subscription id to register the resource providers in.
DESCRIPTION
  nullable    = false

  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "subscription_id must be set"
  }
}

variable "features" {
  type        = set(string)
  default     = []
  description = <<DESCRIPTION
The resource provider features to register, e.g. [`MyFeature`]
DESCRIPTION
  nullable    = false
}
