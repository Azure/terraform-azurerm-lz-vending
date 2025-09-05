locals {
  # subscription_resource_id is the ARM resource ID of the supplied subscription id.
  subscription_resource_id = "/subscriptions/${var.subscription_id}"
}

locals {
  # hub_peering_map is a map of the virtual network hub peerings for those networks
  # with hub peering enabled
  hub_peering_map = {
    for k, v in var.virtual_networks : k => {
      # Peering this network to the remote network
      outbound = {
        name               = coalesce(v.hub_peering_name_tohub, "peer-${uuidv5("url", v.hub_network_resource_id)}")
        this_resource_id   = module.virtual_networks[k].resource_id
        remote_resource_id = v.hub_network_resource_id
        options            = v.hub_peering_options_tohub
      },
      # Peering the remote network to this network
      inbound = {
        name               = coalesce(v.hub_peering_name_fromhub, "peer-${uuidv5("url", local.virtual_network_resource_ids[k])}")
        this_resource_id   = v.hub_network_resource_id
        remote_resource_id = module.virtual_networks[k].resource_id
        options            = v.hub_peering_options_fromhub
      }
      peering_direction = contains(local.valid_peering_directions, coalesce(lower(v.hub_peering_direction), local.peering_direction_both)) ? coalesce(lower(v.hub_peering_direction), local.peering_direction_both) : local.peering_direction_both
    } if v.hub_peering_enabled
  }
  # peering direction constansts
  peering_direction_both    = "both"
  peering_direction_fromhub = "fromhub"
  peering_direction_tohub   = "tohub"
  # allowed values for peering direction
  valid_peering_directions = [local.peering_direction_tohub, local.peering_direction_fromhub, local.peering_direction_both]
  # virtual_networks_resource_ids is a map of the virtual network resource IDs.
  # we construct these to better enable testing of values in the plan
  virtual_network_resource_ids = {
    for k, v in var.virtual_networks : k => "${local.subscription_resource_id}/resourceGroups/${v.resource_group_name}/providers/Microsoft.Network/virtualNetworks/${v.name}"
  }
  # virtual_networks_mesh_peering_map is the data required to create the mesh peerings.
  # That is those peerings between the virtual networks that are supplied in the var.virtual_networks variable
  virtual_networks_mesh_peering_list = flatten([
    for k_src, v_src in local.virtual_network_resource_ids : [
      for k_dst, v_dst in local.virtual_network_resource_ids : {
        source_key              = k_src
        destination_key         = k_dst
        name                    = "peer-${uuidv5("url", v_dst)}"
        this_resource_id        = module.virtual_networks[k_src].resource_id
        remote_resource_id      = v_dst
        allow_forwarded_traffic = var.virtual_networks[k_src].mesh_peering_allow_forwarded_traffic
      } if var.virtual_networks[k_dst].mesh_peering_enabled && k_src != k_dst
    ] if var.virtual_networks[k_src].mesh_peering_enabled
  ])
  vwan_propagated_noneroutetables_resource_ids = {
    for k, v in var.virtual_networks : k => coalescelist(
      [
        for i in v.vwan_propagated_routetables_resource_ids : { id = i }
      ],
      [
        { id = "${v.vwan_hub_resource_id}/hubRouteTables/noneRouteTable" }
      ]
    ) if v.vwan_connection_enabled
  }
  # vwan_propagated_routetables_labels is a map of the virtual network vwan propagated routetables labels
  # for each virtual network that enabled for vwan connectivity.
  vwan_propagated_routetables_labels = {
    for k, v in var.virtual_networks : k => coalescelist(
      v.vwan_propagated_routetables_labels,
      ["default"]
    ) if v.vwan_connection_enabled
  }
  # vwan_propagated_routetables_resource_ids is a map of the virtual network vwan propagated routetable ids
  # for each virtual network that enabled for vwan connectivity.
  vwan_propagated_routetables_resource_ids = {
    for k, v in var.virtual_networks : k => coalescelist(
      [
        for i in v.vwan_propagated_routetables_resource_ids : { id = i }
      ],
      [
        { id = "${v.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable" }
      ]
    ) if v.vwan_connection_enabled
  }
}

locals {
  vhubconnection_body_properties = {
    for k, v in var.virtual_networks : k =>
    merge({
      enableInternetSecurity = v.vwan_security_configuration.secure_internet_traffic
      remoteVirtualNetwork = {
        id = local.virtual_network_resource_ids[k]
      }
      },
      # Only supply routingConfiguration if routing_intent_enabled is set to false
      v.vwan_security_configuration.routing_intent_enabled ? {} : {
        routingConfiguration = {
          associatedRouteTable = {
            id = v.vwan_associated_routetable_resource_id != null ? v.vwan_associated_routetable_resource_id : "${v.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable"
          }
          propagatedRouteTables = {
            ids    = v.vwan_security_configuration.secure_private_traffic ? local.vwan_propagated_noneroutetables_resource_ids[k] : local.vwan_propagated_routetables_resource_ids[k]
            labels = v.vwan_security_configuration.secure_private_traffic ? ["none"] : local.vwan_propagated_routetables_labels[k]
          }
        }
    }) if v.vwan_connection_enabled
  }
}
