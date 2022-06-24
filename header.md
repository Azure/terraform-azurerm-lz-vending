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

None.

## Example

```terraform
module "alz_landing_zone" {
  # Terraform Cloud/Enterprise use
  source  = "Azure/alz-landing-zone/azurerm"
  version = "~>0.0.1"
  # TBC
}
```
