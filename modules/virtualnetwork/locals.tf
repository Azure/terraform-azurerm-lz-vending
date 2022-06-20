locals {
  # subscription_resource_id is the ARM resource ID of the supplied subscription id
  subscription_resource_id = "/subscriptions/${var.subscription_id}"

  hub_network_uuidv5 = uuidv5("url", var.hub_network_resource_id)

  # virtual_network_peering_map is the data required to create the two vnet peerings
  virtual_network_peering_map = var.hub_network_resource_id != "" ? {
    peer-hub-vnet = {
      this_resource_id = azapi_resource.vnet.id
      that_resource_id = var.hub_network_resource_id
    },
    "peer-${local.hub_network_uuidv5}" = {
      this_resource_id = var.hub_network_resource_id
      that_resource_id = azapi_resource.vnet.id
    }
  } : {}
}
