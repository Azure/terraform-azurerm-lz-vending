variable "location" {
  type        = string
  description = "The location of the user-assigned managed identity"
}

variable "name" {
  type        = string
  description = "The name of the user managed identity"
  nullable    = false
}

variable "resource_group_name" {
  type        = string
  description = "The name of the resource group in which to create the user-assigned managed identity"
  nullable    = false

  validation {
    condition     = var.resource_group_name != ""
    error_message = "Resource group name must not be empty."
  }
}

variable "subscription_id" {
  type        = string
  description = "The id of the target subscription. Must be a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."

  validation {
    condition     = can(regex("^^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.subscription_id))
    error_message = "Must be a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

# allow the caller to configure federated credentials by supplying the values verbatim
variable "federated_credentials_advanced" {
  type = map(object({
    name               = string
    subject_identifier = string
    audiences          = optional(set(string), ["api://AzureADTokenExchange"])
    issuer_url         = string
  }))
  default     = {}
  description = <<DESCRIPTION
Configure federated identity credentials, using OpenID Connect, for use scenarios outside GitHub Actions and Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `subject_identifier` - The subject of the token.
- `issuer_url` - the URL of the token issuer, should begin with `https://`
- `audiences` - (optional) a set of strings containing the token audiences, defaults to `["api://AzureADTokenExchange"]`.
DESCRIPTION
}

# allow the caller to easily configure federated credentials for GitHub Actions
variable "federated_credentials_github" {
  type = map(object({
    name         = optional(string)
    organization = string
    repository   = string
    entity       = string
    value        = optional(string)
  }))
  default     = {}
  description = <<DESCRIPTION
Configure federated identity credentials, using OpenID Connect, for use in GitHub actions.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the GitHub organization, e.g. `Azure` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `repository` - the name of the GitHub respository, e.g. `terraform-azurerm-lz-vending` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `entity` - one of 'environment', 'pull_request', 'tag', or 'branch'
- `value` - identifies the `entity` type, e.g. `main` when using entity is `branch`. Should be blank when `entity` is `pull_request`.
DESCRIPTION

  validation {
    condition = alltrue([
      for v in var.federated_credentials_github : contains(["environment", "pull_request", "tag", "branch"], v.entity)
    ])
    error_message = "Entity must be one of 'environment', 'pull_request', 'tag', or 'branch'."
  }
  validation {
    condition = alltrue([
      for v in var.federated_credentials_github : v.value != null
      if v.entity != "pull_request"
    ])
    error_message = "Field 'value' must be specified for all entities except 'pull_request'."
  }
}

# allow the caller to easily configure federated credentials for Terraform Cloud
variable "federated_credentials_terraform_cloud" {
  type = map(object({
    name         = optional(string)
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
  default     = {}
  description = <<DESCRIPTION
Configure federated identity credentials, using OpenID Connect, for use in Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the Terraform Cloud organization.
- `project` - the name of the Terraform Cloud project.
- `workspace` - the name of the Terraform Cloud workspace.
- `run_phase` - one of `plan`, or `apply`.
DESCRIPTION

  validation {
    condition = alltrue([
      for v in var.federated_credentials_terraform_cloud : contains(["apply", "plan"], v.run_phase)
    ])
    error_message = "Field 'run_phase' value must be 'plan' or 'apply'."
  }
}

variable "resource_group_creation_enabled" {
  type        = bool
  default     = true
  description = "Whether to create the supplied resource group for the user-assigned managed identity"
  nullable    = false
}

variable "resource_group_lock_enabled" {
  type        = bool
  default     = true
  description = "Whether to enable resource group lock for the user-assigned managed identity resource group"
  nullable    = false
}

variable "resource_group_lock_name" {
  type        = string
  default     = null
  description = "The name of the resource group lock for the user-assigned managed identity resource group, if `null` will be set to `lock-<resource_group_name>`"
}

variable "resource_group_tags" {
  type        = map(string)
  default     = {}
  description = "The tags to apply to the user-assigned managed identity resource group, if we create it."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "The tags to apply to the user-assigned managed identity"
}
