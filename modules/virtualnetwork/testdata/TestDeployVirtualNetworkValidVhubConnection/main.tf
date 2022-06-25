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

variable "virtual_network_enable_vwan_connection" {
  type = bool
}

resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  parent_id = "/subscriptions/${var.subscription_id}"
  name      = "${var.virtual_network_name}-hub"
  location  = var.virtual_network_location
}

resource "azapi_resource" "vwan" {
  type      = "Microsoft.Network/virtualWans@2021-08-01"
  name      = "${var.virtual_network_name}-vwan"
  location  = var.virtual_network_location
  parent_id = azapi_resource.rg.id
  body = jsonencode({
    properties = {
      type                       = "Standard"
      allowBranchToBranchTraffic = true
      allowVnetToVnetTraffic     = true
      disableVpnEncryption       = false
    }
  })
}

resource "azapi_resource" "vhub" {
  type     = "Microsoft.Network/virtualHubs@2021-08-01"
  name     = "${var.virtual_network_name}-vhub"
  location = var.virtual_network_location
  parent_id = azapi_resource.rg.id
  body = jsonencode({
    properties = {
      addressPrefix = "192.168.100.0/23"
      sku           = "Standard"
      virtualWan = {
        id = azapi_resource.vwan.id
      }
    }
  })
}

module "virtualnetwork_test" {
  source                                 = "../../"
  subscription_id                        = var.subscription_id
  virtual_network_address_space          = var.virtual_network_address_space
  virtual_network_location               = var.virtual_network_location
  virtual_network_resource_group_name    = var.virtual_network_resource_group_name
  virtual_network_name                   = var.virtual_network_name
  virtual_network_enable_vwan_connection = var.virtual_network_enable_vwan_connection
  vwan_hub_resource_id                   = azapi_resource.vhub.id
}
