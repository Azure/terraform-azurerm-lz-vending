# azapi_resource.rg is the resource group that the virtual network will be created in
# the module will create as many as is required by the var.virtual_networks input variable
resource "azapi_resource" "rg" {
  for_each  = { for i in local.resource_group_data : i.name => i }
  parent_id = local.subscription_resource_id
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = each.key
  location  = each.value.location
  tags      = each.value.tags
}

# azapi_resource.rg_lock is an optional resource group lock that can be used
# to prevent accidental deletion.
resource "azapi_resource" "rg_lock" {
  for_each  = { for i in local.resource_group_data : i.name => i if i.lock }
  type      = "Microsoft.Authorization/locks@2017-04-01"
  parent_id = azapi_resource.rg[each.key].id
  name      = coalesce(each.value.lock_name, substr("lock-${each.key}", 0, 90))
  body = jsonencode({
    properties = {
      level = "CanNotDelete"
    }
  })
  depends_on = [
    azapi_resource.vnet,
    azapi_update_resource.vnet,
    azapi_resource.peering_hub_outbound,
    azapi_resource.peering_hub_inbound,
    azapi_resource.peering_mesh,
    azapi_resource.vhubconnection,
  ]
}

# azapi_resource.vnet are the virtual networks that will be created
# lifecycle ignore changes to the body to prevent subnets being deleted
# see #45 for more information
resource "azapi_resource" "vnet" {
  for_each  = var.virtual_networks
  parent_id = "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name}"
  type      = "Microsoft.Network/virtualNetworks@2021-08-01"
  name      = each.value.name
  location  = coalesce(each.value.location, var.location)
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = each.value.address_space
      }
      dhcpOptions = {
        dnsServers = each.value.dns_servers
      }
    }
  })
  tags = each.value.tags
  lifecycle {
    ignore_changes = [body, tags]
  }
  depends_on = [
    azapi_resource.rg,
  ]
}

# azapi_update_resource.vnet are the virtual networks that will be created
# This is a workaround for #45 to allow updates to the virtual network
# without deleting the subnets created elsewhere
resource "azapi_update_resource" "vnet" {
  for_each    = var.virtual_networks
  resource_id = azapi_resource.vnet[each.key].id
  type        = "Microsoft.Network/virtualNetworks@2021-08-01"
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = each.value.address_space
      }
      dhcpOptions = {
        dnsServers = each.value.dns_servers
      }
    }
    tags = each.value.tags
  })
}

# azapi_resource.peering_hub_outbound creates one-way peering from the spoke to the supplied hub virtual network.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_outbound" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.hub_peering_enabled }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = local.hub_peering_map[each.key]["outbound"].this_resource_id
  name      = local.hub_peering_map[each.key]["outbound"].name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = local.hub_peering_map[each.key]["outbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = false
      useRemoteGateways         = each.value.hub_peering_use_remote_gateways
    }
  })
}

# azapi_resource.peering_hub_inbound creates one-way peering from the supplied hub network to the spoke.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering_hub_inbound" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.hub_peering_enabled }
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = local.hub_peering_map[each.key]["inbound"].this_resource_id
  name      = local.hub_peering_map[each.key]["inbound"].name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = local.hub_peering_map[each.key]["inbound"].remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = true
      useRemoteGateways         = false
    }
  })
}

# azapi_resource.peering_mesh creates mesh peerings between the supplied var.virtual_networks.
# They are created if the boolean mesh_peering_enabled is set to true on more than one network.
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
  parent_id = each.value.vwan_hub_resource_id
  name      = coalesce(each.value.vwan_connection_name, "vhc-${uuidv5("url", azapi_resource.vnet[each.key].id)}")
  body = jsonencode({
    properties = {
      enableInternetSecurity = each.value.security_configuration.secure_internet_traffic
      remoteVirtualNetwork = {
        id = local.virtual_network_resource_ids[each.key]
      }
      routingConfiguration = {
        associatedRouteTable = {
          id = each.value.vwan_associated_routetable_resource_id != "" ? each.value.vwan_associated_routetable_resource_id : "${each.value.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable"
        }
        propagatedRouteTables = {
          ids    = local.vwan_propagated_routetables_resource_ids[each.key]
          labels = local.vwan_propagated_routetables_labels[each.key]
        }
      }
    }
  })
}

# azapi_update_resource.vhubdefaultroutetable
resource "azapi_update_resource" "vhubdefaultroutetable" {
  for_each  = { for k, v in var.virtual_networks : k => v if v.security_configuration.secure_internet_traffic }
  type      = "Microsoft.Network/virtualHubs/hubRouteTables@2022-07-01"
  parent_id = each.value.vwan_hub_resource_id
  name = "defaultRouteTable"
  body = jsonencode({
    properties = {
      labels = [
        "default"
      ]
      routes = [
        {
          "name"            = "public_traffic"
          "destinationType" = "CIDR"
          "destinations" = [
            "0.0.0.0/0"
          ]
          "nextHopType" = "ResourceId"
          "nextHop"     = each.value.security_configuration.next_hop
        }
      ]
    }
  })
}