locals {
  role_assignment_definition_id   = can(regex("/providers/Microsoft.Authorization/roleDefinitions", var.role_assignment_definition)) ? var.role_assignment_definition : null
  role_assignment_definition_name = can(regex("/providers/Microsoft.Authorization/roleDefinitions", var.role_assignment_definition)) ? null : var.role_assignment_definition

  role_assignment_condition         = var.role_assignment_condition == "" ? null : var.role_assignment_condition
  role_assignment_condition_version = var.role_assignment_condition == "" ? null : var.role_assignment_condition_version
}
