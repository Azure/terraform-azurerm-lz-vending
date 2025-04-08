data "azapi_resource_list" "role_definitions" {
  count = var.role_assignment_definition_lookup_enabled ? 1 : 0

  parent_id = var.role_assignment_scope
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  response_export_values = {
    results = "value[].{id: id, role_name: properties.roleName}"
  }
}

resource "azapi_resource" "this" {
  type = "Microsoft.Authorization/roleAssignments@2022-04-01"
  body = {
    properties = local.role_assignment_properties
  }
  name      = uuidv5("url", "${var.role_assignment_scope}${var.role_assignment_principal_id}${local.role_assignment_definition_id}")
  parent_id = var.role_assignment_scope

  lifecycle {
    precondition {
      condition     = local.role_assignment_definition_id != null
      error_message = "In `var.role_assignment_definition` - either supply the role assignment definition resource id or a valid role assignment definition name (and make sure that role definition lookup is enabled)."
    }
  }
}
