# The roleassignments module creates role assignments from the data
# supplied in the var.role_assignments variable
module "roleassignment" {
  source                                    = "./modules/roleassignment"
  for_each                                  = { for k, v in var.role_assignments : k => v if var.role_assignment_enabled }
  role_assignment_principal_id              = each.value.principal_id
  role_assignment_definition                = local.role_assignments_definitions[each.key]
  role_assignment_scope                     = each.value.resource_group_scope_key != null ? module.resourcegroup[each.value.resource_group_scope_key].resource_group_resource_id : "${local.subscription_resource_id}${each.value.relative_scope}"
  role_assignment_condition                 = each.value.condition
  role_assignment_condition_version         = each.value.condition_version
  role_assignment_principal_type            = each.value.principal_type
  role_assignment_definition_lookup_enabled = each.value.definition_lookup_enabled
  role_assignment_use_random_uuid           = each.value.use_random_uuid
  enable_telemetry                          = !var.disable_telemetry
}

# The roleassignments_umi module creates role assignments from the data
# supplied in the var.user_managed_identities object role_assignments property
module "roleassignment_umi" {
  source   = "./modules/roleassignment"
  for_each = local.user_managed_identity_role_assignments

  role_assignment_principal_id              = each.value.principal_id
  role_assignment_definition                = each.value.definition
  role_assignment_scope                     = each.value.scope
  role_assignment_condition                 = each.value.condition
  role_assignment_condition_version         = each.value.condition_version
  role_assignment_principal_type            = each.value.principal_type
  role_assignment_definition_lookup_enabled = each.value.definition_lookup_enabled
  role_assignment_use_random_uuid           = each.value.use_random_uuid
  enable_telemetry                          = !var.disable_telemetry
  retry = {
    error_message_regex = [
      "PrincipalNotFound",
    ]
  }
}
