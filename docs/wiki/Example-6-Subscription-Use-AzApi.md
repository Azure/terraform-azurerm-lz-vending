<!-- markdownlint-disable MD041 -->
## Summary

In this example we will highlight a use case for the `subscription_use_azapi` variable.
| :warning: WARNING!          |
|:---------------------------|
| This is not a common use case and we recommend keeping this defaulted to false, unless you are met with a scenario similar to that which is highlighted below.|

## Scenario

When vending subscriptions we must pay attention to the **default management group**. Without any intervention, the default management group will be the tenant root group, but this can be changed in the portal.
See [Setting - Default management group](https://learn.microsoft.com/en-us/azure/governance/management-groups/how-to/protect-resource-hierarchy#setting---default-management-group) for more information.

Consider the following scenario:

- An organisation has explicitly set the default management group to a management group that is **not** the tenant root group.
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

## Behaviour on Management Group Association Drift

In order to maintain the subscription management group association, we must use a data source to retrieve the current association and then use this to recreate the association if it has been moved outside of Terraform.

When drift is detected, the following will occur on the subsequent terraform plan:

```text
Terraform will perform the following actions:

  # module.lz_vending.module.subscription[0].azapi_resource_action.subscription_association[0] will be replaced due to changes in replace_triggered_by
-/+ resource "azapi_resource_action" "subscription_association" {
      ~ id          = "/providers/Microsoft.Management/managementGroups/contoso/subscriptions/00000000-0000-0000-0000-000000000000" -> (known after apply)
      ~ output      = jsonencode({}) -> (known after apply)
        # (4 unchanged attributes hidden)
    }

  # module.lz_vending.module.subscription[0].terraform_data.replacement[0] will be updated in-place
  ~ resource "terraform_data" "replacement" {
        id     = "xxxx"
      ~ input  = true -> false
      ~ output = true -> (known after apply)
    }

Plan: 1 to add, 1 to change, 1 to destroy.
```

Upon apply, this will place the subscription back into the target management group. However, this will result in a idempotency issue on the following plan which results in the following:

```text
Terraform will perform the following actions:

  # module.lz_vending.module.subscription[0].azapi_resource_action.subscription_association[0] will be replaced due to changes in replace_triggered_by
-/+ resource "azapi_resource_action" "subscription_association" {
      ~ id          = "/providers/Microsoft.Management/managementGroups/contoso/subscriptions/00000000-0000-0000-0000-000000000000" -> (known after apply)
      ~ output      = jsonencode({}) -> (known after apply)
        # (4 unchanged attributes hidden)
    }

  # module.lz_vending.module.subscription[0].terraform_data.replacement[0] will be updated in-place
  ~ resource "terraform_data" "replacement" {
        id     = "xxxx"
      ~ input  = false -> true
      ~ output = false -> (known after apply)
    }

Plan: 1 to add, 1 to change, 1 to destroy.
```

This is expected behavior since the `input` argument in `terraform_data.replacement` is monitoring the association between management group and subscription, and will move to false when the two are not aligned triggering a replacement.
The `input` argument will then immediately move back to true on the next plan (assuming no changes outside of terraform) triggering another replacement.

## Drawbacks

This solution is non-standard, and necessarily comes with some caveats:

- Additional resource for updating the display name of the subscription is required.
- Additional resource for updating the subscription tags is required.
- The use of a `PUT` in a `azapi_resource_action` resource is used to continually update the subscription management group association, avoiding a `DELETE` and another `PUT`, which would result in an initial subscription move back to the default management group.
- Artificially managing the lifecycle of the subscription management group association using a data source, which is used to recreate the association if it has been moved outside of Terraform.
  - This involves an idempotency issue on the second run, there after it will behave as expected.
