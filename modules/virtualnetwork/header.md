# ALZ landing zone virtual network submodule

## Overview

Creates a virtual network in the supplied subscription.
Optionally, created bi-directional peering and/or a virtual WAN connection.

## Notes

See [README.md](../../README.md) in the parent module for more information.

## Example

```terraform
module "virtualnetwork" {
  source  = "Azure/lz-vending/azurerm/modules/virtualnetwork"
  version = "~> 0.1.0"

  subscription_id                     = "00000000-0000-0000-0000-000000000000"
  virtual_network_name                = "my-virtual-network"
  virtual_network_resource_group_name = "my-network-rg"
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_location            = "eastus"

  virtual_network_peering_enabled = true
  hub_network_resource_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"
}
```
