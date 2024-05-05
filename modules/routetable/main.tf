resource "azapi_resource" "route_table" {
  for_each = var.route_tables

  type      = "Microsoft.Network/routeTables@2023-04-01"
  parent_id = "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name}"
  name      = each.value.name
  location  = each.value.location
  body = {
    properties = {
      disableBgpRoutePropagation = try(each.value.disable_bgp_route_propagation, false)
      routes = each.value.routes != null ? [
        for r in each.value.routes : {
          name = r.name
          properties = {
            addressPrefix    = r.address_prefix
            nextHopIpAddress = r.next_hop_in_ip_address
            nextHopType      = r.next_hop_type
          }
        }
      ] : null
    }
  }
  schema_validation_enabled = true
  response_export_values    = ["*"]
  tags                      = each.value.tags
}
