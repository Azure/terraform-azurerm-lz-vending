# route table submodule, disabled by default
# Will create a route table, and optionally routes
module "routetable" {
  for_each        = var.route_table_enabled ? local.route_tables : {}
  source          = "./modules/routetable"
  parent_id       = "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name}"
  bgp_route_propagation_enabled = each.value.bgp_route_propagation_enabled
  name                          = each.value.name
  location                      = each.value.location
  routes                        = each.value.routes
  tags                          = each.value.tags
}
