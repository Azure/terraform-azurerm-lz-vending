variable "umi_enabled" {
  description = <<DESCRIPTION
Whether to enable the creation of a user-assigned managed identity.

Requires `umi.name` and `umi.resosurce_group_name` to be non-empty.
DESCRIPTION
  type        = bool
  default     = false
}

variable "user_managed_identities" {
  type = map(object({
    name                            = string
    resource_group_key  = optional(string)
    resource_group_name_existing = optional(string)
    location                        = optional(string)
    tags                            = optional(map(string), {})
    role_assignments = optional(map(object({
      definition                = string
      relative_scope            = optional(string, "")
      condition                 = optional(string)
      condition_version         = optional(string)
      principal_type            = optional(string)
      definition_lookup_enabled = optional(bool, true)
      use_random_uuid           = optional(bool, false)
    })), {})
    federated_credentials_github = optional(map(object({
      name            = optional(string)
      organization    = string
      repository      = string
      entity          = string
      enterprise_slug = optional(string)
      value           = optional(string)
    })), {})
    federated_credentials_terraform_cloud = optional(map(object({
      name         = optional(string)
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
- `resource_group_key`: The resource group key from the resource groups map to create the user assigned identity in. [optional]
- `resource_group_name_existing`: The name of an existing resource group to create the user assigned identity in. [optional]

**One of `resource_group_key` or `resource_group_name_existing` must be specified.**

### Optional fields

- `location`: The location of the user-assigned managed identity. [optional]
- `tags`: The tags to apply to the user-assigned managed identity. [optional]

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
  - `enterprise_slug` - the name of the GitHub Enterprise, e.g. `my-enterprise`. This is optional and only valid when using an enterprise.
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
  validation {
    condition = var.umi_enabled ? alltrue([
      for k, v in var.user_managed_identities : (
        (try(v.resource_group_key, null) != null) != (try(v.resource_group_name_existing, null) != null)
      )
    ]) : true
    error_message = "For each user-managed identity, set exactly one of resource_group_key or resource_group_name_existing."
  }
}
