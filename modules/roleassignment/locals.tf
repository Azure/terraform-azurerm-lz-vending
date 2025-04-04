locals {
  # This local represents the final role definition resource id as required by the roleAssignments resource.
  # It will be either the role definition resource id supplied in var.role_assignment_definition,
  # or the role definition resource id that is looked up based on the role name supplied in var.role_assignment_definition.
  # If the the role definition name cannot the value will be null.
  role_assignment_definition_id = can(regex("/providers/Microsoft.Authorization/roleDefinitions", var.role_assignment_definition)) ? var.role_assignment_definition : local.role_definitions_role_name_to_resource_id
  role_assignments_role_name_to_resource_id = var.role_assignment_definition_lookup_enabled ? {
    for res in data.azapi_resource_list.role_definitions[0].output.results : res.role_name => res.id
  } : {}
  role_definitions_role_name_to_resource_id = var.role_assignment_definition_lookup_enabled ? lookup(
    local.role_assignments_role_name_to_resource_id,
    var.role_assignment_definition,
    null
  ) : null
  role_assignment_properties = merge({
    principalId      = var.role_assignment_principal_id
    roleDefinitionId = local.role_assignment_definition_id
    condition        = var.role_assignment_condition
    conditionVersion = var.role_assignment_condition_version
    },
    var.role_assignment_principal_type != null ? {
      principalType = var.role_assignment_principal_type
  } : {})
}
