variable "subscription_alias_enabled" {
  type        = bool
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

  default = false
}

variable "subscription_alias_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and _.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
  validation {
    condition     = can(regex("^$|^[a-zA-Z0-9-_]{1,63}$", var.subscription_alias_name))
    error_message = "Valid characters are a-z, A-Z, 0-9, -, _."
  }
}

variable "subscription_display_name" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
  validation {
    condition     = can(regex("^$|^[a-zA-Z0-9-_ ]{1,63}$", var.subscription_display_name))
    error_message = "Valid characters are a-z, A-Z, 0-9, -, _, and space."
  }
}

variable "subscription_billing_scope" {
  type        = string
  default     = ""
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
    condition     = can(regex("^$|^/providers/Microsoft.Billing/billingAccounts/.*$", var.subscription_billing_scope))
    error_message = "A valid billing scope starts with /providers/Microsoft.Billing/billingAccounts/ and is case sensitive."
  }
}

variable "subscription_workload" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
  validation {
    condition     = can(regex("^$|^(Production|DevTest)$", var.subscription_workload))
    error_message = "The workload type can be either Production or DevTest and is case sensitive."
  }
}

variable "subscription_management_group_id" {
  type        = string
  default     = ""
  description = <<DESCRIPTION
The destination management group ID for the new subscription.

**Note:** Do not supply the display name.
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.
DESCRIPTION
  validation {
    condition     = can(regex("^$|^[().a-zA-Z0-9_-]{1,90}$", var.subscription_management_group_id))
    error_message = "The management group ID must be between 1 and 90 characters in length and formed of the following characters: a-z, A-Z, 0-9, -, _, (, ), and a period (.)."
  }
}

variable "subscription_management_group_association_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create the `azurerm_management_group_subscription_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.
DESCRIPTION
}

variable "subscription_id" {
  type        = string
  description = <<DESCRIPTION
DESCRIPTION
  default     = ""
  validation {
    condition     = can(regex("^$|^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must be empty, or a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "subscription_tags" {
  type        = map(string)
  description = <<DESCRIPTION
A map of tags to assign to the newly created subscription.
Only valid when `subsciption_alias_enabled` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```
DESCRIPTION
  default     = {}
  validation {
    error_message = "Tag values must contain neither `<>%&\\?/` nor control characters, and must be between 0-256 characters."
    condition = alltrue(
      [for _, v in var.subscription_tags : can(regex("^[^<>%&\\?/[:cntrl:]]{0,256}$", v))]
    )
  }
  validation {
    error_message = "Tag name must contain neither `<>%&\\?/` nor control characters, and must be between 0-512 characters."
    condition = alltrue(
      [for k, _ in var.subscription_tags : can(regex("^[^<>%&\\?/[:cntrl:]]{0,512}$", k))]
    )
  }
}

variable "subscription_use_azapi" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to use the azapi_resource resource to create the subscription alias. This includes the subscription alias in the management group.
DESCRIPTION
}

variable "subscription_update_existing" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to update an existing subscription with the supplied tags and display name.
If enabled, the following must also be supplied:
- `subscription_id`
DESCRIPTION
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

variable "subscription_budgets" {
  type = map(object({
    amount            = number
    time_grain        = optional(string, "Monthly")
    time_period_start = string
    time_period_end   = string
    notifications = map(object({
      enabled       = bool
      operator      = string # EqualTo, GreaterThan, GreaterThanOrEqualTo
      threshold     = number # 0-1000 percent
      thresholdType = string # Actual, Forecasted
      contactEmails = optional(list(string), [])
      contactRoles  = optional(list(string), [])
      contactGroups = optional(list(string), [])
      locale        = optional(string, "en-us")
    }))
  }))
  validation {
    condition = alltrue(
      [for _, v in var.subscription_budgets : contains(["Annually", "BillingAnnual", "BillingMonth", "BillingQuarter", "Monthly", "Quarterly"], v.time_grain)]
    )
    error_message = "Time period must be one of Annually, BillingAnnual, BillingMonth, BillingQuarter, Monthly, or Quarterly. BillingMonth, BillingQuarter, and BillingAnnual are only supported by WD customers."
  }
  validation {
    condition = alltrue(
      [for _, v in var.subscription_budgets : length(keys(v.notifications)) <= 5]
    )
    error_message = "Maximum number of notifications per budget is 5."
  }
  validation {
    condition = alltrue(
      [for _, v in var.subscription_budgets : timecmp(v.time_period_start, timestamp()) == 1]
    )
    error_message = "Start date should be in the future."
  }
  validation {
    condition = alltrue(
      [for _, v in var.subscription_budgets : timecmp(v.time_period_start, v.time_period_end) == -1]
    )
    error_message = "Start date should be earlier than end date."
  }
  default     = {}
  description = <<DESCRIPTION
The budgets to create for the subscription using the AzApi provider.

time_period_start and time_period_end must be UTC in RFC3339 format, e.g. 2018-05-13T07:44:12Z.

Example value:

```terraform
subscription_budgets = {
  budget1 = {
    amount            = 150
    time_grain        = "Monthly"
    time_period_start = "2024-01-01T00:00:00Z"
    time_period_end   = "2027-12-31T23:59:59Z"
    notifications = {
      eightypercent = {
        enabled       = true
        operator      = "GreaterThan"
        threshold     = "80"
        thresholdType = "Actual"
        contactEmails = ["john@contoso.com"]
      }
      budgetexceeded = {
        enabled       = true
        operator      = "GreaterThan"
        threshold     = "120"
        thresholdType = "Forecasted"
        contactGroups = ["Owner"]
      }
    }
  }
}
```
DESCRIPTION
}
