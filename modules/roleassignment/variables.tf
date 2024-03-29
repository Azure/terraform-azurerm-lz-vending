variable "role_assignment_principal_id" {
  type        = string
  description = <<DESCRIPTION
The principal (object) ID of the role assignment.
Note, for a service principal, this is not the application id.

Can be user, group or service principal.
DESCRIPTION
  validation {
    condition     = can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.role_assignment_principal_id))
    error_message = "Must a GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. All letters must be lowercase."
  }
}

variable "role_assignment_scope" {
  type        = string
  description = <<DESCRIPTION
The scope of the role assignment.

Must begin with `/subscriptions/{subscription-id}` to avoid accidentally creating a role assignment at higher scopes.
DESCRIPTION
  validation {
    condition     = can(regex("^/subscriptions/[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}", var.role_assignment_scope))
    error_message = "Must begin with a subscription scope, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`. All letters must be lowercase in the subscription id."
  }
}

variable "role_assignment_definition" {
  type        = string
  description = <<DESCRIPTION
Either the role definition resource id, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c`.
Or, the role definition name, e.g. `Contributor`.
DESCRIPTION
}

variable "role_assignment_condition" {
  type        = string
  description = <<DESCRIPTION
(Optional) The condition that limits the resources that the role can be assigned to.
DESCRIPTION
}

variable "role_assignment_condition_version" {
  type        = string
  description = <<DESCRIPTION
The version of the condition. Possible values are `""`, 1.0 or 2.0. If `""`, null will be set in role_assignment_condition and role_assignment_condition_version.
DESCRIPTION
  validation {
    condition     = contains(["", "1.0", "2.0"], var.role_assignment_condition_version)
    error_message = "Must be version 1.0 or 2.0."
  }
}
