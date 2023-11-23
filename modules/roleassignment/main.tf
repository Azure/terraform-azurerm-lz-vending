resource "azurerm_role_assignment" "this" {
  scope                = var.role_assignment_scope
  principal_id         = var.role_assignment_principal_id
  role_definition_id   = local.role_assignment_definition_id
  role_definition_name = local.role_assignment_definition_name
  condition            = var.role_assignment_condition
  condition_version    = var.role_assignment_condition_version
}
