variable "name" {
  description = "The name of the user managed identity"
  type        = string
}

variable "location" {
  description = "The location of the user managed identity"
  type        = string
}

variable "parent_id" {
  description = "The parent resource ID (resource group or subscription)"
  type        = string
}

variable "tags" {
  description = "Tags to apply to the user managed identity"
  type        = map(string)
  default     = {}
}

variable "federated_credentials_github" {
  description = "GitHub federated credentials"
  type = map(object({
    organization = string
    repository   = string
    entity       = string
    value        = optional(string)
  }))
  default = {}
}

variable "federated_credentials_terraform_cloud" {
  description = "Terraform Cloud federated credentials"
  type = map(object({
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
  default = {}
}

variable "federated_credentials_advanced" {
  description = "Advanced federated credentials"
  type = map(object({
    name               = string
    subject_identifier = string
    issuer_url         = string
  }))
  default = {}
}
