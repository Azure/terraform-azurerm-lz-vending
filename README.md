<!-- BEGIN_TF_DOCS -->
# Terraform landing zone vending module for Azure

## Overview

The landing zone Terraform module is designed to accelerate deployment of the individual landing zones within an Azure tenant.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Hub & spoke networking
- Virtual WAN networking
- Role assignments

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

```terraform
module "alz_landing_zone" {
  source  = "Azure/lz-vending/azurerm"
  version = "~>0.1.0"

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_location            = "eastus"
  virtual_network_resource_group_name = "my-network-rg"

  # virtual network peering
  virtual_network_peering_enabled = true
  hub_network_resource_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"

  # role assignments
  role_assignment_enabled = true
  role_assignments = [
    # using role definition name, created at subscription scope
    {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Contributor"
      relative_scope = ""
    },
    # using a custom role definition
    {
      principal_id   = "11111111-1111-1111-1111-111111111111"
      definition     = "/providers/Microsoft.Management/MyMg/providers/Microsoft.Authorization/roleDefinitions/ffffffff-ffff-ffff-ffff-ffffffffffff"
      relative_scope = ""
    },
    # using relative scope (to the created or supplied subscription)
    {
      principal_id   = "00000000-0000-0000-0000-000000000000"
      definition     = "Owner"
      relative_scope = "/resourceGroups/MyRg"
    },
  ]
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.0.0)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (>= 0.3.0)

## Modules

The following Modules are called:

### <a name="module_roleassignment"></a> [roleassignment](#module\_roleassignment)

Source: ./modules/roleassignment

Version:

### <a name="module_subscription"></a> [subscription](#module\_subscription)

Source: ./modules/subscription

Version:

### <a name="module_virtualnetwork"></a> [virtualnetwork](#module\_virtualnetwork)

Source: ./modules/virtualnetwork

Version:

<!-- markdownlint-disable MD013 -->
## Required Inputs

No required inputs.

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_hub_network_resource_id"></a> [hub\_network\_resource\_id](#input\_hub\_network\_resource\_id)

Description: The resource ID of the virtual network in the hub to which the created virtual network will be peered.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`

Leave blank to create the virtual network without peering.

Type: `string`

Default: `""`

### <a name="input_location"></a> [location](#input\_location)

Description: The location of resources deployed by this module.

Type: `string`

Default: `""`

### <a name="input_role_assignment_enabled"></a> [role\_assignment\_enabled](#input\_role\_assignment\_enabled)

Description: Whether to create role assignments.  
If enabled, supply the list of role assignments in `var.role_assignments`.

Type: `bool`

Default: `false`

### <a name="input_role_assignments"></a> [role\_assignments](#input\_role\_assignments)

Description: Supply a list of objects containing the details of the role assignments to create.

Object fields:

- `principal_id`: The directory/object id of the principal to assign the role to.
- `definition`: The role definition to assign. Either use the name or the role definition resource id.
- `relative_scope`: Scope relative to the created subscription. Leave blank for subscription scope.

E.g.

```terraform
role_assignments = [
  # Example using role definition name:
  {
    principal_id   = "00000000-0000-0000-0000-000000000000",
    definition     = "Contributor",
    relative_scope = "",
  },
  # Example using role definition id and RG scope:
  {
    principal_id   = "11111111-1111-1111-1111-111111111111",
    definition     = "/providers/Microsoft.Management/managementGroups/mymg/providers/Microsoft.Authorization/roleDefinitions/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    relative_scope = "/resourceGroups/MyRg",
  }
]
```

Type:

```hcl
list(object({
    principal_id   = string,
    definition     = string,
    relative_scope = string,
  }))
```

Default: `[]`

### <a name="input_subscription_alias_enabled"></a> [subscription\_alias\_enabled](#input\_subscription\_alias\_enabled)

Description: Whether to create a new subscription using the subscription alias resource.

If enabled, the following must also be supplied:

- `subscription_alias_name`
- `subscription_display_name`
- `subscription_billing_scope`
- `subscription_workload`

Optionally, supply the following to enable the placement of the subscription into a management group:

- `subscription_management_group_id`
- `subscription_management_group_association_enabled`

If disabled, supply the `subscription_id` variable instead.

Type: `bool`

Default: `false`

### <a name="input_subscription_alias_name"></a> [subscription\_alias\_name](#input\_subscription\_alias\_name)

Description: The name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, - and \_.  
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_billing_scope"></a> [subscription\_billing\_scope](#input\_subscription\_billing\_scope)

Description: The billing scope for the new subscription alias.

A valid billing scope starts with `/providers/Microsoft.Billing/billingAccounts/` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_display_name"></a> [subscription\_display\_name](#input\_subscription\_display\_name)

Description: The display name of the subscription alias.

The string must be comprised of a-z, A-Z, 0-9, -, \_ and space.  
The maximum length is 63 characters.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: An existing subscription id.

Use this when you do not want the module to create a new subscription.  
But do want to manage the management group membership.

A GUID should be supplied in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.  
All letters must be lowercase.

When using this, `subscription_management_group_association_enabled` should be enabled,  
and `subscription_management_group_id` should be supplied.

You may also supply an empty string if you want to create a new subscription alias.  
In this scenario, `subscription_alias_enabled` should be set to `true` and the following other variables must be supplied:

- `subscription_alias_name`
- `subscription_alias_display_name`
- `subscription_alias_billing_scope`
- `subscription_alias_workload`

Type: `string`

Default: `""`

### <a name="input_subscription_management_group_association_enabled"></a> [subscription\_management\_group\_association\_enabled](#input\_subscription\_management\_group\_association\_enabled)

Description: Whether to create the `azurerm_management_group_association` resource.

If enabled, the `subscription_management_group_id` must also be supplied.

