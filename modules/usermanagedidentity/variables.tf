variable "umi_enabled" {
  description = <<DESCRIPTION
Whether to enable the creation of a user-assigned managed identity.
DESCRIPTION
  type        = bool
  default     = false
}

variable "location" {
  type        = string
  description = "The location of the user-assigned managed identity"
}

variable "name" {
  type        = string
  description = "The name of the user managed identity"
  nullable    = false
}

variable "parent_id" {
  type        = string
  description = "The ID of the parent resource to which this user-assigned managed identity."

  validation {
    condition     = var.umi_enabled ? length(var.parent_id) > 0 : true
    error_message = "The parent_id must not be empty."
  }
  validation {
    condition     = can(regex("^/subscriptions/[a-fA-F0-9-]+/resourceGroups/[a-zA-Z0-9-_.()]+$", var.parent_id))
    error_message = "The parent_id must be a valid Azure Resource Group ID."
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
    name            = optional(string)
    organization    = string
    repository      = string
    entity          = string
    enterprise_slug = optional(string)
    value           = optional(string)
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
- `enterprise_slug` - the name of the GitHub Enterprise, e.g. `my-enterprise`. This is optional and only valid when using an enterprise.
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

variable "tags" {
  type        = map(string)
  default     = {}
  description = "The tags to apply to the user-assigned managed identity"
}
