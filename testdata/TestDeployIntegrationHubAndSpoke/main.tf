terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.7.0"
    }
    azapi = {
      source  = "Azure/azapi"
      version = ">= 0.3.0"
    }
  }
}

resource "azurerm_resource_group" "hub" {
  name     = "${var.virtual_network_name}-hub"
  location = var.location
}

resource "azurerm_virtual_network" "hub" {
  name                = "${var.virtual_network_name}-hub"
  location            = azurerm_resource_group.hub.location
  resource_group_name = azurerm_resource_group.hub.name
  address_space       = ["192.168.0.0/23"]
}

module "alz_landing_zone" {
  source = "../../"

  location = var.location

  # subscription variables
  subscription_alias_enabled = var.subscription_alias_enabled
  subscription_billing_scope = var.subscription_billing_scope
  subscription_display_name  = var.subscription_display_name
  subscription_alias_name    = var.subscription_alias_name
  subscription_workload      = var.subscription_workload

  # virtual network variables
  virtual_network_enabled             = var.virtual_network_enabled
  virtual_network_address_space       = var.virtual_network_address_space
  virtual_network_name                = var.virtual_network_name
  virtual_network_resource_group_name = var.virtual_network_resource_group_name

  # virtual network peering
  virtual_network_peering_enabled     = var.virtual_network_peering_enabled
  virtual_network_use_remote_gateways = var.virtual_network_use_remote_gateways
  hub_network_resource_id             = azurerm_virtual_network.hub.id
}
