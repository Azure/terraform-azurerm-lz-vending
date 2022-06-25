variable "subscription_alias_enabled" {
  default     = false
  type        = bool
  description = <<DESCRIPTION
Whether the creation of a new subscripion alias is enabled or not.

If it is disabled, the `subscription_id` variable must be supplied instead.
  DESCRIPTION
}

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
}
