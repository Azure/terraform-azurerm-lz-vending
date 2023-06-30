# Landing zone resource provider submodule

## Overview

Registers resource providers and features. Must be performed after all other provisioning is complete.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "resourceproviders" {
  source  = "Azure/lz-vending/azurerm/modules/resourceprovider"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  resource_provider = "Microsoft.PowerBI"
  features          = ["DailyPrivateLinkServicesForPowerBI"]
}
```
