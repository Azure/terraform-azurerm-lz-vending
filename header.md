# ALZ landing zone module

## Overview

The landing zone Terraform module is designed to accelerate deployment of the individual landing zones into the ALZ conceptual architecture.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Hub & spoke networking
- Virtual WAN networking
- More to come!

## Notes

TBC.

## Example

```terraform
module "alz_landing_zone" {
  source  = "Azure/alz-landing-zone/azurerm"
  version = "~>0.1.0"

  # subscription variables
  subscription_alias_enabled       = true
  subscription_alias_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_alias_display_name  = "my-subscription-display-name"
  subscription_alias_name          = "my-subscription-alias"
  subscription_alias_workload      = "Production"

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24]
  virtual_network_location            = "eastus"
  virtual_network_resource_group_name = "my-network-rg"

  # virtual network peering
  virtual_network_peering_enabled = true
  hub_network_resource_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"
}
```
