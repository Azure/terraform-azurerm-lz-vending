module "roleassignment" {
  source                       = "./modules/roleassignment"
  for_each                     = local.role_assignments_map
  role_assignment_principal_id = each.value.principal_id
  role_assignment_definition   = each.value.definition
  role_assignment_scope        = each.value.scope
}
