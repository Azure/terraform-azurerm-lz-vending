# Landing zone subscription submodule

## Overview

Creates a subscription alias
Optionally:

- Associates the resulting subscription to a management group 
- Creates the Microsoft Defender for Cloud (DFC) security contact and enables notifications

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "subscription" {
  source  = "Azure/lz-vending/azurerm/modules/subscription"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_alias_billing_scope       = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_alias_display_name        = "my-subscription-display-name"
  subscription_alias_name                = "my-subscription-alias"
  subscription_alias_workload            = "Production"
  subscription_alias_management_group_id = "mymg"
  subscription_dfc_contact_enabled       = true
  subscription_dfc_contact = {
    notifications_by_role = ["Owner", "Contributor"]
    emails                = "john@microsoft.com;jane@microsoft.com"
    phone                 = "+1-555-555-5555"
    alert_notifications   = "Medium"
  }
}
```
