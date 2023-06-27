# ALZ landing zone resource group submodule

## Overview

Creates the supplied resource group in the specified location and subscription.

Useful in a subscription vending scenario to pre-create the `NetworkWatcherRG` so we may delete it prior to subscription cancellation. Also useful for pre-creating resource groups for the purposes of RBAC delegation.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "resourcegroups" {
  source  = "Azure/lz-vending/azurerm/modules/resourcegroups"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  location        = "eastus"
  name            = "rg-test"
}
```
