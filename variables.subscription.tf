variable "subscription_alias_enabled" {
  default     = false
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
}

variable "subscription_alias_name" {
  type        = string
  nullable    = false
  default     = ""
  description = <<DESCRIPTION
The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and _.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_display_name" {
  type        = string
  nullable    = false
  default     = ""
  description = <<DESCRIPTION
The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, _ and space.
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_billing_scope" {
  type        = string
  nullable    = false
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
}

variable "subscription_workload" {
  type        = string
  nullable    = false
  default     = ""
  description = <<DESCRIPTION
The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.
DESCRIPTION
}

variable "subscription_management_group_id" {
  type        = string
  nullable    = false
  default     = ""
  description = <<DESCRIPTION
  The destination management group ID for the new subscription.

**Note:** Do not supply the display name.
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.
DESCRIPTION
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
  nullable    = false
  default     = ""
  description = <<DESCRIPTION
An existing subscription id.

Use this when you do not want the module to create a new subscription.
But do want to manage the management group membership.

A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
All letters must be lowercase.

When using this, `subscription_management_group_association_enabled` should be enabled,
and `subscription_management_group_id` should be supplied.

You may also supply an empty string if you want to create a new subscription alias.
In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:

- `subscription_alias_name`
- `subscription_alias_display_name`
- `subscription_alias_billing_scope`
- `subscription_alias_workload`
DESCRIPTION
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
  nullable    = false
  default     = {}
}

variable "subscription_use_azapi" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create a new subscription using the azapi provider. This may be required if the principal running
terraform does not have the required permissions to create a subscription under the default management group.
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
    time_grain        = string
    time_period_start = string
    time_period_end   = string
    notifications = optional(map(object({
      enabled        = bool
      operator       = string
      threshold      = number
      threshold_type = optional(string, "Actual")
      contact_emails = optional(list(string), [])
      contact_roles  = optional(list(string), [])
      contact_groups = optional(list(string), [])
      locale         = optional(string, "en-us")
    })), {})
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
      [
        for _, v in var.subscription_budgets :
        can(regex("^[0-9]{4}-[0-9]{2}-01T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$", v.time_period_start))
      ]
    )
    error_message = "Start date should be in the format yyyy-MM-01THH:mm:ssZ."
  }
  validation {
    condition = alltrue(
      [
        for _, v in var.subscription_budgets :
        timecmp(v.time_period_start, v.time_period_end) == -1 && can(regex("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$", v.time_period_end))
      ]
    )
    error_message = "Start date should be earlier than end date and in the format yyyy-MM-ddTHH:mm:ssZ."
  }
  validation {
    condition = alltrue(flatten(
      [
        for _, v in var.subscription_budgets :
        [
          for _, n in v.notifications :
          contains(["GreaterThan", "GreaterThanOrEqualTo"], n.operator)
        ]

      ]
    ))
    error_message = "Operator must be one of GreaterThan or GreaterThanOrEqualTo."
  }
  validation {
    condition = alltrue(flatten(
      [
        for _, v in var.subscription_budgets :
        [
          for _, n in v.notifications :
          contains(["Actual", "Forecasted"], n.threshold_type)
        ]

      ]
    ))
    error_message = "Threshold type must be one of Actual or Forecasted."
  }
  validation {
    condition = alltrue(flatten(
      [
        for _, v in var.subscription_budgets :
        [
          for _, n in v.notifications :
          n.threshold >= 0 && n.threshold <= 1000
        ]

      ]
    ))
    error_message = "Threshold must be between 0 and 1000."
  }
  validation {
    condition = alltrue(flatten(
      [
        for _, v in var.subscription_budgets :
        [
          for _, n in v.notifications :
          can(regex("^[a-z]{2}-[a-z]{2}$", n.locale))
        ]

      ]
    ))
    error_message = "Locale must be in the format xx-xx."
  }
  validation {
    condition = alltrue(flatten(
      [
        for _, v in var.subscription_budgets :
        [
          for _, n in v.notifications :
          length(n.contact_emails) > 0 || length(n.contact_roles) > 0 || length(n.contact_groups) > 0
        ]

      ]
    ))
    error_message = "At least one of contact_emails, contact_roles, or contact_groups must be supplied."
  }
  default     = {}
  description = <<DESCRIPTION
Map of budgets to create for the subscription.

- `amount` - The total amount of cost to track with the budget.
- `time_grain` - The time grain for the budget. Must be one of Annually, BillingAnnual, BillingMonth, BillingQuarter, Monthly, or Quarterly.
- `time_period_start` - The start date for the budget.
- `time_period_end` - The end date for the budget.
- `notifications` - The notifications to create for the budget.
  - `enabled` - Whether the notification is enabled.
  - `operator` - The operator for the notification. Must be one of GreaterThan or GreaterThanOrEqualTo.
  - `threshold` - The threshold for the notification. Must be between 0 and 1000.
  - `threshold_type` - The threshold type for the notification. Must be one of Actual or Forecasted.
  - `contact_emails` - The contact emails for the notification.
  - `contact_roles` - The contact roles for the notification.
  - `contact_groups` - The contact groups for the notification.
  - `locale` - The locale for the notification. Must be in the format xx-xx.


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
        enabled        = true
        operator       = "GreaterThan"
        threshold      = 80
        threshold_type = "Actual"
        contact_emails = ["john@contoso.com"]
      }
      budgetexceeded = {
        enabled        = true
        operator       = "GreaterThan"
        threshold      = 120
        threshold_type = "Forecasted"
        contact_groups = ["Owner"]
      }
    }
  }
}
```
DESCRIPTION
}
