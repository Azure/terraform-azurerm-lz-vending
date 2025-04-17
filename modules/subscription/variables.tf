variable "subscription_alias_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create a new subscription using the subscription alias resource.

If enabled, the following must also be supplied:

- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`

Optionally, supply the following to enable the placement of the subscription into a management group:

- `subscription_management_group_id`
- `subscription_management_group_association_enabled`

If disabled, supply the `subscription_id` variable to use an existing subscription instead.

> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.
DESCRIPTION
  nullable    = false
}

variable "subscription_alias_name" {
  type        = string
  default     = null
  description = <<DESCRIPTION
The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and _.
The maximum length is 63 characters.

You may also supply a null string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION

  validation {
    condition     = var.subscription_alias_name != null ? length(var.subscription_alias_name) <= 64 && !can(regex("[<>;|]", var.subscription_alias_name)) : true
    error_message = "Subscription Alias must either `null`, or be less or equal to 64 characters in length and cannot contain the characters `<`, `>`, `;`, or `|`"
  }
  validation {
    error_message = "Value must not be null if `subscription_alias_enabled` is set to `true`."
    condition     = var.subscription_alias_enabled ? var.subscription_alias_name != null : true
  }
  validation {
    error_message = "Value must not be null if `subscription_update_existing` is set to `true`."
    condition     = var.subscription_update_existing ? var.subscription_alias_name != null : true
  }
}

variable "subscription_billing_scope" {
  type        = string
  default     = null
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

E.g.

- For CustomerLed and FieldLed, e.g. MCA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
- For PartnerLed, e.g. MPA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}`
- For Legacy EA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION

  validation {
    condition     = var.subscription_billing_scope != null ? can(regex("^$|^/providers/Microsoft.Billing/billingAccounts/.*$", var.subscription_billing_scope)) : true
    error_message = "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive."
  }
  validation {
    error_message = "Value must not be null if `subscription_alias_enabled` is set to `true`."
    condition     = var.subscription_alias_enabled ? var.subscription_billing_scope != null : true
  }
}

variable "subscription_display_name" {
  type        = string
  default     = null
  description = <<DESCRIPTION
The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
The maximum length is 64 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION

  validation {
    condition     = var.subscription_display_name != null ? length(var.subscription_display_name) > 0 && length(var.subscription_display_name) <= 64 && !can(regex("[<>;|]", var.subscription_display_name)) : true
    error_message = "Subscription Name must be between 1 and 64 characters in length and cannot contain the characters `<`, `>`, `;`, or `|`"
  }
  validation {
    error_message = "Value must not be null if `subscription_alias_enabled` is set to `true`."
    condition     = var.subscription_alias_enabled ? var.subscription_display_name != null : true
  }
}

variable "subscription_id" {
  type        = string
  default     = null
  description = <<DESCRIPTION
DESCRIPTION

  validation {
    condition     = var.subscription_id != null ? can(regex("^$|^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id)) : true
    error_message = "Must be null, or a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "subscription_management_group_association_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create the subscription_association resource.

If enabled, the `subscription_management_group_id` must also be supplied.
DESCRIPTION
  nullable    = false
}

variable "subscription_management_group_id" {
  type        = string
  default     = null
  description = <<DESCRIPTION
The destination management group ID for the new subscription.

**Note:** Do not supply the display name.
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.
DESCRIPTION

  validation {
    condition     = var.subscription_management_group_id != null ? can(regex("^$|^[().a-zA-Z0-9_-]{1,90}$", var.subscription_management_group_id)) : true
    error_message = "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.)."
  }
  validation {
    error_message = "Value must not be null if `subscription_management_group_association_enabled` is set to `true`."
    condition     = var.subscription_management_group_association_enabled ? var.subscription_management_group_id != null : true
  }
}

variable "subscription_tags" {
  type        = map(string)
  default     = {}
  description = <<DESCRIPTION
A map of tags to assign to the newly created subscription.
Only valid when `subscription_alias_enabled` OR `subscription_update_existing` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```
DESCRIPTION

  validation {
    error_message = "Tag values must be between 0-256 characters."
    condition = alltrue(
      [for _, v in var.subscription_tags : can(regex("^.{0,256}$", v))]
    )
  }
  validation {
    error_message = "Tag name must contain neither `<>%&\\?/` nor control characters, and must be between 0-512 characters."
    condition = alltrue(
      [for k, _ in var.subscription_tags : can(regex("^[^<>%&\\?/[:cntrl:]]{0,512}$", k))]
    )
  }
}

variable "subscription_update_existing" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to update an existing subscription with the supplied tags and display name.
Must be set to `false` if `subscription_alias_enabled` is set to `true`.

If enabled, the following must also be supplied:
- `subscription_id`
DESCRIPTION
  nullable    = false

  validation {
    error_message = "Value must not be true if `subscription_id` is `null`."
    condition     = var.subscription_update_existing ? var.subscription_id != null : true
  }
  validation {
    error_message = "Value must not be true if `subscription_alias_enabled` is `true`."
    condition     = var.subscription_update_existing ? !var.subscription_alias_enabled : true
  }
}

variable "subscription_workload" {
  type        = string
  default     = null
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION

  validation {
    condition     = var.subscription_workload != null ? can(regex("^$|^(Production|DevTest)$", var.subscription_workload)) : true
    error_message = "The workload type can be either Production or DevTest and is case sensitive."
  }
  validation {
    error_message = "Value must not be null if `subscription_alias_enabled` is set to `true`."
    condition     = var.subscription_alias_enabled ? var.subscription_workload != null : true
  }
}

variable "wait_for_subscription_before_subscription_operations" {
  type = object({
    create  = optional(string, "30s")
    destroy = optional(string, "0s")
  })
  default     = {}
  description = <<DESCRIPTION
The duration to wait after vending a subscription before performing subscription operations.
DESCRIPTION
}
