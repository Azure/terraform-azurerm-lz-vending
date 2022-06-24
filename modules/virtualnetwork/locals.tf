locals {
  # subscription_resource_id is the ARM resource ID of the supplied subscription id.
  subscription_resource_id = "/subscriptions/${var.subscription_id}"

  # hub_network_uuidv5 generates a unique, but predictable uuid based on the
  # resource id of the cretaed virtual network.
  hub_network_uuidv5  = uuidv5("url", var.hub_network_resource_id)
  this_network_uuidv5 = uuidv5("url", azapi_resource.vnet.id)

  # virtual_network_peering_map is the data required to create the two vnet peerings.
  # If the supplied hub virtual network is an empty string, the map will be empty,
  # resulting on no resources being created.
  virtual_network_peering_map = var.hub_network_resource_id != "" ? {
    # Peering this network to the remote network
    "peer-${local.hub_network_uuidv5}" = {
      this_resource_id   = azapi_resource.vnet.id
      remote_resource_id = var.hub_network_resource_id
    },
    # Peering the remote network to this network
    "peer-${local.this_network_uuidv5}" = {
      this_resource_id   = var.hub_network_resource_id
      remote_resource_id = azapi_resource.vnet.id
    }
  } : {}

  # vhub_connection_map is the data required to create the virtual wan hub connection.
  # If the supplied vwan hub is an empty string, the set will be empty,
  # resulting on no resource being created.
  vhub_connection_set = var.vwan_hub_resource_id != "" ? toset([
    "vhubcon-${local.hub_network_uuidv5}"
  ]) : toset([])

  # vwan_propagated_routetables_resource_ids generates the routetable resource ids for the vhub connection
  # if not specified by the input variable, it will be the default routetable for the vwan hub.
  vwan_propagated_routetables_resource_ids = length(var.virtual_network_vwan_propagated_routetables_resource_ids) > 0 ? [
    for i in var.virtual_network_vwan_propagated_routetables_resource_ids : { id = i }
  ] : [
    { id = "${var.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable" }
  ]

  # vwan_propagated_labels generates the propagated route labels for the vhub connection
  # if not specified by the input variable, it will be set to a list with a single item: `default`.
  vwan_propagated_labels = length(var.virtual_network_vwan_propagated_routetables_labels) > 0 ? var.virtual_network_vwan_propagated_routetables_labels : ["default"]
}
