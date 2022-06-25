terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = ">= 0.3.0"
    }
  }
}



variable "subscription_id" {
  type = string
}

variable "virtual_network_name" {
  type = string
}

variable "virtual_network_resource_group_name" {
  type = string
}

variable "virtual_network_location" {
  type = string
}

variable "virtual_network_address_space" {
  type = list(string)
}

variable "virtual_network_use_remote_gateways" {
  type = bool
}

variable "virtual_network_enable_peering" {
  type = bool
}

resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  parent_id = "/subscriptions/${var.subscription_id}"
  name      = "${var.virtual_network_name}-hub"
  location  = var.virtual_network_location
}

resource "azapi_resource" "hub" {
  type      = "Microsoft.Network/virtualNetworks@2021-08-01"
  name      = "${var.virtual_network_name}-hub"
  parent_id = azapi_resource.rg.id
  location  = var.virtual_network_location
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = [
          "192.168.0.0/23"
        ]
      }
      subnets = [
        {
          name = "default"
          properties = {
            addressPrefix = "192.168.0.0/24"
          }
        },
        {
          name = "GatewaySubnet"
          properties = {
            addressPrefix = "192.168.1.0/24"
          }
        }
      ]
    }
  })
}

# resource "azapi_resource" "pip" {
#   type      = "Microsoft.Network/publicIPAddresses@2021-08-01"
#   name      = "${var.virtual_network_name}-pip"
#   parent_id = azapi_resource.rg.id
#   location  = var.virtual_network_location
#   body = jsonencode({
#     properties = {
#       publicIPAllocationMethod = "Static"
#       deleteOption             = "Delete"
#     }
#     sku = {
#       name = "Standard"
#       tier = "Regional"
#     }
#   })
# }

# resource "azapi_resource" "gw" {
#   type      = "Microsoft.Network/virtualNetworkGateways@2021-08-01"
#   name      = "${var.virtual_network_name}-gw"
#   parent_id = azapi_resource.rg.id
#   location  = var.virtual_network_location
#   body = jsonencode({
#     properties = {
#       activeActive = false
#       enableBgp    = false
#       gatewayType  = "Vpn"
#       vpnType      = "RouteBased"
#       vpnGatewayGeneration = "Generation2"
#       sku = {
#         name = "VpnGw2"
#         tier = "VpnGw2"
#       }
#       ipConfigurations = [
#         {
#           name = "ipconfig"
#           properties = {
#             privateIPAllocationMethod = "Dynamic"
#             publicIPAddress = {
#               id = azapi_resource.pip.id
#             }
#             subnet = {
#               id = "${azapi_resource.hub.id}/subnets/GatewaySubnet"
#             }
#           }
#         }
#       ]
#     }
#   })
# }

module "virtualnetwork_test" {
  source                              = "../../"
  subscription_id                     = var.subscription_id
  virtual_network_address_space       = var.virtual_network_address_space
  virtual_network_location            = var.virtual_network_location
  virtual_network_resource_group_name = var.virtual_network_resource_group_name
  virtual_network_name                = var.virtual_network_name
  virtual_network_enable_peering      = true
  hub_network_resource_id             = azapi_resource.hub.id
  virtual_network_use_remote_gateways = var.virtual_network_use_remote_gateways
}
