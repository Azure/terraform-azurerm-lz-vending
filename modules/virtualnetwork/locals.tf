locals {
  # subscription_resource_id is the ARM resource ID of the supplied subscription id.
  subscription_resource_id = "/subscriptions/${var.subscription_id}"

  # hub_network_uuidv5 generates a unique, but predictable uuid based on the
  # resource id of the created virtual network.
  hub_network_uuidv5  = uuidv5("url", var.hub_network_resource_id)
  this_network_uuidv5 = uuidv5("url", azapi_resource.vnet.id)

  # virtual_network_peering_map is the data required to create the two vnet peerings.
  # If the supplied hub virtual network is an empty string, the map will be empty,
  # resulting in no resources being created.
  virtual_network_peering_map = var.virtual_network_peering_enabled ? {
    # Peering this network to the remote network
    "outbound" = {
      name               = "peer-${local.hub_network_uuidv5}"
      this_resource_id   = azapi_resource.vnet.id
      remote_resource_id = var.hub_network_resource_id
    },
    # Peering the remote network to this network
    "inbound" = {
      name               = "peer-${local.this_network_uuidv5}"
      this_resource_id   = var.hub_network_resource_id
      remote_resource_id = azapi_resource.vnet.id
    }
  } : {}

  # virtual_network_resource_id is the Azure resource id of the virtual network.
  # Although we could use the azapi_resource.vnet.id, this is not known until after the resource is created.
  # Therefore we construct this using the input vars to improve known attributes in unit testing.
  virtual_network_resource_id = "/subscriptions/${var.subscription_id}/resourceGroups/${var.virtual_network_resource_group_name}/providers/Microsoft.Network/virtualNetworks/${var.virtual_network_name}"

  # vhub_connection_map is the data required to create the virtual wan hub connection.
  # If the supplied vwan hub is an empty string, the set will be empty,
  # resulting on no resource being created.
  vhub_connection_set = var.virtual_network_vwan_connection_enabled ? toset([
    "this"
  ]) : toset([])

  # vwan_propagated_routetables_resource_ids generates the routetable resource ids for the vhub connection
  # if not specified by the input variable, it will be the default routetable for the vwan hub.
  vwan_propagated_routetables_resource_ids = coalescelist(
    [
      for i in var.virtual_network_vwan_propagated_routetables_resource_ids : { id = i }
    ],
    [
      { id = "${var.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable" }
    ]
  )

  # vwan_propagated_labels generates the propagated route labels for the vhub connection
  # if not specified by the input variable, it will be set to a list with a single item: `default`.
  vwan_propagated_labels = coalescelist(var.virtual_network_vwan_propagated_routetables_labels, ["default"])
}

locals {
  # virtual_networks_data contains additional calculated data required to create the virtual networks
  virtual_networks_data = {
    for k, v in var.virtual_networks : k => {
      # virtual_network_resource_id is the Azure resource id of the virtual network.
      # Although we could use the azapi_resource.vnet.id, this is not known until after the resource is created.
      # Therefore we construct this using the input vars to improve known attributes in unit testing.
      virtual_network_resource_id = "/subscriptions/${v.subscription_id}/resourceGroups/${v.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${v.name}"

      # virtual_network_peering_map is the data required to create the two vnet peerings.
      # If the supplied hub virtual network is an empty string, the map will be empty,
      # resulting in no resources being created.
      hub_peering_map = v.hub_peering_enabled ? {
        # Peering this network to the remote network
        "outbound" = {
          name               = coalesce(v.hub_peering_name_tohub, "peer-${uuidv5("url", v.hub_network_resource_id)}")
          this_resource_id   = azapi_resource.vnet[k].id
          remote_resource_id = v.hub_network_resource_id
        },
        # Peering the remote network to this network
        "inbound" = {
          name               = coalesce(v.hub_peering_name_fromhub, "peer-${uuidv5("url", azapi_resource.vnet[k].id)}")
          this_resource_id   = v.hub_network_resource_id
          remote_resource_id = azapi_resource.vnet[k].id
        }
      } : {}

      # vwan_propagated_routetables_resource_ids generates the routetable resource ids for the vhub connection
      # if not specified by the input variable, it will be the default routetable for the vwan hub.
      vwan_propagated_routetables_resource_ids = coalescelist(
        [
          for i in v.vwan_propagated_routetables_resource_ids : { id = i }
        ],
        [
          { id = "${v.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable" }
        ]
      )
    }
  }

  # virtual_networks_mesh_peering_map is the data required to create the mesh peerings.
  # That is those peerings between the virtual networks that are supplied in the var.virtual_networks variable
  virtual_networks_mesh_peering_list = flatten([
    for k_src, v_src in var.virtual_networks_data : [
      for k_dst, v_dst in var.virtual_networks_data : {
        source_key              = k_src
        destination_key         = k_dst
        name                    = "peer-${uuidv5("url", v_dst.virtual_network_resource_id)}"
        this_resource_id        = azapi_resource.vnet[k_src].id
        remote_resource_id      = v_dst.virtual_network_resource_id
        allow_forwarded_traffic = v_src.mesh_peering_allow_forwarded_traffic
      } if v_dst.mesh_peering_enabled && k_src != k_dst
     ] if v_src.mesh_peering_enabled
  ])
}
