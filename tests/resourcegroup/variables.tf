variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

variable "location" {
  description = "The location of the resource group"
  type        = string
}

variable "subscription_id" {
  description = "The subscription ID"
  type        = string
}

variable "tags" {
  description = "Tags to apply to the resource group"
  type        = map(string)
  default     = {}
}

variable "lock_enabled" {
  description = "Whether to enable resource group lock"
  type        = bool
  default     = false
}
