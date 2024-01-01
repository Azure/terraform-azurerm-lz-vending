<!-- markdownlint-disable MD041 -->
## Summary

In this example we will highlight a use case for the `subscription_use_azapi` variable.
| :warning: WARNING!          |
|:---------------------------|
| This is not a common use case and we recommend keeping this defaulted to false, unless you are met with a scenario similar to that which is highlighted below.|

## Scenario

When vending subscriptions we must pay attention to the **default management group**. Without any intervention, the default management group will be the root management group, but this can be changed in the portal.
See [Setting - Default management group](https://learn.microsoft.com/en-us/azure/governance/management-groups/how-to/protect-resource-hierarchy#setting---default-management-group) for more information.

Consider the following scenario:

- An organisation has explicitly set the default management group to a management group that is **not** the root management group.
- The principal vending subscriptions has the necessary permissions on the `contoso` management group.
- The principal vending subscriptions has **no** permissions on the default management group.

In this scenario, the vending process will fail because, with the `azurerm` provider, the subscription is firstly created in the default management group, and then moved to the target management group. But the principal as no access to the default management group which is a necessary pre-requisite for moving the subscription.
See [Moving management groups and subscriptions](https://learn.microsoft.com/en-us/azure/governance/management-groups/overview#moving-management-groups-and-subscriptions) for more information.

To work around this issue, we have an additional variable `subscription_use_azapi` which when set to `true` will use the `azapi` provider to create the subscription and furthermore will be able to create the subscription in the target management group directly, bypassing the default management group.

### Example

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  location = "northeurope"

  # subscription variables
  subscription_alias_enabled = true
  subscription_alias_name    = "mylz"
  subscription_display_name  = "mylz"
  subscription_use_azapi     = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_workload      = "DevTest"

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "contoso"
}
```

## Drawbacks

This solution is non-standard, and necessarily comes with some caveats:

- Additional resource for updating the display name of the subscription is required.
- Additional resource for updating the subscription tags is required.
- The use of a `PUT` in a `azapi_resource_action` resource is used to continually update the subscription management group association, avoiding a `DELETE` and another `PUT`, which would result in an initial subscription move back to the default management group.
- Artificially managing the lifecycle of the subscription management group association using a data source, which is used to recreate the association if it has been moved outside of Terraform.
