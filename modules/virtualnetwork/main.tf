# azapi_resource.rg is the resource group that the virtual network will be created in
resource "azapi_resource" "rg" {
  parent_id = local.subscription_resource_id
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = var.virtual_network_resource_group_name
  location  = var.virtual_network_location
}

# azapi_resource.vnet is the virtual network that will be created
resource "azapi_resource" "vnet" {
  parent_id = azapi_resource.rg.id
  type      = "Microsoft.Network/virtualNetworks@2021-08-01"
  name      = var.virtual_network_name
  location  = azapi_resource.rg.location
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = var.virtual_network_address_space
      }
    }
  })
}

# azapi_resource.peerings creates two-way peering from the spoke to the supplied hub virtual network.
# They are not created if the hub virtual network is an empty string.
resource "azapi_resource" "peering" {
  for_each  = local.virtual_network_peering_map
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = each.value.this_resource_id
  name      = each.value.name
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = each.value.remote_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = true
      allowGatewayTransit       = each.value.this_resource_id == azapi_resource.vnet.id ? false : true
      useRemoteGateways         = (each.value.this_resource_id == azapi_resource.vnet.id) && var.virtual_network_use_remote_gateways ? true : false
    }
  })
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
resource "azapi_resource" "vhubconnection" {
  for_each  = local.vhub_connection_set
  type      = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2021-08-01"
  parent_id = var.vwan_hub_resource_id
  name      = each.key
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = azapi_resource.vnet.id
      }
      routingConfiguration = {
        associatedRouteTable = {
          id = var.virtual_network_vwan_routetable_resource_id != "" ? var.virtual_network_vwan_routetable_resource_id : "${var.vwan_hub_resource_id}/hubRouteTables/defaultRouteTable"
        }
        propagatedRouteTables = {
          ids    = local.vwan_propagated_routetables_resource_ids
          labels = local.vwan_propagated_labels
        }
      }
    }
  })
}

# Subnet resources not in scope due to complexity of creation, e.g. route tables/nsgs
# resource "azapi_resource" "subnet" {
#   for_each = var.virtual_network_subnets
#   parent_id = azapi_resource.vnet.id
#   type = "Microsoft.Network/virtualNetworks/subnets@2021-08-01"
#   name = each.key
#   body = jsonencode({
#     properties = {
#       addressPrefix = each.value.address_prefix
#     }
#   })
# }
