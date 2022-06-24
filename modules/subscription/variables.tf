variable "subscription_alias_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  The name of the subscription alias.

  The string must be comprised of a-z, A-Z, 0-9, - and _.
  The maximum length is 63 characters.

  You may also supply an empty string if you do not want to create a new subscription alias.
  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied.
  DESCRIPTION
  validation {
    condition     = can(regex("^$|^[a-zA-Z0-9-_]{1,63}$", var.subscription_alias_name))
    error_message = "Valid characters are a-z, A-Z, 0-9, -, _."
  }
}

variable "subscription_alias_display_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  The display name of the subscription alias.

  The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
  The maximum length is 63 characters.

  You may also supply an empty string if you do not want to create a new subscription alias.
  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied.
  DESCRIPTION
  validation {
    condition     = can(regex("^$|^[a-zA-Z0-9-_ ]{1,63}$", var.subscription_alias_display_name))
    error_message = "Valid characters are a-z, A-Z, 0-9, -, _, and space."
  }
}

variable "subscription_alias_billing_scope" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  The billing scope for the new subscription alias.

  A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

  You may also supply an empty string if you do not want to create a new subscription alias.
  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied.
  DESCRIPTION
  validation {
    condition     = can(regex("^$|^/providers/Microsoft.Billing/billingAccounts/.*$", var.subscription_alias_billing_scope))
    error_message = "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive."
  }
}

variable "subscription_alias_workload" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  The billing scope for the new subscription alias.

  The workload type can be either `Production` or `DevTest` and is case sensitive.

  You may also supply an empty string if you do not want to create a new subscription alias.
  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied.
  DESCRIPTION
  validation {
    condition     = can(regex("^$|^(Production|DevTest)$", var.subscription_alias_workload))
    error_message = "The workload type can be either Production or DevTest and is case sensitive."
  }
}

variable "subscription_alias_management_group_id" {
  type    = string
  default = ""
  description = <<DESCRIPTION
  The destination management group ID for the new subscription.

  **Note:** Do not supply the display name.
  The management group ID forms part of the Azure resource ID. E.g.,
  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`.
  DESCRIPTION
  validation {
    condition = can(regex("^$|^[().a-zA-Z0-9_-]{1,90}$", var.subscription_alias_management_group_id))
    error_message = "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.)."
  }
}

variable "subscription_id" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
  An existing subscription id.

  Use this when you do not want the nmodule to create a new subscription.

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
