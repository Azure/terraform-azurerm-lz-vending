locals {
  # This local represents the final role definition resource id as required by the roleAssignments resource.
  # It will be either the role definition resource id supplied in var.role_assignment_definition,
  # or the role definition resource id that is looked up based on the role name supplied in var.role_assignment_definition.
  # If the the role definition name cannot the value will be null.
  is_definition_resource_id           = can(regex("(?i)/providers/Microsoft.Authorization/roleDefinitions", var.role_assignment_definition))
  role_definition_name_to_resource_id = lookup(module.role_definitions.role_definition_name_to_resource_id, var.role_assignment_definition, null)
  role_definition_id                  = local.is_definition_resource_id ? var.role_assignment_definition : local.role_definition_name_to_resource_id

  role_assignment_properties = merge({
    principalId      = var.role_assignment_principal_id
    roleDefinitionId = local.role_definition_id
    condition        = var.role_assignment_condition
    conditionVersion = var.role_assignment_condition_version
    },
    var.role_assignment_principal_type != null ? {
      principalType = var.role_assignment_principal_type
  } : {})
}
