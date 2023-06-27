# ALZ landing zone network watcher resource group submodule

## Overview

Creates the supplied resource groups in the specified location and subscription.

Useful in a subscription vending scenario to pre-create the `NetworkWatcherRG` so we may delete it prior to subscription cancellation.

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
