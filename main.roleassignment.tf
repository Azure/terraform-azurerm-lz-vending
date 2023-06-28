# The roleassignments module creates role assignments from the data
# supplied in the var.role_assignments variable
module "roleassignment" {
  source = "./modules/roleassignment"
  depends_on = [
    module.subscription,
    module.virtualnetwork,
    module.resourcegroup,
    module.usermanagedidentity
  ]
  for_each                     = { for k, v in var.role_assignments : k => v if var.role_assignment_enabled }
  role_assignment_principal_id = each.value.principal_id
  role_assignment_definition   = each.value.definition
  role_assignment_scope        = "${local.subscription_resource_id}${each.value.relative_scope}"
}

# The roleassignments_umi module creates role assignments from the data
# supplied in the var.uni_role_assignments variable
module "roleassignment_umi" {
  source = "./modules/roleassignment"
  depends_on = [
    module.subscription,
    module.virtualnetwork,
    module.resourcegroup,
    module.usermanagedidentity
  ]
  for_each                     = { for k, v in var.umi_role_assignments : k => v if var.umi_enabled }
  role_assignment_principal_id = one(module.usermanagedidentity).principal_id
  role_assignment_definition   = each.value.definition
  role_assignment_scope        = "${local.subscription_resource_id}${each.value.relative_scope}"
}
