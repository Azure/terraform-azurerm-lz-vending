# azapi_resource.rg is the resource group that the virtual network will be created in
resource "azapi_resource" "rg" {
  parent_id = local.subscription_resource_id
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = var.virtual_network_resource_group_name
  location  = var.virtual_network_location
}

# azapi_resource.rg_lock is an optional resource group lock that can be used
# to prevent accidental deletion.
resource "azapi_resource" "rg_lock" {
  count     = var.virtual_network_resource_lock_enabled ? 1 : 0
  parent_id = azapi_resource.rg.id
  type      = "Microsoft.Authorization/locks@2017-04-01"
  name      = substr("lock-${var.virtual_network_resource_group_name}", 0, 90)
  body = jsonencode({
    properties = {
      level = "CanNotDelete"
    }
  })
  depends_on = [
    azapi_resource.vnet,
    azapi_update_resource.vnet,
    azapi_resource.peering,
    azapi_resource.vhubconnection,
  ]
}

# azapi_resource.vnet is the virtual network that will be created
# lifecycle ignore changes to the body to prevent subnets being deleted
# see #45 for more information
resource "azapi_resource" "vnet" {
  for_each  = var.virtual_networks
  parent_id = azapi_resource.rg.id
  type      = "Microsoft.Network/virtualNetworks@2021-08-01"
  name      = each.value.name
  location  = each.value.location
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = each.value.virtual_network_address_space
      }
    }
  })
  tags = each.value.tags
  lifecycle {
    ignore_changes = [body, tags]
  }
}

# azapi_update_resource.vnet is the virtual network that will be created
# This is a workaround for #45 to allow updates to the virtual network
# without deleting the subnets created elsewhere
resource "azapi_update_resource" "vnet" {
  for_each    = var.virtual_networks
  resource_id = azapi_resource.vnet[each.key].id
  type        = "Microsoft.Network/virtualNetworks@2021-08-01"
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = var.virtual_network_address_space
      }
    }
    tags = each.value.tags
  })
}

# azapi_resource.peerings creates two-way peering from the spoke to the supplied hub virtual network.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_outbound" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.hub_peering_enabled }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = local.virtual_networks_data[each.key].hub_peering_map["outbound"].this_resource_id
  name      = local.virtual_networks_data[each.key].hub_peering_map["outbound"].name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = local.virtual_networks_data[each.key].hub_peering_map["outbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = false
      useRemoteGateways         = each.value.hub_peering_use_remote_gateways
    }
  })
}

# azapi_resource.peerings creates two-way peering from the spoke to the supplied hub virtual network.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_inbound" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.hub_peering_enabled }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = local.virtual_networks_data[each.key].hub_peering_map["inbound"].this_resource_id
  name      = local.virtual_networks_data[each.key].hub_peering_map["inbound"].name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = local.virtual_networks_data[each.key].hub_peering_map["outbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = true
      useRemoteGateways         = false
    }
  })
}

# azapi_resource.peering_mesh creates mesh peerings between the supplied var.virtual_networks.
# They are created if the boolean mesh_peering_enabled is set to true.
resource "azapi_resource" "peering_mesh" {
  for_each  = { for i in local.virtual_networks_mesh_peering_list : "${i.source_key}-${i.destination_key}" => i }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2022-05-01"
  parent_id = each.value.this_resource_id
  name      = each.value.name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = each.value.remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = each.value.allow_forwarded_traffic
      allowGatewayTransit       = false
      useRemoteGateways         = false
    }
  })
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
resource "azapi_resource" "vhubconnection" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.vwan_connection_enabled }
  type      = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2021-08-01"
  parent_id = v.vwan_hub_resource_id
  name      = "vhc-${uuidv5("url", azapi_resource.vnet[each.key].id)}"
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = local.virtual_network_resource_id
      }
      routingConfiguration = {
        associatedRouteTable = {
          id = each.value.virtual_network_vwan_associated_routetable_resource_id != "" ? each.value.virtual_network_vwan_associated_routetable_resource_id : "${each.value.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable"
        }
        propagatedRouteTables = {
          ids    = local.virtual_networks_data[each.key].vwan_propagated_routetables_resource_ids
          labels = each.value.vwan_propagated_labels
        }
      }
    }
  })
}
