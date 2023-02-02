# The roleassignments module creates role assignments from the data
# supplied in the var.role_assignments variable.
# var.role_assignment_enabled must also be set to true.
module "roleassignment" {
  source = "./modules/roleassignment"
  depends_on = [
    module.subscription,
    module.virtualnetwork,
  ]
  for_each                     = var.role_assignments
  role_assignment_principal_id = each.value.principal_id
  role_assignment_definition   = each.value.definition
  role_assignment_scope        = each.value.scope
}
