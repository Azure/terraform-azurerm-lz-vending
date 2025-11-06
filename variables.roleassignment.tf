variable "role_assignment_enabled" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to create role assignments.
If enabled, supply the list of role assignments in `var.role_assignments`.
DESCRIPTION
}

variable "role_assignments" {
  type = map(object({
    principal_id              = string,
    definition                = string,
    relative_scope            = optional(string, "")
    resource_group_scope_key  = optional(string)
    condition                 = optional(string)
    condition_version         = optional(string)
    principal_type            = optional(string)
    definition_lookup_enabled = optional(bool, false)
    use_random_uuid           = optional(bool, false)
  }))
  default     = {}
  description = <<DESCRIPTION
Supply a map of objects containing the details of the role assignments to create.

Object fields:

- `principal_id`: The directory/object id of the principal to assign the role to.
- `definition`: The role definition to assign. Either use the name or the role definition resource id. If supplying a definition ID, use a *scopeless* role definition ID (i.e. starting with `/providers/Microsoft.Authorization/roleDefinitions/`).
- `relative_scope`: (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.
- `resource_group_scope_key`: (optional) The resource group key from the resource groups map to use as the scope for the role assignment. If supplied, this takes precedence over `relative_scope`.
- `condition`: (optional) A condition to apply to the role assignment. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.
- `condition_version`: (optional) The version of the condition syntax. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.
- `principal_type`: (optional) The type of the principal. Can be `"User"`, `"Group"`, `"Device"`, `"ForeignGroup"`, or `"ServicePrincipal"`.
- `definition_lookup_enabled`: (optional) Whether to look up the role definition resource id from the the Azure API. Default is `false`, where we use a static module of role definitions.
- `use_random_uuid`: (optional) Whether to use a random UUID for the role assignment name. Default is `false`. If set to `true`, the role assignment name will be a random UUID, otherwise it will be a deterministic UUID based on the scope, principal id, and role definition id.

E.g.

```terraform
role_assignments = {
  # Example using role definition name:
  contributor_user = {
    principal_id      = "00000000-0000-0000-0000-000000000000",
    definition        = "Contributor",
    relative_scope    = "",
    condition         = "(!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'} AND NOT SubOperationMatches{'Blob.List'})",
    condition_version = "2.0",
  },
  # Example using role definition id and RG scope:
  myrg_custom_role = {
    principal_id   = "11111111-1111-1111-1111-111111111111",
    definition     = "/providers/Microsoft.Management/managementGroups/mymg/providers/Microsoft.Authorization/roleDefinitions/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    relative_scope = "/resourceGroups/MyRg",
  }
}
```
DESCRIPTION
  nullable    = false
  default     = {}

  validation {
    error_message = "If definition is a role definition ID, it must start with /providers/Microsoft.Authorization/roleDefinitions/ to be a scopeless role definition ID."
    condition = alltrue([for ra in values(var.role_assignments) : (
      strcontains(lower(ra.definition), lower("/providers/Microsoft.Authorization/roleDefinitions/")) == false ||
      can(regex("^/providers/Microsoft\\.Authorization/roleDefinitions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", ra.definition))
    )])
  }
}
