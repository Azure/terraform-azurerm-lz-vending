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
    principal_id   = string,
    definition     = string,
    relative_scope = optional(string, ""),
  }))
  description = <<DESCRIPTION
Supply a map of objects containing the details of the role assignments to create.

Object fields:

- `principal_id`: The directory/object id of the principal to assign the role to.
- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: (optional) Scope relative to the created subscription. Omit, or leave blank for subscription scope.

E.g.

```terraform
role_assignments = {
  # Example using role definition name:
  contributor_user = {
    principal_id   = "00000000-0000-0000-0000-000000000000",
    definition     = "Contributor",
    relative_scope = "",
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
