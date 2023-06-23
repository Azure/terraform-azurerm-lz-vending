# ALZ landing zone virtual network submodule

## Overview

Creates a NetworkWatcherRG resource group in the specified location and subscription.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "networkwatcherrg" {
  source  = "Azure/lz-vending/azurerm/modules/networkwatcherrg"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  location        = "eastus"
}
```
