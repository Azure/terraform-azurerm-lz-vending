<!-- BEGIN_TF_DOCS -->
# ALZ landing zone module

## Overview

The landing zone Terraform module is designed to accelerate deployment of the individual landing zones into the ALZ conceptual architecture.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- (backlogged) Hub & spoke networking
- (backlogged) Virtual WAN networking
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

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0.0 |
| <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) | >= 0.3.0 |

## Modules

No modules.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_subscription_alias_billing_scope"></a> [subscription\_alias\_billing\_scope](#input\_subscription\_alias\_billing\_scope) | The billing scope for the new subscription alias.<br><br>  A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_display_name"></a> [subscription\_alias\_display\_name](#input\_subscription\_alias\_display\_name) | The display name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_enabled"></a> [subscription\_alias\_enabled](#input\_subscription\_alias\_enabled) | Whether the creation of a new subscripion alias is enabled or not.<br><br>  If it is disabled, the `subscription_id` variable must be supplied instead. | `bool` | `true` | no |
| <a name="input_subscription_alias_management_group_id"></a> [subscription\_alias\_management\_group\_id](#input\_subscription\_alias\_management\_group\_id) | The destination management group ID for the new subscription.<br><br>  **Note:** Do not supply the display name.<br>  The management group ID forms part of the Azure resource ID. E.g.,<br>  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`. | `string` | `""` | no |
| <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name) | The name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, - and \_.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_workload"></a> [subscription\_alias\_workload](#input\_subscription\_alias\_workload) | The billing scope for the new subscription alias.<br><br>  The workload type can be either Production or DevTest and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id) | An existing subscription id.<br><br>  Use this when you do not want the nmodule to create a new subscription.<br><br>  A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.<br>  All letters must be lowercase.<br><br>  You may also supply an empty string if you want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:<br><br>  - `subscription_alias_name`<br>  - `subscription_alias_display_name`<br>  - `subscription_alias_billing_scope`<br>  - `subscription_alias_workload` | `string` | `""` | no |

## Resources

| Name | Type |
|------|------|
| [azapi_resource.subscription_alias](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id) | The subscription\_id output allows other modules to use the generated or supplied subscription id. |

<!-- markdownlint-disable MD041 -->
## Contributing
<!-- markdownline-enable -->

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit [https://cla.opensource.microsoft.com](https://cla.opensource.microsoft.com).

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Developing the Module

See [DEVELOPER.md](DEVELOPER.md).

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft
trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.

<!-- END_TF_DOCS -->