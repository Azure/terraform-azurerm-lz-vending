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

resource "azapi_resource" "peerings" {
  for_each  = local.virtual_network_peering_map
  type      = "Microsoft.Network/virtualNetworks/virtualNetworkPeerings@2021-08-01"
  parent_id = each.value.this_resource_id
  name      = each.key
  body = jsonencode({
    properties = {
      remoteVirtualNetwork = {
        id = each.value.that_resource_id
      }
      allowVirtualNetworkAccess = true
      allowForwardedTraffic     = each.value.this_resource_id == azapi_resource.vnet.id ? true : false
      allowGatewayTransit       = each.value.this_resource_id == azapi_resource.vnet.id ? false : true
      useRemoteGateways         = each.value.this_resource_id == azapi_resource.vnet.id ? true : false
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
