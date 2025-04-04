variable "budget_amount" {
  type        = number
  description = "The total amount of cost to track with the budget."
  nullable    = false
}

variable "budget_name" {
  type        = string
  description = "The name of the budget."
  nullable    = false
}

variable "budget_scope" {
  type        = string
  description = "The scope of the budget."
  nullable    = false
}

variable "budget_time_grain" {
  type        = string
  description = "The time grain of the budget."
  nullable    = false

  validation {
    condition     = contains(["Annually", "BillingAnnual", "BillingMonth", "BillingQuarter", "Monthly", "Quarterly"], var.budget_time_grain)
    error_message = "Time period must be one of Annually, BillingAnnual, BillingMonth, BillingQuarter, Monthly, or Quarterly."
  }
}

variable "budget_time_period" {
  type = object({
    start_date = string
    end_date   = string
  })
  description = "The time period of the budget."
  nullable    = false

  validation {
    condition     = can(regex("^[0-9]{4}-[0-9]{2}-01T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$", var.budget_time_period.start_date))
    error_message = "Start date should be in the format yyyy-MM-01THH:mm:ssZ."
  }
  validation {
    condition     = timecmp(var.budget_time_period.start_date, var.budget_time_period.end_date) == -1 && can(regex("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$", var.budget_time_period.end_date))
    error_message = "Start date should be earlier than end date and in the format yyyy-MM-ddTHH:mm:ssZ."
  }
}

variable "budget_notifications" {
  type = map(object({
    enabled        = bool
    operator       = string
    threshold      = number
    threshold_type = optional(string, "Actual")
    contact_emails = optional(list(string), [])
    contact_roles  = optional(list(string), [])
    contact_groups = optional(list(string), [])
    locale         = optional(string, "en-us")
  }))
  default     = {}
  description = "The notifications for the budget."
  nullable    = false

  validation {
    condition     = length(keys(var.budget_notifications)) <= 5
    error_message = "Maximum number of notifications per budget is 5."
  }
  validation {
    condition = alltrue([
      for notification in var.budget_notifications : contains(["GreaterThan", "GreaterThanOrEqualTo"], notification.operator)
    ])

    error_message = "Operator must be one of GreaterThan or GreaterThanOrEqualTo."
  }
  validation {
    condition = alltrue([
      for notification in var.budget_notifications :
      contains(["Actual", "Forecasted"], notification.threshold_type)
    ])
    error_message = "Threshold type must be one of Actual or Forecasted."
  }
  validation {
    condition = alltrue([
      for notification in var.budget_notifications : notification.threshold >= 0 && notification.threshold <= 1000
    ])
    error_message = "Threshold must be between 0 and 1000."
  }
  validation {
    condition = alltrue([
      for notification in var.budget_notifications : can(regex("^[a-z]{2}-[a-z]{2}$", notification.locale))
    ])
    error_message = "Locale must be in the format xx-xx."
  }
  validation {
    condition = alltrue([
      for notification in var.budget_notifications : length(notification.contact_emails) > 0 || length(notification.contact_roles) > 0 || length(notification.contact_groups) > 0
    ])
    error_message = "At least one of contact_emails, contact_roles, or contact_groups must be supplied."
  }
}
