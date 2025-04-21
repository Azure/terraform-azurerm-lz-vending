# module.networksecuritygroup uses the local submodule to create
# as many network security groups as is required by the var.network_security_groups input variable
# and any nested security rules within the network security group.
module "networksecuritygroup" {
  
  source             = "./modules/networksecuritygroup"

  subscription_id = local.subscription_id

  for_each = { for nsg_k, nsg_v in var.network_security_groups : nsg_k => nsg_v if var.network_security_group_enabled }

  name     = each.value.name
  location = coalesce(each.value.location, var.location)
  resource_group_name = each.value.resource_group_name
  tags     = each.value.tags

  security_rules = each.value.security_rules

  depends_on = [
    module.resourcegroup
  ]
}