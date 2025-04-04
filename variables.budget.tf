variable "budget_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether to create budgets.
If enabled, supply the list of budgets in `var.budgets`.
DESCRIPTION
  default     = false
}

variable "budgets" {
  type = map(object({
    amount            = number
    time_grain        = string
    time_period_start = string
    time_period_end   = string
    relative_scope    = optional(string, null)
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
  default     = {}
  description = <<DESCRIPTION
Map of budgets to create for the subscription.

- `amount` - The total amount of cost to track with the budget.
- `time_grain` - The time grain for the budget. Must be one of Annually, BillingAnnual, BillingMonth, BillingQuarter, Monthly, or Quarterly.
- `time_period_start` - The start date for the budget.
- `time_period_end` - The end date for the budget.
- `relative_scope` - (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.
- `notifications` - (optional) The notifications to create for the budget.
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
        contact_roles  = ["Owner"]
      }
    }
  }
}
```
DESCRIPTION
}
