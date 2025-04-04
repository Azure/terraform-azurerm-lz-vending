variable "role_assignment_enabled" {
  type        = bool
  description = <<DESCRIPTION
Whether to create role assignments.
If enabled, supply the list of role assignments in `var.role_assignments`.
DESCRIPTION
  default     = false
}

variable "role_assignments" {
  type = map(object({
    principal_id              = string,
    definition                = string,
    relative_scope            = optional(string, null)
    condition                 = optional(string, null)
    condition_version         = optional(string, null)
    principal_type            = optional(string, null)
    definition_lookup_enabled = optional(bool, true)
  }))
  description = <<DESCRIPTION
Supply a map of objects containing the details of the role assignments to create.

Object fields:

- `principal_id`: The directory/object id of the principal to assign the role to.
- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.
- `condition`: (optional) A condition to apply to the role assignment. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.
- `condition_version`: (optional) The version of the condition syntax. See [Conditions Custom Security Attributes](https://learn.microsoft.com/azure/role-based-access-control/conditions-custom-security-attributes) for more details.
- `principal_type`: (optional) The type of the principal. Can be `"User"`, `"Group"`, `"Device"`, `"ForeignGroup"`, or `"ServicePrincipal"`.
- `definition_lookup_enabled`: (optional) Whether to look up the role definition resource id from the role definition name. If disabled, the `definition` must be a role definition resource id. Default is `true`.


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
}
