# module.networksecuritygroup uses the local submodule to create
# as many network security groups as is required by the var.network_security_groups input variable
# and any nested security rules within the network security group.
module "networksecuritygroup" {

  source = "./modules/networksecuritygroup"

  for_each = { for nsg_k, nsg_v in var.network_security_groups : nsg_k => nsg_v if var.network_security_group_enabled }

  name     = each.value.name
  location = coalesce(each.value.location, var.location)
  parent_id = coalesce(
    can(module.resourcegroup[each.value.resource_group_key].resource_group_resource_id) ? module.resourcegroup[each.value.resource_group_key].resource_group_resource_id : null,
    each.value.resource_group_name_existing != null ? "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name_existing}" : null
  )
  tags = each.value.tags

  security_rules = each.value.security_rules
}
