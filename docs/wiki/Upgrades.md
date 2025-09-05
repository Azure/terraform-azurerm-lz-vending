<!-- markdownlint-disable MD041 -->
This module uses [Semantic Versioning (SemVer)](https://semver.org/s) versioning.

Given a version number `MAJOR.MINOR.PATCH`, we increment the:

* `MAJOR` version when we make incompatible/breaking changes,
* `MINOR` version when we add functionality in a backwards compatible manner, and
* `PATCH` version when we make backwards compatible bug fixes.

## Upgrade process

If you are upgrading to a new `MINOR` or `PATCH` release, you will not see any breaking changes.
If you are using the Terraform registry, you can update the version number for the module.

For a new `MAJOR` release, you will see breaking changes.
We will publish guidance in the release notes on GitHub.

See the [release notes](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for more information.

## Upgrading from v1.x to v2.x

v2 of the module makes large-scale changes to the virtual networking capabilities of the module.
We therefore recommend that you keep any existing instances of the module at v1, and use v2 going forward for new instances.
If you would like multiple vnets in the same subscription using v1 of the module you can use the pattern [described here](https://github.com/Azure/terraform-azurerm-lz-vending/issues/97#issuecomment-1240712419)

## Upgrading from v2.x to v3.x

v3 of the module makes changes the the `role_assignments` variable, changing the format of the variable from a list of objects `list(object({...}))` to a map of objects `map(object({...}))`.
This change fixes [#153](https://github.com/Azure/terraform-azurerm-lz-vending/issues/153).

> **Due to the map key changing, all role assignments will be deleted and re-created.**

By way of explanation, in order to run the `for_each` loop on the role assignments, we need to use either a set or a map.
If using a map, Terraform needs to know all keys at plan time, they cannot be 'known after apply'.
Previously we converted the list of objects into a map of objects, using the `uuidv5()` function to generate predictable map keys from the inputs.
Unfortunately this caused issues when any of the inputs to the `uuidv5()` function were not known at plan time, in this case it was the principal id.

Rather than revert to a set, where ordering can be an issue, we decided to change the input variable to be a map from the outset.
This does mean a small change is required, you must specify a map key. This can be anything but do not use a reference to other object.

### v2.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = [
    {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  ]
}
```

### v3.x `role_assignments` syntax

```terraform
module "lz_vending" {
  source  = "..."
  version = "..."

  # (other input variables hidden)

  role_assignments = {
    contrib_to_group = {
      principal_id   = "..."
      definition     = "contributor"
      relative_scope = ""
    }
  }
}
```

## Upgrading from v4.x to v5.x

###  Terraform version

We now require a minimum of Terraform version 1.8.

### Provider Versions

We now require a minimum of AzureRM version 4.0 and AzAPI version 2.2.

### Resource Groups

We have removed the boolean input variable to create the network watcher resource group.
Instead, use `var.resource_groups` to specify the resource groups to create.

We have used the `moved {}` block to move the resource in state.
If you previously deployed the network watcher resource group, please modify the value of `var.resource_groups` to include the existing resource group. The key ***must*** be `NetworkWatcherRG` You can use the following example:

```hcl
resource_groups = {
  NetworkWatcherRG = {
    name     = "NetworkWatcherRG"
    location = "your-location"
    tags     = {} # add tags here
  }
}
```

## Virtual WAN

When joining virtual networks to a Virtual WAN hub, the behaviour with routing intent has changed.
Previously the AzAPI provider allowed us to use `ignore_body_properties` to dynamically ignore parts of the resource body
With AzAPI v2 this is no longer possible, so we have to use the `lifecycle` block to ignore changes.
However, as ignore changes is not able to be user configurable, we have had to split the virtual hub connections into two separate resources.

In order to avoid destroying and re-creating the virtual hub connections, you will have to use the `moved {}` block to move the resource in state.
We are unable to do this for you because we do not know the specific instances of the resources that require moving.

```hcl
moved {
  from = module.YOUR_MODULE_ALIAS.module.virtualnetwork.azapi_resource.vhubconnection["instance_name"]
  to   = module.YOUR_MODULE_ALIAS.module.virtualnetwork.azapi_resource.vhubconnection_routing_intent["instance_name"]
}
```

## Upgrading from v5.x to v6.x

### Resource Group Module

We have removed the resource group and lock provisioning capability from the virtual network and user assigned identity submodules. Add new resource group objects to the `resource_groups` map and set the root module variable `resource_group_creation_enabled` to `true`.

```hcl
resource_group_creation_enabled = true
resource_groups = {
  vnetrg = {
    name         = local.network_rg
    location     = var.location
    lock_enabled = true
    lock_name    = "lock-network-${local.component_name}-01"
  }
  mainrg = {
    name     = local.application_rg
    location = var.location
  }
  identityrg = {
    name         = local.identity_rg
    location     = var.location
    lock_enabled = true
    lock_name    = "lock-umi-${local.component_name}-plan-01"
  }
}
```

Add Terraform moved blocks for your resource group and resource group locks as shown below.

```hcl
# VNET
moved {
  from = module.lz_vending.module.virtualnetwork[0].azapi_resource.rg["<resource-group-name-value>"]
  to   = module.lz_vending.module.resourcegroup["<resource-groups-map-key-name>"].azapi_resource.rg
}
# VNET LOCK
moved {
  from = module.lz_vending.module.virtualnetwork[0].azapi_resource.rg_lock["<resource-group-lock-name-value>"]
  to   = module.lz_vending.module.resourcegroup["<resource-groups-map-key-name>"].azapi_resource.rg_lock[0]
}
# UMI
moved {
  from = module.lz_vending.module.usermanagedidentity["<user-managed-identity-map-key-name>"].azapi_resource.rg[0]
  to   = module.lz_vending.module.resourcegroup["<resource-groups-map-key-name>"].azapi_resource.rg
}
# UMI LOCK
moved {
  from = module.lz_vending.module.usermanagedidentity["<user-managed-identity-map-key-name>"].azapi_resource.rg_lock[0]
  to   = module.lz_vending.module.resourcegroup["<resource-groups-map-key-name>"].azapi_resource.rg_lock[0]
}
```

### Virtual Network Module

For virtual networks change the `resource_group_name` attribute to `resource_group_key` and change the value to the key name that corresponds to the object in the `resource_groups` map. This is to maintain consistency throughout the submodules.

**NOTE:** If you do not wish to manage the resource group creation with this module at all, then you may specify the `resource_group_name_existing` attribute instead to leverage an already existing resource group within the subscription. This module first checks for the `resource_group_key` and then the `resource_group_name_existing` as a fallback.

```hcl
virtual_network_enabled = true
virtual_networks = {
  primary = {
    name                    = "vnet-${local.component_name}-01"
    address_space           = [var.spoke_vnet_address_space]
    resource_group_key      = "vnetrg"
    hub_peering_enabled     = true
    hub_network_resource_id = var.hub_network_id
    hub_peering_direction   = "both"
    dns_servers             = []
    subnets = {
      subnet1 = {
        name = "snet-${local.component_name}-01"
        #address_prefix                               = module.ip_calc.address_prefixes["default"]
        address_prefixes                              = [module.ip_calc.address_prefixes["default"]]
        private_endpoint_network_policies             = "Disabled"
        private_link_service_network_policies_enabled = false
        route_table = {
          key_reference = "HubNetwork"
        }
        network_security_group = {
          key_reference = "default"
        }
        service_endpoints               = []
        default_outbound_access_enabled = false
      }
      subnet2 = {
        name = "snet-${local.component_name}-pe-02"
        #address_prefix                               = module.ip_calc.address_prefixes["private_endpoint"]
        address_prefixes                              = [module.ip_calc.address_prefixes["private_endpoint"]]
        private_endpoint_network_policies             = "Disabled"
        private_link_service_network_policies_enabled = false
        default_outbound_access_enabled               = false
        service_endpoints                             = []
        delegations                                   = []
      }
    }
    hub_peering_options_tohub = {
      allow_forwarded_traffic      = true
      allow_gateway_transit        = false
      allow_virtual_network_access = true
      peer_complete_vnets          = true
      use_remote_gateways          = false
    }
    hub_peering_options_fromhub = {
      allow_forwarded_traffic      = true
      allow_gateway_transit        = false
      allow_virtual_network_access = true
      peer_complete_vnets          = true
      use_remote_gateways          = false
    }
  }
}
```

### Route Table Module

For route tables change the `resource_group_name` attribute to `resource_group_key` and change the value to the key name that corresponds to the object in the `resource_groups` map. This is to maintain consistency throughout the submodules.

**NOTE:** If you do not wish to manage the resource group creation with this module at all, then you may specify the `resource_group_name_existing` attribute instead to leverage an already existing resource group within the subscription. This module first checks for the `resource_group_key` and then the `resource_group_name_existing` as a fallback.

```hcl
route_table_enabled = true
route_tables = {
  HubNetwork = {
    name                          = "rt-${local.component_name}-01"
    location                      = var.location
    resource_group_key            = "vnetrg"
    bgp_route_propagation_enabled = false
    routes = {
      FirewallDefaultRoute = {
        name                   = "${var.application_short_name}-to-firewall"
        address_prefix         = "0.0.0.0/0"
        next_hop_type          = "VirtualAppliance"
        next_hop_in_ip_address = var.hub_fw_ip
      }
    }
  }
}
```

### Network Security Group Module

For network security groups, change the `resource_group_name` attribute to `resource_group_key` and change the value to the key name that corresponds to the object in the `resource_groups` map. This is to maintain consistency throughout the submodules.

**NOTE:** If you do not wish to manage the resource group creation with this module at all, then you may specify the `resource_group_name_existing` attribute instead to leverage an already existing resource group within the subscription. This module first checks for the `resource_group_key` and then the `resource_group_name_existing` as a fallback.

```hcl
network_security_group_enabled = true
network_security_groups = {
  default = {
    name               = "nsg-${local.component_name}-01"
    location           = var.location
    resource_group_key = "vnetrg"
    security_rules = {
      allow_outbound = {
        name                         = "allow-spoke-outbound"
        priority                     = 100
        direction                    = "Outbound"
        access                       = "Allow"
        protocol                     = "Tcp"
        source_port_ranges           = ["80", "443"]
        destination_port_ranges      = ["80", "443"]
        source_address_prefixes      = [var.spoke_vnet_address_space]
        destination_address_prefixes = [var.hub_network_address_prefix]
        description                  = "Allow spoke outbound traffic to FW"
      }
      allow_inbound = {
        name                         = "allow-spoke-inbound"
        priority                     = 100
        direction                    = "Inbound"
        access                       = "Allow"
        protocol                     = "Tcp"
        source_port_ranges           = ["80", "443"]
        destination_port_ranges      = ["80", "443"]
        source_address_prefixes      = [var.hub_network_address_prefix]
        destination_address_prefixes = [var.spoke_vnet_address_space]
        description                  = "Allow spoke inbound traffic to FW"
      }
    }
  }
}
```

### User Managed Identity Module

For user managed identities change the `resource_group_name` attribute to `resource_group_key` and change the value to the key name that corresponds to the object in the `resource_groups` map. This is to maintain consistency throughout the submodules.

**NOTE:** If you do not wish to manage the resource group creation with this module at all, then you may specify the `resource_group_name_existing` attribute instead to leverage an already existing resource group within the subscription. This module first checks for the `resource_group_key` and then the `resource_group_name_existing` as a fallback.

```hcl
umi_enabled = true
user_managed_identities = {
  plan = {
    name               = "umi-${local.component_name}-plan-01"
    location           = var.location
    resource_group_key = "identityrg"
    tags               = local.tags
    role_assignments = {
      reader = {
        definition      = "Reader"
        relative_scope  = ""
        use_random_uuid = true
      }
    }
  },
  apply = {
    name               = "umi-${local.component_name}-apply-01"
    location           = var.location
    resource_group_key = "identityrg"
    tags               = local.tags
    role_assignments = {
      apply = {
        definition      = "Contributor"
        relative_scope  = ""
        use_random_uuid = true
      }
    }
  },
  app = {
    name               = "umi-${local.component_name}-app-01"
    location           = var.location
    resource_group_key = "mainrg"
    tags               = local.tags
  }
}
```
