variable "umi_enabled" {
  description = <<DESCRIPTION
Whether to enable the creation of a user-assigned managed identity.

Requires `umi_name` and `umi_resosurce_group_name` to be non-empty.
DESCRIPTION
  type        = bool
  default     = false
}

variable "user_managed_identities" {
  type = map(object({
    name                        = string
    resource_group_name         = string
    location                    = optional(string, null)
    tags                        = optional(map(string), {})
    resource_group_tags         = optional(map(string), {})
    resource_group_lock_enabled = optional(bool, true)
    resource_group_lock_name    = optional(string, null)
    role_assignments = optional(map(object({
      definition                = string
      relative_scope            = optional(string, null)
      condition                 = optional(string, null)
      condition_version         = optional(string, null)
      principal_type            = optional(string, null)
      definition_lookup_enabled = optional(bool, true)
    })), {})
    federated_credentials_github = optional(map(object({
      name         = optional(string, null)
      organization = string
      repository   = string
      entity       = string
      value        = optional(string, null)
    })), {})
    federated_credentials_terraform_cloud = optional(map(object({
      name         = optional(string, null)
      organization = string
      project      = string
      workspace    = string
      run_phase    = string
    })), {})
    federated_credentials_advanced = optional(map(object({
      name               = string
      subject_identifier = string
      issuer_url         = string
      audiences          = optional(set(string), ["api://AzureADTokenExchange"])
    })), {})
  }))
  default     = {}
  description = <<DESCRIPTION
A map of user-managed identities to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply. The value is a map of attributes.

### Required fields

- `name`: The name of the user-assigned managed identity. [required]
- `resource_group_name`: The name of the resource group to create the user-assigned managed identity in. [required]

### Optional fields

- `location`: The location of the user-assigned managed identity. [optional]
- `tags`: The tags to apply to the user-assigned managed identity. [optional]
- `resource_group_tags`: The tags to apply to the user-assigned managed identity resource group, if we create it. [optional]
- `resource_group_lock_enabled`: Whether to enable resource group lock for the user-assigned managed identity resource group. [optional]
- `resource_group_lock_name`: The name of the resource group lock for the user-assigned managed identity resource group, if blank will be set to `lock-<resource_group_name>`. [optional]

### Role Based Access Control (RBAC)

The following fields are used to configure role assignments for the user-assigned managed identity.
- `role_assignments`: A map of role assignments to create for the user-assigned managed identity. [optional] - See `role_assignments` variable for details.

### Federated Credentials

The following fields are used to configure federated identity credentials, using OpenID Connect, for use in GitHub actions, Azure DevOps pipelines, and Terraform Cloud.

#### GitHub Actions

- `federated_credentials_github`: A map of federated credentials to create for the user-assigned managed identity. [optional]
  - `name` - the name of the federated credential resource, the last segment of the Azure resource id.
  - `organization` - the name of the GitHub organization, e.g. `Azure` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
  - `repository` - the name of the GitHub respository, e.g. `terraform-azurerm-lz-vending` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
  - `entity` - one of 'environment', 'pull_request', 'tag', or 'branch'
  - `value` - identifies the `entity` type, e.g. `main` when using entity is `branch`. Should be blank when `entity` is `pull_request`.

#### Terraform Cloud

- `federated_credentials_terraform_cloud`: A map of federated credentials to create for the user-assigned managed identity. [optional]
  - `name` - the name of the federated credential resource, the last segment of the Azure resource id.
  - `organization` - the name of the Terraform Cloud organization.
  - `project` - the name of the Terraform Cloud project.
  - `workspace` - the name of the Terraform Cloud workspace.
  - `run_phase` - one of `plan`, or `apply`.

#### Advanced Federated Credentials

- `federated_credentials_advanced`: A map of federated credentials to create for the user-assigned managed identity. [optional]
  - `name`: The name of the federated credential resource, the last segment of the Azure resource id.
  - `subject_identifier`: The subject of the token.
  - `issuer_url`: The URL of the token issuer, should begin with `https://`
  - `audience`: (optional) The token audience, defaults to `api://AzureADTokenExchange`.
DESCRIPTION
}
