# The roleassignments module creates role assignments from the data
# supplied in the var.role_assignments variable
module "roleassignment" {
  source = "./modules/roleassignment"
  depends_on = [
    module.resourcegroup,
    module.subscription,
    module.usermanagedidentity,
    module.virtualnetwork,
  ]
  for_each                          = { for k, v in var.role_assignments : k => v if var.role_assignment_enabled }
  role_assignment_principal_id      = each.value.principal_id
  role_assignment_definition        = each.value.definition
  role_assignment_scope             = "${local.subscription_resource_id}${each.value.relative_scope}"
  role_assignment_condition         = each.value.condition
  role_assignment_condition_version = each.value.condition_version
}

# The roleassignments_umi module creates role assignments from the data
# supplied in the var.user_managed_identities object role_assignments property
module "roleassignment_umi" {
  source = "./modules/roleassignment"
  depends_on = [
    module.resourcegroup,
    module.subscription,
    module.usermanagedidentity,
    module.virtualnetwork,
  ]
  for_each = local.user_managed_identity_role_assignments

  role_assignment_principal_id      = each.value.principal_id
  role_assignment_definition        = each.value.definition
  role_assignment_scope             = each.value.scope
  role_assignment_condition         = each.value.condition
  role_assignment_condition_version = each.value.condition_version
}
