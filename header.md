# Terraform landing zone vending module for Azure

[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/Azure/terraform-azurerm-lz-vending.svg)](http://isitmaintained.com/project/Azure/terraform-azurerm-lz-vending "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/Azure/terraform-azurerm-lz-vending.svg)](http://isitmaintained.com/project/Azure/terraform-azurerm-lz-vending "Percentage of issues still open")
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/Azure/terraform-azurerm-lz-vending/badge)](https://scorecard.dev/viewer/?uri=github.com/Azure/terraform-azurerm-lz-vending)

## Overview

The landing zone Terraform module is designed to accelerate deployment of individual landing zones within an Azure tenant.
We use the [AzureRM](https://registry.terraform.io/providers/hashicorp/azurerm/latest) and [AzAPI](https://registry.terraform.io/providers/azure/azapi/latest) providers to create the subscription and deploy the resources in a single `terraform apply` step.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
  - Microsoft Defender for Cloud (DFC) security contact
- Networking - deploy multiple vnets with:
  - Hub & spoke connectivity (peering to a hub network)
  - vWAN connectivity
  - Mesh peering (peering between spokes)
- Role assignments
- Resource provider (and feature) registration
- Resource group creation
- User assigned managed identity creation
  - Federated credential configuration for GitHub Actions, Terraform Cloud, and other providers.

> When creating virtual network peerings, be aware of the [limit of peerings per virtual network](https://learn.microsoft.com/azure/azure-resource-manager/management/azure-subscription-service-limits?toc=%2Fazure%2Fvirtual-network%2Ftoc.json#azure-resource-manager-virtual-networking-limits).

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Change log

Please see the [GitHub releases pages](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for change log information.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

The below example created a landing zone subscription with two virtual networks.
One virtual network is in the default location of the subscription, the other is in a different location.

The virtual networks are peered with the supplied hub network resource ids, they are also peered with each other using the mesh peering option.

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # Set the default location for resources
  location = "westeurope"

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  network_watcher_resource_group_enabled = true

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "Corp"

  # defender for cloud variables
  subscription_dfc_contact_enabled = true
  subscription_dfc_contact = {
    emails = "john@microsoft.com;jane@microsoft.com"
  }

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks = {
    one = {
      name                    = "my-vnet"
      address_space           = ["192.168.1.0/24"]
      hub_peering_enabled     = true
      hub_network_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"
      mesh_peering_enabled    = true
    }
    two = {
      name                    = "my-vnet2"
      location                = "northeurope"
      address_space           = ["192.168.2.0/24"]
      hub_peering_enabled     = true
      hub_network_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network2"
      mesh_peering_enabled    = true
    }
  }

  umi_enabled             = true
  umi_name                = "umi"
  umi_resource_group_name = "rg-identity"
  umi_role_assignments = {
    myrg-contrib = {
      definition     = "Contributor"
      relative_scope = "/resourceGroups/MyRg"
    }
  }

  resource_group_creation_enabled = true
  resource_groups = {
    myrg = {
      name     = "MyRg"
      location = "westeurope"
    }
  }

  # role assignments
  role_assignment_enabled = true
  role_assignments = {
    # using role definition name, created at subscription scope
    contrib_user_sub = {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Contributor"
      relative_scope = ""
    },
    # using a custom role definition
    custdef_sub_scope = {
      principal_id   = "11111111-1111-1111-1111-111111111111"
      definition     = "/providers/Microsoft.Management/MyMg/providers/Microsoft.Authorization/roleDefinitions/ffffffff-ffff-ffff-ffff-ffffffffffff"
      relative_scope = ""
    },
    # using relative scope (to the created or supplied subscription)
    rg_owner = {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Owner"
      relative_scope = "/resourceGroups/MyRg"
    },
  }
}
```
