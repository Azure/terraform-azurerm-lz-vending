locals {
  # subscription_id is the id of the subscription into which resources will be created.
  # We pick the created sub id first, if it exists, otherwise we pick the subscription_id variable.
  subscription_id = coalesce(local.subscription_module_output_subscription_id, var.subscription_id)
  # subscription_module_output_subscription_id is either the output of the subscription module,
  # or if disabled, a null.
  # Needed to avoid errors in local.subscription_id when referencing a module instance that does not exists.
  subscription_module_output_subscription_id = try(module.subscription[0].subscription_id, null)
  # subscription_module_output_subscription_id is either the output of the subscription module,
  # or if disabled, a null.
  # Needed to avoid errors in local.subscription_id when referencing a module instance that does not exists.
  subscription_module_output_subscription_resource_id = try(module.subscription[0].subscription_resource_id, null)
  # subscription_resource_id is the Azure resource id of the subscription into which resources will be created.
  # We use the created sub resource id first, if it exists, otherwise we pick the subscription_id variable.
  # If this is blank then the subscription submodule is disabled an no subscription id has been supplied as an input variable.
  subscription_resource_id = coalesce(local.subscription_module_output_subscription_resource_id, local.supplied_subscription_resource_id)
  # subscription_resource_id is the Azure resource id of the subscription id that was supplied in the input variables.
  # If var.subscription_id is empty, then we will return an empty string so that we can correctly coalesce the subscription_resource_id output.
  supplied_subscription_resource_id = var.subscription_id == null ? null : "/subscriptions/${var.subscription_id}"
  # umi_client_ids is a map of client ids for the user managed identities created, if the module has been enabled.
  # This is used in the outputs.tf file to return the umi client ids.
  umi_client_ids = var.umi_enabled ? { for k, v in module.usermanagedidentity : k => v.client_id } : {}
  # umi_principal_ids is a map of principal ids for the user managed identities created, if the module has been enabled.
  # This is used in the outputs.tf file to return the umi principal ids.
  umi_principal_ids = var.umi_enabled ? { for k, v in module.usermanagedidentity : k => v.principal_id } : {}
  # umi_resource_ids is a map of user managed identities created, if the module has been enabled.
  # This is used in the outputs.tf file to return the umi resource ids.
  umi_resource_ids = var.umi_enabled ? { for k, v in module.usermanagedidentity : k => v.resource_id } : {}
  # umi_tenant_ids is a map of tenant ids for the user managed identities created, if the module has been enabled.
  # This is used in the outputs.tf file to return the umi tenant ids. Since there can my duplicate tenant ids,
  # we should only return unique values.
  umi_tenant_ids = var.umi_enabled ? { for k, v in module.usermanagedidentity : k => v.tenant_id } : {}
  # user_managed_identity_role_assignments is a list of objects containing the identity information after the user managed identities are created, if the module has been enabled.
  # since var.user_managed_identities is a map that contains the role assignments maps, we need to use a for loop to extract the values from the nested map.
  # using https://github.com/Azure/terraform-robust-module-design/blob/main/nested_maps/flatten_nested_map/main.tf as a reference.
  user_managed_identity_role_assignments = {
    for item in flatten(
      [
        for umi_k, umi_v in var.user_managed_identities : [
          for role_k, role_v in umi_v.role_assignments : {
            umi_key  = umi_k
            role_key = role_k
            role_assignment = {
              principal_id              = module.usermanagedidentity[umi_k].principal_id
              definition                = role_v.definition
              scope                     = "${local.subscription_resource_id}${role_v.relative_scope}"
              condition                 = role_v.condition
              condition_version         = role_v.condition_version
              principal_type            = role_v.principal_type
              definition_lookup_enabled = role_v.definition_lookup_enabled
            }
          }
        ]
      ]
    ) : "${item.umi_key}/${item.role_key}" => item.role_assignment
  }
  # resource_group_ids is a map of resource groups created, if the module has been enabled.
  # This is used in the outputs.tf file to return the resource group ids.
  virtual_network_resource_group_ids = var.virtual_network_enabled ? module.virtualnetwork[0].resource_group_resource_ids : {}
  # virtual_networks_merged is a map of virtual networks created, if the module has been enabled.
  # This is used in the outputs.tf file to return the virtual network resource ids.
  virtual_network_resource_ids = var.virtual_network_enabled ? module.virtualnetwork[0].virtual_network_resource_ids : {}
}
