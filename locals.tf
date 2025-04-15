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

  # This virtual_networks varialbe is used internally to consume the mapped subnet properties for dependencies on resources such as 
  # route tables today but at some point network security groups as well.
  virtual_networks = [
  for vnet_k, vnet_v in var.virtual_networks : {
    name                    = vnet_v.name
    address_prefixes        = vnet_v.address_prefixes
    resource_group_name     = vnet_v.resource_group_name
    location                = vnet_v.location
    dns_servers             = vnet_v.dns_servers

    flow_timeout_in_minutes = vnet_v.flow_timeout_in_minutes

    ddos_protection_enabled = vnet_v.ddos_protection_enabled
    ddos_protection_plan_id = vnet_v.ddos_protection_plan_id

    subnets                 = local.virtual_network_subnets_map[vnet_k].subnets
  }
]


# virtual_network_subnets_map is a map with the virtual network name as the key and the subnet properties re-mapped
# to support the reference of route tables created by the LZ-Vending route table module. This type of reference will
# also be required with network security group created is added to this accelerator.
  virtual_network_subnets_map = {
    for item in flatten(
      [
        for vnet_k, vnet_v in var.virtual_networks : [
          for subnet_k, subnet_v in vnet_v.subnets : {
            vnet_key  = vnet_k
            subnet_key = subnet_k
            subnets = {
              name                                          = subnet_v.name
              address_prefixes                              = subnet_v.address_prefixes
              nat_gateway                                   = subnet_v.nat_gateway
              network_security_group                        =  subnet_v.network_security_group
              private_endpoint_network_policies             = subnet_v.private_endpoint_network_policies
              private_link_service_network_policies_enabled = subnet_v.private_link_service_network_policies_enabled
              route_table                                   = try(coalesce(try(subnet_v.route_table.id, null), try(local.virtual_network_subnet_route_table_available_resource_ids[try(subnet_v.route_table.name_reference, null)]), null), null)
              default_outbound_access_enabled               = subnet_v.default_outbound_access_enabled
              service_endpoints                             = subnet_v.service_endpoints
              service_endpoint_policies                     = subnet_v.service_endpoint_policies
              delegation                                    = try(subnet_v.delegation, null)
            }
          }
        ]
      ]
    ) : "${item.vnet_key}" => item.subnets
  }


  # virtual_network_subnet_route_table_available_resource_ids is a map of route table names and resource ids.
  # The need for this is within the LZ-Vending module there route table may be created but the user would not know
  # the resource id in advance, in such case they could specify the name in the `name_reference` property of the
  # virtual network subnet's route table object.

  virtual_network_subnet_route_table_available_resource_ids = { for rt_k, rt_v in module.routetable : element(split("/", rt_v.route_table_resource_id), -1) => rt_v.route_table_resource_id }


  # resource_group_ids is a map of resource groups created, if the module has been enabled.
  # This is used in the outputs.tf file to return the resource group ids.
  virtual_network_resource_group_ids = var.virtual_network_enabled ? module.virtualnetwork[0].resource_group_resource_ids : {}
  # virtual_networks_merged is a map of virtual networks created, if the module has been enabled.
  # This is used in the outputs.tf file to return the virtual network resource ids.
  virtual_network_resource_ids = var.virtual_network_enabled ? module.virtualnetwork[0].virtual_network_resource_ids : {}

  # route_table_routes is a list of objects containing the routes converted from a map to a list information after the user managed identities are created, if the module has been enabled.
  # since var.user_managed_identities is a map that contains the role assignments maps, we need to use a for loop to extract the values from the nested map.
  route_tables = {
    for item in flatten(
      [
        for rt_k, rt_v in var.route_tables : {
          rt_key = rt_k
          name = rt_v.name
          location = rt_v.location
          resource_group_name = rt_v.resource_group_name
          bgp_route_propagation_enabled = rt_v.bgp_route_propagation_enabled
          tags = rt_v.tags
          routes = [for k, v in rt_v.routes : v]
        }
      ]
    ) : item.rt_key => item
  }
}
