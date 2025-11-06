# Landing zone virtual network submodule

## Overview

Creates multiple virtual networks in the supplied subscription.
Optionally:

- Creates bi-directional peering and/or a virtual WAN connection
- Creates peerings between the virtual networks (mesh peering)

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "virtualnetwork" {
  source  = "Azure/lz-vending/azurerm/modules/virtualnetwork"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  virtual_networks = {
    vnet1 = {
      name                = "myvnet"
      address_space       = ["192.168.0.0/24", "10.0.0.0/24"]
      location            = "westeurope"
      resource_group_name = "myrg"
    },
    vnet2 = {
      name                = "myvnet2"
      address_space       = ["192.168.1.0/24", "10.0.1.0/24"]
      location            = "northeurope"
      resource_group_name = "myrg2"
    }
  }
}
```
