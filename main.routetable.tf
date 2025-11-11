# route table submodule, disabled by default
# Will create a route table, and optionally routes
module "routetable" {
  source   = "./modules/routetable"
  for_each = var.route_table_enabled ? local.route_tables : {}

  location                      = each.value.location
  name                          = each.value.name
  parent_id                     = "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name}"
  bgp_route_propagation_enabled = each.value.bgp_route_propagation_enabled
  routes                        = each.value.routes
  tags                          = each.value.tags
}
