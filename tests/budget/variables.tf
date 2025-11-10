variable "budget_name" {
  description = "The name of the budget"
  type        = string
}

variable "budget_scope" {
  description = "The scope of the budget (subscription or resource group)"
  type        = string
}

variable "budget_amount" {
  description = "The total amount of the budget"
  type        = number
}

variable "budget_time_grain" {
  description = "The time covered by the budget (Monthly, Quarterly, Annually)"
  type        = string
}

variable "budget_time_period" {
  description = "The time period for the budget"
  type = object({
    start_date = string
    end_date   = optional(string)
  })
}

variable "budget_notifications" {
  description = "Notifications for the budget"
  type = map(object({
    enabled         = bool
    operator        = string
    threshold       = number
    threshold_type  = string
    contact_emails  = optional(list(string))
    contact_roles   = optional(list(string))
    contact_groups  = optional(list(string))
    locale          = optional(string)
  }))
  default = {}
}
