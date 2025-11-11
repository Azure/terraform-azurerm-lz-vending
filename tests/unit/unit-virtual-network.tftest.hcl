# Virtual Network Module Basic Tests
# Tests basic VNet creation without deploying

mock_provider "azurerm" {}
mock_provider "azapi" {}
mock_provider "modtem" {}
mock_provider "time" {}

variables {
  location        = "uksouth"
  subscription_id = "00000000-0000-0000-0000-000000000000"
  enable_telemetry = false
}

run "valid_two_vnets" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                = "primary-vnet"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name_existing = "primary-rg"
      }
      secondary = {
        name                = "secondary-vnet"
        address_space       = ["192.168.1.0/24"]
        location            = "northeurope"
        resource_group_name_existing = "secondary-rg"
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) == 2
    error_message = "Expected 2 virtual networks to be created"
  }
}
