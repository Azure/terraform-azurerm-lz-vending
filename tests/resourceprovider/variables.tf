variable "resource_provider" {
  description = "The resource provider to register"
  type        = string
}

variable "features" {
  description = "List of resource provider features to register"
  type        = list(string)
  default     = []
}

variable "subscription_id" {
  description = "The subscription ID"
  type        = string
}
