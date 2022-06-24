<!-- BEGIN_TF_DOCS -->
# ALZ landing zone subscription submodule

## Overview

Creates a subscription alias, optionally in the specified management group.

## Notes

See [README.md](../../README.md) in the parent module for more information.

## Example

```terraform
# TBC
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0.0 |
| <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) | >= 0.3.0 |

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_subscription_alias_billing_scope"></a> [subscription\_alias\_billing\_scope](#input\_subscription\_alias\_billing\_scope) | The billing scope for the new subscription alias.<br><br>  A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_display_name"></a> [subscription\_alias\_display\_name](#input\_subscription\_alias\_display\_name) | The display name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_management_group_id"></a> [subscription\_alias\_management\_group\_id](#input\_subscription\_alias\_management\_group\_id) | The destination management group ID for the new subscription.<br><br>  **Note:** Do not supply the display name.<br>  The management group ID forms part of the Azure resource ID. E.g.,<br>  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`. | `string` | `""` | no |
| <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name) | The name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, - and \_.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_workload"></a> [subscription\_alias\_workload](#input\_subscription\_alias\_workload) | The billing scope for the new subscription alias.<br><br>  The workload type can be either `Production` or `DevTest` and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |

## Resources

| Name | Type |
|------|------|
| [azapi_resource.subscription_alias](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id) | The subscription\_id is the id of the newly created subscription. |
| <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id) | The subscription\_resource\_id output is the Azure resource id for the newly created subscription. |

<!-- markdownlint-enable -->

<!-- END_TF_DOCS -->