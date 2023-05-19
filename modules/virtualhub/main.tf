resource "azapi_resource" "routingintent" {
  for_each  = { for k, v in var.virtual_hubs : k => v }
  type = "Microsoft.Network/virtualHubs/routingIntent@2022-07-01"
  name = "vhub-ukwest_RoutingIntent"
  parent_id = each.value.virtual_hub_id
  body = jsonencode({
    properties = {
      routingPolicies = [
        {
          destinations = [
            "Internet"
          ]
          name = "PublicTraffic"
          nextHop = each.value.intent_based_routing_next_hop_firewall
        }
      ]
    }
  })
}