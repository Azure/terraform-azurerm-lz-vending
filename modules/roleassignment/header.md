# Landing zone role assignment submodule

## Overview

Creates a role assignment at subscription or lower scope.
Module is designed to be instantiated many times, once per role assignment.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "roleassignment" {
  source          = "Azure/lz-vending/azurerm/modules/roleassignment"
  version         = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints
  role_definition = "Owner"
  scope           = "/subscriptions/00000000-0000-0000-0000-000000000000"
  principal_id    = "00000000-0000-0000-0000-000000000000"
}
```
