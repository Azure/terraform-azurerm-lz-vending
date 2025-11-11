# Integration Tests for LZ Vending Module
# These tests validate complete scenarios with multiple components

run "integration_hub_and_spoke" {
  command = plan

  variables {
    location                        = "northeurope"
    subscription_alias_enabled      = true
    subscription_display_name       = "test-subscription-alias"
    subscription_alias_name         = "test-subscription-alias"
    subscription_workload           = "Production"
    subscription_billing_scope      = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
    subscription_tags = {
      test-tag   = "test-value"
      test-tag-2 = "test-value-2"
    }
    resource_group_creation_enabled = true
    resource_groups = {
      primary = {
        name     = "primary-rg"
        location = "westeurope"
      }
      secondary = {
        name     = "secondary-rg"
        location = "westeurope"
      }
    }
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                      = "primary-vnet"
        address_space             = ["192.168.0.0/24"]
        location                  = "westeurope"
        resource_group_key        = "primary"
        hub_peering_enabled       = true
        hub_network_resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
      }
    }
  }

  assert {
    condition     = module.subscription[0].subscription_resource_id != null
    error_message = "Subscription should be created"
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) >= 1
    error_message = "At least one virtual network should be planned"
  }
}

run "integration_vwan" {
  command = plan

  variables {
    location                        = "northeurope"
    subscription_alias_enabled      = true
    subscription_display_name       = "test-subscription-alias"
    subscription_alias_name         = "test-subscription-alias"
    subscription_workload           = "Production"
    subscription_billing_scope      = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
    subscription_tags = {
      test-tag = "test-value"
    }
    resource_group_creation_enabled = true
    resource_groups = {
      primary = {
        name     = "primary-rg"
        location = "westeurope"
      }
      secondary = {
        name     = "secondary-rg"
        location = "westeurope"
      }
    }
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                   = "primary-vnet"
        address_space          = ["192.168.0.0/24"]
        location               = "westeurope"
        resource_group_key     = "primary"
        vwan_hub_resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualHubs/testhub"
        vwan_connection_enabled = true
      }
    }
  }

  assert {
    condition     = module.subscription[0].subscription_resource_id != null
    error_message = "Subscription should be created"
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) >= 1
    error_message = "At least one virtual network should be planned"
  }
}

run "integration_subscription_and_roleassignment_only" {
  command = plan

  variables {
    location                        = "northeurope"
    subscription_alias_enabled      = true
    subscription_display_name       = "test-subscription-alias"
    subscription_alias_name         = "test-subscription-alias"
    subscription_workload           = "Production"
    subscription_billing_scope      = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
    virtual_network_enabled         = false
    role_assignment_enabled         = true
    role_assignments = {
      ra = {
        principal_id   = "00000000-0000-0000-0000-000000000000"
        definition     = "Owner"
        relative_scope = ""
      }
    }
    resource_group_creation_enabled = false
  }

  assert {
    condition     = module.subscription[0].subscription_resource_id != null
    error_message = "Subscription should be created"
  }
}

run "integration_existing_subscription_hub_and_spoke" {
  command = plan

  variables {
    location                        = "northeurope"
    subscription_alias_enabled      = false
    subscription_id                 = "00000000-0000-0000-0000-000000000000"
    resource_group_creation_enabled = true
    resource_groups = {
      primary = {
        name     = "primary-rg"
        location = "westeurope"
      }
      secondary = {
        name     = "secondary-rg"
        location = "westeurope"
      }
    }
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name                      = "primary-vnet"
        address_space             = ["192.168.0.0/24"]
        location                  = "westeurope"
        resource_group_key        = "primary"
        hub_network_resource_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
        hub_peering_enabled       = true
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) >= 1
    error_message = "At least one virtual network should be planned"
  }
}

run "integration_resource_groups_only" {
  command = plan

  variables {
    subscription_id                 = "00000000-0000-0000-0000-000000000000"
    location                        = "westeurope"
    resource_group_creation_enabled = true
    resource_groups = {
      NetworkWatcherRG = {
        location = "westeurope"
        name     = "NetworkWatcherRG"
      }
      rg1 = {
        location = "westeurope"
        name     = "rg1"
      }
    }
  }

  assert {
    condition     = length(keys(var.resource_groups)) == 2
    error_message = "Expected 2 resource groups to be defined"
  }
}

run "integration_vnet_with_route_table" {
  command = plan

  variables {
    location                        = "northeurope"
    subscription_alias_enabled      = true
    subscription_display_name       = "test-subscription-alias"
    subscription_alias_name         = "test-subscription-alias"
    subscription_workload           = "Production"
    subscription_billing_scope      = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
    subscription_tags = {
      test-tag   = "test-value"
      test-tag-2 = "test-value-2"
    }
    resource_group_creation_enabled = true
    resource_groups = {
      primary = {
        name     = "primary-rg"
        location = "westeurope"
      }
      secondary = {
        name     = "secondary-rg"
        location = "westeurope"
      }
    }
    virtual_network_enabled = true
    virtual_networks = {
      primary = {
        name               = "primary-vnet"
        address_space      = ["192.168.0.0/24"]
        location           = "westeurope"
        resource_group_key = "primary"
        subnets = {
          primary = {
            name             = "primary-subnet"
            address_prefixes = ["192.168.0.0/25"]
            route_table = {
              key_reference = "primary"
            }
          }
          secondary = {
            name             = "secondary-subnet"
            address_prefixes = ["192.168.0.128/25"]
            route_table = {
              id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/primary-rg/providers/Microsoft.Network/routeTables/primary-route-table"
            }
          }
        }
      }
    }
    route_table_enabled = true
    route_tables = {
      primary = {
        name               = "primary-route-table"
        resource_group_key = "primary"
        location           = "westeurope"
      }
      default = {
        name                         = "default-route-table"
        resource_group_name_existing = "primary-rg"
        location                     = "westeurope"
      }
    }
  }

  assert {
    condition     = module.subscription[0].subscription_resource_id != null
    error_message = "Subscription should be created"
  }

  assert {
    condition     = length(keys(module.virtualnetwork[0].virtual_network_resource_ids)) >= 1
    error_message = "At least one virtual network should be planned"
  }
}
