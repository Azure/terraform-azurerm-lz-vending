# Terraform landing zone vending module for Azure

## Overview

The landing zone Terraform module is designed to accelerate deployment of the individual landing zones within an Azure tenant.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Hub & spoke networking
- Virtual WAN networking
- Role assignments

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Change log

Please see the [GitHub releases pages](https://github.com/Azure/terraform-azurerm-lz-vending/releases/latest) for change log information.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

```terraform
module "alz_landing_zone" {
  source  = "Azure/lz-vending/azurerm"
  version = "~>0.1.0"

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_location            = "eastus"
  virtual_network_resource_group_name = "my-network-rg"

  # virtual network peering
  virtual_network_peering_enabled = true
  hub_network_resource_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"

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
