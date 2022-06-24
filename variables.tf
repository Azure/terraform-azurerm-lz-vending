variable "location" {
  type = string
  description = <<DESCRIPTION
    The location of resources deployed by this module.
  DESCRIPTION
  default = ""
}

variable "subscription_id" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  An existing subscription id.

  Use this when you do not want the module to create a new subscription.

  A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
  All letters must be lowercase.

  You may also supply an empty string if you want to create a new subscription alias.
  In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:

  - `subscription_alias_name`
  - `subscription_alias_display_name`
  - `subscription_alias_billing_scope`
  - `subscription_alias_workload`
  DESCRIPTION
  validation {
    condition     = can(regex("^$|^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must be empty, or a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}
