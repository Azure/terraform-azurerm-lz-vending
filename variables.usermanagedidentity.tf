variable "umi_enabled" {
  description = "Whether to enable the creation of a user-assigned managed identity."
  type        = bool
  default     = false
}

variable "umi_name" {
  description = "The name of the user-assigned managed identity"
  type        = string
}

variable "umi_tags" {
  description = "The tags to apply to the user-assigned managed identity"
  type        = map(string)
  default     = {}
}

variable "umi_resource_group_creation_enabled" {
  description = "Whether to create the supplied resource group for the user-assigned managed identity"
  type        = bool
  default     = true
}

variable "umi_resource_group_name" {
  description = "The name of the resource group in which to create the user-assigned managed identity"
  type        = string
}

variable "umi_resource_group_tags" {
  description = "The tags to apply to the user-assigned managed identity resource group, if we create it."
  type        = map(string)
  default     = {}
}

variable "umi_resource_group_lock_enabled" {
  description = "Whether to enable resource group lock for the user-assigned managed identity resource group"
  type        = bool
  default     = true
}

variable "umi_resource_group_lock_name" {
  description = "The name of the resource group lock for the user-assigned managed identity resource group, if blank will be set to `lock-<resource_group_name>`"
  type        = string
  default     = ""
}

# allow the caller to easily configure federated credentials for GitHub Actions
variable "umi_federated_credentials_github" {
  type = map(object({
    name         = optional(string, "")
    organization = string
    repository   = string
    entity       = string
    value        = optional(string, "")
  }))
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
  default     = {}
}

# allow the caller to easily configure federated credentials for Terraform Cloud
variable "umi_federated_credentials_terraform_cloud" {
  type = map(object({
    name         = optional(string, "")
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
}

# allow the caller to configure federated credentials by supplying the values verbatim
variable "umi_federated_credentials_advanced" {
  type = map(object({
    name               = string
    subject_identifier = string
    issuer_url         = string
    audience           = optional(string, "api://AzureADTokenExchange")
  }))
  default     = {}
  description = <<DESCRIPTION
Configure federated identity credentials, using OpenID Connect, for use scenarios outside GitHub Actions and Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `subject_identifier` - The subject of the token.
- `issuer_url` - the URL of the token issuer, should begin with `https://`
- `audience` - (optional) the token audience, defaults to `api://AzureADTokenExchange`.
DESCRIPTION
}

variable "umi_role_assignments" {
  type = map(object({
    definition     = string
    relative_scope = string
  }))
  default     = {}
  description = <<DESCRIPTION
Supply a map of objects containing the details of the role assignments to create for the user-assigned managed identity.
This will be merged with the other role assignments specified in `var.role_assignments`

Requires both `var.umi_enabled` and `var.role_assignment_enabled` to be `true`.

Object fields:

- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: Scope relative to the created subscription. Leave blank for subscription scope.
DESCRIPTION
}
