# Virtual Network Module Basic Tests
# Tests basic VNet creation without deploying

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

run "vnets_with_custom_dns" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                = "primary-vnet"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name_existing = "primary-rg"
        dns_servers         = ["1.2.3.4", "4.3.2.1"]
      }
      secondary = {
        name                = "secondary-vnet"
        address_space       = ["192.168.1.0/24"]
        location            = "northeurope"
        resource_group_name_existing = "secondary-rg"
        dns_servers         = ["8.8.8.8"]
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) == 2
    error_message = "Expected 2 virtual networks"
  }
}

run "vnets_with_tags" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                = "primary-vnet"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name_existing = "primary-rg"
        tags = {
          tag1 = "value1"
          tag2 = "2"
        }
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
    error_message = "Expected 2 virtual networks with tags"
  }
}

run "vnets_with_subnets" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                = "primary-vnet"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name_existing = "primary-rg"
        subnets = {
          default = {
            name              = "snet-default"
            address_prefixes  = ["192.168.0.0/26"]
          }
          privateendpoint = {
            name              = "snet-privateendpoint"
            address_prefixes  = ["192.168.0.64/26"]
          }
        }
      }
      secondary = {
        name                = "secondary-vnet"
        address_space       = ["192.168.1.0/24"]
        location            = "northeurope"
        resource_group_name_existing = "secondary-rg"
        subnets = {
          default = {
            name              = "snet-default"
            address_prefixes  = ["192.168.1.0/26"]
          }
        }
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) == 2
    error_message = "Expected 2 virtual networks with subnets"
  }
}

run "vnet_with_mesh_peering" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                  = "primary-vnet"
        address_space         = ["192.168.0.0/24"]
        location              = "westeurope"
        resource_group_name_existing   = "primary-rg"
        mesh_peering_enabled  = true
      }
      secondary = {
        name                                      = "secondary-vnet"
        address_space                             = ["192.168.1.0/24"]
        location                                  = "northeurope"
        resource_group_name_existing                       = "secondary-rg"
        mesh_peering_enabled                      = true
        mesh_peering_allow_forwarded_traffic      = true
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) == 2
    error_message = "Expected 2 virtual networks with mesh peering"
  }
}

run "vnet_with_hub_peering" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                      = "primary-vnet"
        address_space             = ["192.168.0.0/24"]
        location                  = "westeurope"
        resource_group_name_existing       = "primary-rg"
        hub_network_resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet2"
        hub_peering_enabled       = true
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
    error_message = "Expected 2 virtual networks, one with hub peering"
  }
}

run "vnet_with_ddos_protection" {
  command = plan

  variables {
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                      = "primary-vnet"
        address_space             = ["192.168.0.0/24"]
        location                  = "westeurope"
        resource_group_name_existing       = "primary-rg"
        ddos_protection_enabled   = true
        ddos_protection_plan_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test_rg/providers/Microsoft.Network/ddosProtectionPlans/test-ddos-plan"
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
    error_message = "Expected 2 virtual networks, one with DDoS protection"
  }
}
