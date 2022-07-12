# ALZ landing zone role assignment submodule

## Overview

Creates a role assignment at subscription or lower scope.

## Notes

See [README.md](../../README.md) in the parent module for more information.

## Example

```terraform
module "roleassignment" {
  source          = "Azure/lz-vending/azurerm/modules/roleassignment"
  version         = "~> 0.1.0"
  role_definition = "Owner"
  scope           = "/subscriptions/00000000-0000-0000-0000-000000000000"
  principal_id    = "00000000-0000-0000-0000-000000000000"
}
```