Type: `bool`

Default: `false`

### <a name="input_subscription_management_group_id"></a> [subscription\_management\_group\_id](#input\_subscription\_management\_group\_id)

Description:   The destination management group ID for the new subscription.

**Note:** Do not supply the display name.  
The management group ID forms part of the Azure resource ID. E.g.,
`/providers/Microsoft.Management/managementGroups/{managementGroupId}`.

Type: `string`

Default: `""`

### <a name="input_subscription_tags"></a> [subscription\_tags](#input\_subscription\_tags)

Description: A map of tags to assign to the newly created subscription.  
Only valid when `subsciption_alias_enabled` is set to `true`.

Example value:

```terraform
subscription_tags = {
  mytag  = "myvalue"
  mytag2 = "myvalue2"
}
```

Type: `map(string)`

Default: `{}`

### <a name="input_subscription_workload"></a> [subscription\_workload](#input\_subscription\_workload)

Description: The billing scope for the new subscription alias.

The workload type can be either `Production` or `DevTest` and is case sensitive.

You may also supply an empty string if you do not want to create a new subscription alias.  
In this scenario, `subscription_enabled` should be set to `false` and `subscription_id` must be supplied.

Type: `string`

Default: `""`

### <a name="input_virtual_network_address_space"></a> [virtual\_network\_address\_space](#input\_virtual\_network\_address\_space)

Description: The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/8","172.16.0.0/12"]`.

Type: `list(string)`

Default: `[]`

### <a name="input_virtual_network_enabled"></a> [virtual\_network\_enabled](#input\_virtual\_network\_enabled)

Description: Enables and disables the virtual network submodule.

Type: `bool`

Default: `false`

### <a name="input_virtual_network_location"></a> [virtual\_network\_location](#input\_virtual\_network\_location)

Description: The location of the virtual network.

Use this to override the default location defined by `var.location`.  
Leave blank to use the default location.

Type: `string`

Default: `""`

### <a name="input_virtual_network_name"></a> [virtual\_network\_name](#input\_virtual\_network\_name)

Description: The name of the virtual network.

Type: `string`

Default: `""`

### <a name="input_virtual_network_peering_enabled"></a> [virtual\_network\_peering\_enabled](#input\_virtual\_network\_peering\_enabled)

Description: Whether to enable peering with the supplied hub virtual network.  
Enables a hub & spoke networking topology.

If enabled the `hub_network_resource_id` must also be suppled.

Type: `bool`

Default: `false`

### <a name="input_virtual_network_resource_group_name"></a> [virtual\_network\_resource\_group\_name](#input\_virtual\_network\_resource\_group\_name)

Description: The name of the resource group to create the virtual network in.

Type: `string`

Default: `""`

### <a name="input_virtual_network_use_remote_gateways"></a> [virtual\_network\_use\_remote\_gateways](#input\_virtual\_network\_use\_remote\_gateways)

Description: Enables the use of remote gateways for the virtual network.

Applies to hub and spoke (vnet peerings).

Type: `bool`

Default: `true`

### <a name="input_virtual_network_vwan_connection_enabled"></a> [virtual\_network\_vwan\_connection\_enabled](#input\_virtual\_network\_vwan\_connection\_enabled)

Description: Whether to enable connection with supplied vwan hub.  
Enables a vwan networking topology.

If enabled the `vwan_hub_resource_id` must also be supplied.

Type: `bool`

Default: `false`

### <a name="input_virtual_network_vwan_propagated_routetables_labels"></a> [virtual\_network\_vwan\_propagated\_routetables\_labels](#input\_virtual\_network\_vwan\_propagated\_routetables\_labels)

Description: The list of virtual WAN labels to advertise the routes to.

Leave blank to use the `default` label.

Type: `list(string)`

Default: `[]`

### <a name="input_virtual_network_vwan_propagated_routetables_resource_ids"></a> [virtual\_network\_vwan\_propagated\_routetables\_resource\_ids](#input\_virtual\_network\_vwan\_propagated\_routetables\_resource\_ids)

Description: The list of route table resource ids to advertise routes to.

Leave blank to use the `defaultRouteTable`.

Type: `list(string)`

Default: `[]`

### <a name="input_virtual_network_vwan_routetable_resource_id"></a> [virtual\_network\_vwan\_routetable\_resource\_id](#input\_virtual\_network\_vwan\_routetable\_resource\_id)

Description: The resource ID of the virtual network route table to use for the virtual network.

Leave blank to use the `defaultRouteTable`.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable`

Type: `string`

Default: `""`

### <a name="input_vwan_hub_resource_id"></a> [vwan\_hub\_resource\_id](#input\_vwan\_hub\_resource\_id)

Description: The resource ID of the vwan hub to which the virtual network will be connected.  
E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`

Leave blank to create a virtual network without a vwan hub connection.

Type: `string`

Default: `""`

## Resources

No resources.

## Outputs

The following outputs are exported:

### <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id)

Description: The subscription\_id is the Azure subscription id that resources have been deployed into.

### <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id)

Description: The subscription\_resource\_id is the Azure subscription resource id that resources have been deployed into

<!-- markdownlint-enable -->
<!-- markdownlint-disable MD041 -->
## Contributing
<!-- markdownlint-enable -->

This project welcomes contributions and suggestions.
Most contributions require you to agree to a Contributor License Agreement (CLA)
declaring that you have the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.opensource.microsoft.com](https://cla.opensource.microsoft.com).

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment).
Simply follow the instructions provided by the bot.
You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Developing the Module

See [DEVELOPER.md](DEVELOPER.md).

## Trademarks

This project may contain trademarks or logos for projects, products, or services.
Authorized use of Microsoft trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
<!-- END_TF_DOCS -->