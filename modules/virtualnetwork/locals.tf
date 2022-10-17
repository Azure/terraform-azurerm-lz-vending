locals {
  # subscription_resource_id is the ARM resource ID of the supplied subscription id.
  subscription_resource_id = "/subscriptions/${var.subscription_id}"
}

locals {
  # virtual_networks_data contains additional calculated data required to create the virtual networks
  virtual_networks_data = {
    for k, v in var.virtual_networks : k => {
      # virtual_network_resource_id is the Azure resource id of the virtual network.
      # Although we could use the azapi_resource.vnet.id, this is not known until after the resource is created.
      # Therefore we construct this using the input vars to improve known attributes in unit testing.
      virtual_network_resource_id = "${local.subscription_resource_id}/resourceGroups/${v.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${v.name}"

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

      vwan_propagated_routetables_labels = coalescelist(v.vwan_propagated_routetables_labels, ["default"])
    }
  }

  # virtual_networks_mesh_peering_map is the data required to create the mesh peerings.
  # That is those peerings between the virtual networks that are supplied in the var.virtual_networks variable
  virtual_networks_mesh_peering_list = flatten([
    for k_src, v_src in local.virtual_networks_data : [
      for k_dst, v_dst in local.virtual_networks_data : {
        source_key              = k_src
        destination_key         = k_dst
        name                    = "peer-${uuidv5("url", v_dst.virtual_network_resource_id)}"
        this_resource_id        = azapi_resource.vnet[k_src].id
        remote_resource_id      = v_dst.virtual_network_resource_id
        allow_forwarded_traffic = v_src.mesh_peering_allow_forwarded_traffic
      } if var.virtual_networks[k_dst].mesh_peering_enabled && k_src != k_dst
    ] if var.virtual_networks[k_src].mesh_peering_enabled
  ])
}

locals {
  resource_group_data = toset([
    for k, v in var.virtual_networks : {
      name      = v.resource_group_name
      location  = v.location
      lock      = v.resource_group_lock_enabled
      lock_name = v.resource_group_lock_name
      tags      = v.resource_group_tags
    } if v.resource_group_creation_enabled
  ])
}
