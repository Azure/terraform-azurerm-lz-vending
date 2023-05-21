# azapi_resource.routingintent ensure routing intent is being configured for a firewall that belongs to the subscription configured in the caller
resource "azapi_resource" "routingintent" {
  for_each  = { for k, v in var.virtual_hubs : k => v if length(regexall(var.subscription_id, v.vhub_firewall_resource_id)) > 0 }
  type      = "Microsoft.Network/virtualHubs/routingIntent@2022-07-01"
  name      = "${reverse(split("/", each.value.vhub_firewall_resource_id))[0]}_RoutingIntent"
  parent_id = each.value.vwan_hub_resource_id
  body = jsonencode({
    properties = {
      routingPolicies = (each.value.intent_based_internet_traffic_enabled && each.value.intent_based_private_traffic_enabled ?
        [merge(local.public_routing_policy, { nextHop = each.value.vhub_firewall_resource_id }), merge(local.private_routing_policy, { nextHop = each.value.vhub_firewall_resource_id })] :
        each.value.intent_based_internet_traffic_enabled && !each.value.intent_based_private_traffic_enabled ?
        [merge(local.public_routing_policy, { nextHop = each.value.vhub_firewall_resource_id })] :
        !each.value.intent_based_internet_traffic_enabled && each.value.intent_based_private_traffic_enabled ?
      [merge(local.private_routing_policy, { nextHop = each.value.vhub_firewall_resource_id })] : [])
    }
  })
}
