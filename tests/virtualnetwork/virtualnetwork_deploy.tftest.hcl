# Virtual Network Deployment Tests
# These tests actually deploy resources to Azure
# Only run these when TERRATEST_DEPLOY environment variable is set

# These tests require:
# - AZURE_SUBSCRIPTION_ID to be set
# - AZURE_TENANT_ID to be set
# - Azure authentication configured
# - Resource groups to be created beforehand

run "deploy_basic_vnets" {
  command = apply

  variables {
    subscription_id  = "00000000-0000-0000-0000-000000000000"  # Set via env var at runtime
    enable_telemetry = false
    virtual_networks = {
      primary = {
        name                = "test-vnet-primary"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name = "test-rg-primary"
      }
      secondary = {
        name                = "test-vnet-secondary"
        address_space       = ["192.168.1.0/24"]
        location            = "northeurope"
        resource_group_name = "test-rg-secondary"
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork.virtual_network_resource_ids)) == 2
    error_message = "Should have created 2 virtual networks"
  }

  assert {
    condition     = module.virtualnetwork.virtual_network_resource_ids["primary"] != ""
    error_message = "Primary VNet resource ID should not be empty"
  }

  assert {
    condition     = module.virtualnetwork.virtual_network_resource_ids["secondary"] != ""
    error_message = "Secondary VNet resource ID should not be empty"
  }
}

run "deploy_vnets_with_subnets" {
  command = apply

  variables {
    subscription_id  = "00000000-0000-0000-0000-000000000000"
    enable_telemetry = false
    virtual_networks = {
      primary = {
        name                = "test-vnet-primary"
        address_space       = ["192.168.0.0/24"]
        location            = "westeurope"
        resource_group_name = "test-rg-primary"
        subnets = {
          default = {
            name                                          = "snet-default"
            address_prefixes                              = ["192.168.0.0/26"]
            private_link_service_network_policies_enabled = false
            private_endpoint_network_policies             = "Disabled"
          }
        }
      }
      secondary = {
        name                = "test-vnet-secondary"
        address_space       = ["192.168.1.0/24"]
        location            = "northeurope"
        resource_group_name = "test-rg-secondary"
        subnets = {
          default = {
            name                            = "snet-default"
            address_prefixes                = ["192.168.1.0/26"]
            default_outbound_access_enabled = true
            service_endpoints               = ["Microsoft.Storage"]
          }
          containers = {
            name             = "snet-containers"
            address_prefixes = ["192.168.1.64/26"]
            delegations = [
              {
                name = "Microsoft.ContainerInstance/containerGroups"
                service_delegation = {
                  name = "Microsoft.ContainerInstance/containerGroups"
                }
              }
            ]
          }
        }
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork.virtual_network_resource_ids)) == 2
    error_message = "Should have created 2 virtual networks with subnets"
  }
}

run "deploy_vnets_with_mesh_peering" {
  command = apply

  variables {
    subscription_id  = "00000000-0000-0000-0000-000000000000"
    enable_telemetry = false
    virtual_networks = {
      primary = {
        name                  = "test-vnet-primary"
        address_space         = ["192.168.0.0/24"]
        location              = "westeurope"
        resource_group_name   = "test-rg-primary"
        mesh_peering_enabled  = true
      }
      secondary = {
        name                 = "test-vnet-secondary"
        address_space        = ["192.168.1.0/24"]
        location             = "northeurope"
        resource_group_name  = "test-rg-secondary"
        mesh_peering_enabled = true
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork.virtual_network_resource_ids)) == 2
    error_message = "Should have created 2 virtual networks with mesh peering"
  }
}
