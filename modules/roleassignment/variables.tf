variable "role_assignment_condition" {
  type        = string
  description = <<DESCRIPTION
(Optional) The condition that limits the resources that the role can be assigned to.
DESCRIPTION
  default     = null
}

variable "role_assignment_condition_version" {
  type        = string
  description = <<DESCRIPTION
The version of the condition. Possible values are `null`, 1.0 or 2.0. If `null` then `role_assignment_condition` will also be null.
DESCRIPTION

  validation {
    condition     = var.role_assignment_condition_version != null ? contains(["1.0", "2.0"], var.role_assignment_condition_version) : true
    error_message = "Must be version 1.0 or 2.0."
  }
  default = null
}

variable "role_assignment_definition" {
  type        = string
  description = <<DESCRIPTION
Either the role definition resource id, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c`.
Or, the role definition name, e.g. `Contributor`.
DESCRIPTION
}

variable "role_assignment_principal_id" {
  type        = string
  description = <<DESCRIPTION
The principal (object) ID of the role assignment.
Note, for a service principal, this is not the application id.

Can be user, group or service principal.
DESCRIPTION

  validation {
    condition     = var.role_assignment_principal_id == "skip" ? true : can(regex("^[a-f\\d]{4}(?:[a-f\\d]{4}-){4}[a-f\\d]{12}$", var.role_assignment_principal_id))
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

variable "role_assignment_definition_lookup_enabled" {
  type        = bool
  default     = true
  description = <<DESCRIPTION
Whether to look up the role definition resource id from the role definition name.
If disabled, the `role_assignment_definition` must be a role definition resource id.
DESCRIPTION
  nullable    = false
}

variable "role_assignment_principal_type" {
  type        = string
  default     = null
  description = <<DESCRIPTION
Required when using attribute based access control (ABAC).
The type of principal. Can be `User`, `Group`, `ServicePrincipal`, `Device`, or `ForeignGroup`.
DESCRIPTION

  validation {
    condition     = var.role_assignment_principal_type != null ? can(regex("^(User|Group|ServicePrincipal|Device|ForeignGroup)$", var.role_assignment_principal_type)) : true
    error_message = "Must be one of User, Group, ServicePrincipal, Device, or ForeignGroup."
  }
}

variable "role_assignment_use_random_uuid" {
  type        = bool
  default     = false
  description = <<DESCRIPTION
Whether to use a random UUID for the role assignment name.

> NOTE: Use this option to prevent unknown values causing role assignments to be recreated on every plan/apply. However make sure to use a new module call (UUID) if you change the properties of a role assignment.
DESCRIPTION
  nullable    = false
}
