variable "name" {
  description = "The name of the user managed identity"
  type        = string
}

variable "location" {
  description = "The location of the user managed identity"
  type        = string

}

variable "tags" {
  description = "The tags to apply to the user managed identity"
  type        = map(string)
  default     = {}
}

variable "subscription_id" {
  description = "The subscription id"
  type        = string
  validation {
    condition     = can(regex("^^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must be a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "resource_group_creation_enabled" {
  description = "Whether to create the supplied resource group"
  type        = bool
  default     = true
}

variable "resource_group_name" {
  description = "The name of the resource group in which to create the user managed identity"
  type        = string
}

variable "resource_group_tags" {
  description = "The tags to apply to the resource group, if we create it."
  type        = map(string)
  default     = {}
}

variable "resource_group_lock_enabled" {
  description = "Whether to enable resource group lock"
  type        = bool
  default     = true
}

variable "resource_group_lock_name" {
  description = "The name of the resource group lock, if blank will be set to `lock-<resource_group_name>`"
  type        = string
  default     = ""
}

# allow the caller to easily configure federated credentials for GitHub Actions
variable "federated_credentials_github" {
  type = map(object({
    name         = optional(string, "")
    organization = string
    repository   = string
    entity       = string
    value        = optional(string, "")
  }))
  default = {}

  validation {
    condition = alltrue([
      for v in var.federated_credentials_github : contains(["environment", "pull_request", "tag", "branch"], v.entity)
    ])
    error_message = "Entity must be one of 'environment', 'pull_request', 'tag', or 'branch'."
  }

  validation {
    condition = alltrue([
      for v in var.federated_credentials_github : v.value != ""
      if v.entity != "pull_request"
    ])
    error_message = "Field 'value' must be specified for all entities except 'pull_request'."
  }
}

# allow the caller to easily configure federated credentials for Terraform Cloud
variable "federated_credentials_terraform_cloud" {
  type = map(object({
    name         = optional(string, "")
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
  default = {}

  validation {
    condition = alltrue([
      for v in var.federated_credentials_terraform_cloud : contains(["apply", "plan"], v.run_phase)
    ])
    error_message = "Field 'run_phase' value must be 'plan' or 'apply'."
  }
}

# allow the caller to configure federated credentials by supplying the values verbatim
variable "federated_credentials_advanced" {
  type = map(object({
    name               = string
    subject_identifier = string
    audience           = optional(string, "api://AzureADTokenExchange")
    issuer_url         = string
  }))
  default = {}
}
