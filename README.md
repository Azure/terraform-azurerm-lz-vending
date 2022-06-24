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

| Name | Source | Version |
|------|--------|---------|
| <a name="module_subscription"></a> [subscription](#module\_subscription) | ./modules/subscription | n/a |
| <a name="module_virtual_network"></a> [virtual\_network](#module\_virtual\_network) | ./modules/virtualnetwork | n/a |

<!-- markdownlint-disable MD013 -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_hub_network_resource_id"></a> [hub\_network\_resource\_id](#input\_hub\_network\_resource\_id) | The resource ID of the virtual network in the hub to which the created virtual network will be peered.<br><br>    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`<br><br>    Leave blank to create the virtual network without peering. | `string` | `""` | no |
| <a name="input_location"></a> [location](#input\_location) | The location of resources deployed by this module. | `string` | `""` | no |
| <a name="input_subscription_alias_billing_scope"></a> [subscription\_alias\_billing\_scope](#input\_subscription\_alias\_billing\_scope) | The billing scope for the new subscription alias.<br><br>  A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_display_name"></a> [subscription\_alias\_display\_name](#input\_subscription\_alias\_display\_name) | The display name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_enabled"></a> [subscription\_alias\_enabled](#input\_subscription\_alias\_enabled) | Whether the creation of a new subscripion alias is enabled or not.<br><br>  If it is disabled, the `subscription_id` variable must be supplied instead. | `bool` | `false` | no |
| <a name="input_subscription_alias_management_group_id"></a> [subscription\_alias\_management\_group\_id](#input\_subscription\_alias\_management\_group\_id) | The destination management group ID for the new subscription.<br><br>  **Note:** Do not supply the display name.<br>  The management group ID forms part of the Azure resource ID. E.g.,<br>  `/providers/Microsoft.Management/managementGroups/{managementGroupId}`. | `string` | `""` | no |
| <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name) | The name of the subscription alias.<br><br>  The string must be comprised of a-z, A-Z, 0-9, - and \_.<br>  The maximum length is 63 characters.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_alias_workload"></a> [subscription\_alias\_workload](#input\_subscription\_alias\_workload) | The billing scope for the new subscription alias.<br><br>  The workload type can be either `Production` or `DevTest` and is case sensitive.<br><br>  You may also supply an empty string if you do not want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `false` and `subscription_id` must be supplied. | `string` | `""` | no |
| <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id) | An existing subscription id.<br><br>  Use this when you do not want the module to create a new subscription.<br><br>  A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.<br>  All letters must be lowercase.<br><br>  You may also supply an empty string if you want to create a new subscription alias.<br>  In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:<br><br>  - `subscription_alias_name`<br>  - `subscription_alias_display_name`<br>  - `subscription_alias_billing_scope`<br>  - `subscription_alias_workload` | `string` | `""` | no |
| <a name="input_virtual_network_address_space"></a> [virtual\_network\_address\_space](#input\_virtual\_network\_address\_space) | The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/16","172.16.0.0/12"]`. | `list(string)` | `[]` | no |
| <a name="input_virtual_network_enabled"></a> [virtual\_network\_enabled](#input\_virtual\_network\_enabled) | Enables and disables the virtual network submodule. | `bool` | `false` | no |
| <a name="input_virtual_network_location"></a> [virtual\_network\_location](#input\_virtual\_network\_location) | The location of the virtual network.<br><br>    Use this to override the default location defined by `var.location`.<br>    Leave blank to use the default location. | `string` | `""` | no |
| <a name="input_virtual_network_name"></a> [virtual\_network\_name](#input\_virtual\_network\_name) | The name of the virtual network. | `string` | `""` | no |
| <a name="input_virtual_network_resource_group_name"></a> [virtual\_network\_resource\_group\_name](#input\_virtual\_network\_resource\_group\_name) | The name of the resource group to create the virtual network in. | `string` | `""` | no |
| <a name="input_virtual_network_use_remote_gateways"></a> [virtual\_network\_use\_remote\_gateways](#input\_virtual\_network\_use\_remote\_gateways) | Enables the use of remote gateways for the virtual network.<br><br>    Applies to both hub and spoke (vnet peerings) as well as virtual WAN connections. | `bool` | `true` | no |
| <a name="input_virtual_network_vwan_propagated_routetables_labels"></a> [virtual\_network\_vwan\_propagated\_routetables\_labels](#input\_virtual\_network\_vwan\_propagated\_routetables\_labels) | The list of virtual WAN labels to advertise the routes to.<br><br>    Leave blank to use the `default` label. | `list(string)` | `[]` | no |
| <a name="input_virtual_network_vwan_propagated_routetables_resource_ids"></a> [virtual\_network\_vwan\_propagated\_routetables\_resource\_ids](#input\_virtual\_network\_vwan\_propagated\_routetables\_resource\_ids) | The list of route table resource ids to advertise routes to.<br><br>    Leave blank to use the `defaultRouteTable.<br>` | `list(string)` | `[]` | no |
| <a name="input_virtual_network_vwan_routetable_resource_id"></a> [virtual\_network\_vwan\_routetable\_resource\_id](#input\_virtual\_network\_vwan\_routetable\_resource\_id) | The resource ID of the virtual network route table to use for the virtual network.<br><br>    Leave blank to use the `defaultRouteTable`.<br><br>    E.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable | `string` | `""` | no |
| <a name="input_vwan_hub_resource_id"></a> [vwan\_hub\_resource\_id](#input\_vwan\_hub\_resource\_id) | The resource ID of the vwan hub to which the virtual network will be connected.<br><br>    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`<br><br>    Leave blank to create a virtual network without a vwan hub connection. | `string` | `""` | no |

## Resources

No resources.

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id) | The subscription\_id is the Azure subscription id that resources have been deployed into. |
| <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id) | The subscription\_resource\_id is the Azure subscription resource id that resources have been deployed into |

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