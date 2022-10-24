# Terraform landing zone vending module for Azure

## Overview

The landing zone Terraform module is designed to accelerate deployment of individual landing zones within an Azure tenant.
We use the [AzureRM][azurem_provider] and [AzAPI][azapi_provider] providers to create the subscription and deploy the resources in a single `terrafom apply` step.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Networking - deploy multiple vnets with:
  - Hub & spoke connectivity
  - vWAN connectivity
  - Mesh peering (peering between spokes)
- Role assignments

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Change log

Please see the [GitHub releases pages](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for change log information.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "Corp"

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks = {
    one = {
      name                    = "my-vnet"
      location                = "westeurope"
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

  # role assignments
  role_assignment_enabled = true
  role_assignments = [
    # using role definition name, created at subscription scope
    {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Contributor"
      relative_scope = ""
    },
    # using a custom role definition
    {
      principal_id   = "11111111-1111-1111-1111-111111111111"
      definition     = "/providers/Microsoft.Management/MyMg/providers/Microsoft.Authorization/roleDefinitions/ffffffff-ffff-ffff-ffff-ffffffffffff"
      relative_scope = ""
    },
    # using relative scope (to the created or supplied subscription)
    {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Owner"
      relative_scope = "/resourceGroups/MyRg"
    },
  ]
}
```

[azurem_provider]: https://registry.terraform.io/providers/hashicorp/azurerm/latest
[azapi_provider]: https://registry.terraform.io/providers/azure/azapi/latest
