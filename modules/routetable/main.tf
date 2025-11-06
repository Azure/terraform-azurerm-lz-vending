resource "azapi_resource" "route_table" {
  location  = var.location
  name      = var.name
  parent_id = var.parent_id
  type      = "Microsoft.Network/routeTables@2023-04-01"
  body = {
    properties = {
      disableBgpRoutePropagation = !var.bgp_route_propagation_enabled
      routes = [
        for r in var.routes : {
          name = r.name
          properties = {
            addressPrefix    = r.address_prefix
            nextHopIpAddress = r.next_hop_in_ip_address
            nextHopType      = r.next_hop_type
          }
        }
      ]
    }
  }
  tags = var.tags
}
