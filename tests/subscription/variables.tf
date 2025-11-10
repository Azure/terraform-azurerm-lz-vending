variable "subscription_alias_enabled" {
  description = "Whether to enable the subscription alias"
  type        = bool
  default     = false
}

variable "subscription_alias_name" {
  description = "The name of the subscription alias"
  type        = string
  default     = ""
}

variable "subscription_billing_scope" {
  description = "The billing scope for the subscription"
  type        = string
  default     = ""
}

variable "subscription_display_name" {
  description = "The display name for the subscription"
  type        = string
  default     = ""
}

variable "subscription_workload" {
  description = "The workload type for the subscription (Production or DevTest)"
  type        = string
  default     = "Production"
}

variable "subscription_id" {
  description = "The ID of an existing subscription"
  type        = string
  default     = ""
}

variable "subscription_management_group_id" {
  description = "The management group ID to associate the subscription with"
  type        = string
  default     = ""
}

variable "subscription_management_group_association_enabled" {
  description = "Whether to enable management group association"
  type        = bool
  default     = false
}

variable "subscription_tags" {
  description = "Tags to apply to the subscription"
  type        = map(string)
  default     = {}
}

variable "subscription_update_existing" {
  description = "Whether to update an existing subscription"
  type        = bool
  default     = false
}
