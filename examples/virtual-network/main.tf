terraform {
  required_version = "~> 1.12"
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.5"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azapi_client_config" "current" {}

module "lz-vending" {
  source                          = "../../"
  subscription_id                 = data.azapi_client_config.current.subscription_id
  location                        = "swedencentral"
  resource_group_creation_enabled = true
  resource_groups = {
    rg1 = {
      name     = "rg-vending-001"
      location = "swedencentral"
    }
  }
  virtual_network_enabled = true
  virtual_networks = {
    vnet1 = {
      name               = "vnet-vending-001"
      address_space      = ["192.168.0.0/16"]
      location           = "swedencentral"
      resource_group_key = "rg1"
      subnets = {
        subnet1 = {
          name             = "snet-vending-001"
          address_prefixes = ["192.168.0.0/24"]
        }
      }
    }
  }
}
