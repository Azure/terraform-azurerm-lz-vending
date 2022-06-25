# ALZ landing zone subscription submodule

## Overview

Creates a subscription alias, optionally in the specified management group.

## Notes

See [README.md](../../README.md) in the parent module for more information.

## Example

```terraform
module "subscription" {
  source  = "Azure/alz-landing-zone/azurerm/modules/subscription"
  version = "~> 0.1.0"

  subscription_alias_billing_scope       = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_alias_display_name        = "my-subscription-display-name"
  subscription_alias_name                = "my-subscription-alias"
  subscription_alias_workload            = "Production"
  subscription_alias_management_group_id = "mymg"
}
```
