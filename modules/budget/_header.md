# Landing zone budget submodule

## Overview

Creates a budget in Azure. Designed to be instantiated multiple times to create multiple budgets.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "budget" {
  source  = "Azure/lz-vending/azurerm/modules/roleassignment"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  budget_name       = "budget1"
  budget_amount     = 100
  budget_scope      = "/subscriptions/00000000-0000-0000-0000-000000000000"
  budget_time_grain = "Monthly"
  budget_time_period = {
    start_date = "2024-01-01"
    end_date   = "2025-01-01"
  }
}
```
