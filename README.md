<!-- BEGIN_TF_DOCS -->
# Terraform landing zone vending module for Azure

## Overview

The landing zone Terraform module is designed to accelerate deployment of individual landing zones within an Azure tenant.
We use the [AzureRM][azurem\_provider] and [AzAPI][azapi\_provider] providers to create the subscription and deploy the resources in a single `terrafom apply` step.

The module is designed to be instantiated many times, once for each desired landing zone.

This is currently split logically into the following capabilities:

- Subscription creation and management group placement
- Hub & spoke networking
- Virtual WAN networking
- Role assignments

We would like feedback on what's missing in the module.
Please raise an [issue](https://github.com/Azure/terraform-azurerm-lz-vending/issues) if you have any suggestions.

## Change log

Please see the [GitHub releases pages](https://github.com/Azure/terraform-azurerm-lz-vending/releases) for change log information.

## Notes

Please see the content in the [wiki](https://github.com/Azure/terraform-azurerm-lz-vending/wiki) for more detailed information.

## Example

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "my-subscription-display-name"
  subscription_alias_name    = "my-subscription-alias"
  subscription_workload      = "Production"

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "Corp"

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_location            = "eastus"
  virtual_network_name                = "myvnet"
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

[azurem\_provider]: https://registry.terraform.io/providers/hashicorp/azurerm/latest
[azapi\_provider]: https://registry.terraform.io/providers/azure/azapi/latest

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.3.0)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (>= 1.0.0)

- <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) (>= 3.7.0)

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

The following input variables are required:

### <a name="input_virtual_networks"></a> [virtual\_networks](#input\_virtual\_networks)

Description: A map of the virtual networks to create. The map key must be known at the plan stage, e.g. must not be calculated and known only after apply.

### Required fields

- `name`: The name of the virtual network. [required]
- `address_space`: The address space of the virtual network as a list of strings in CIDR format, e.g. ["192.168.0.0/24, 10.0.0.0/24"]. [required]
- `location`: The location of the virtual network. [required]
- `resource_group_name`: The name of the resource group to create the virtual network in. [required]

### Hub network peering values

The following values configure bi-directional hub & spoke peering for the given virtual network.

- `hub_peering_enabled`: Whether to enable hub peering. [optional]
- `hub_network_resource_id`: The resource ID of the hub network to peer with. [optional - but required if hub\_peering\_enabled is `true`]
- `hub_peering_name_tohub`: The name of the peering to the hub network. [optional - leave empty to use calculated name]
- `hub_peering_name_fromhub`: The name of the peering from the hub network. [optional - leave empty to use calculated name]
- `hub_peering_use_remote_gateways`: Whether to use remote gateways for the hub peering. [optional - default true]

### Mesh peering values

Mesh peering is the capability to create a bi-directional peerings between all supplied virtual networks in `var.virtual_networks`.  
Peerings will only be created between virtual networks with the `mesh_peering_enabled` value set to `true`.

- `mesh_peering_enabled`: Whether to enable mesh peering for this virtual network. Must be enabled on more than one virtual network for any peerings to be created. [optional]
- `mesh_peering_allow_forwarded_traffic`: Whether to allow forwarded traffic for the mesh peering. [optional - default false]

### Other peerings

TODO

### Resource group values

A resource group will be created for

- `resource_group_creation_enabled`: Whether to create a resource group for the virtual network. [optional - default true]
- `resource_group_lock_enabled`: Whether to create a `CanNotDelete` resource lock on the resource group. [optional - default true]
- `resource_group_lock_name`: The name of the resource lock. [optional - leave empty to use calculated name]
- `resource_group_tags`: A map of tags to apply to the resource group, e.g. `{ mytag = "myvalue", mytag2 = "myvalue2" }`. [optional - default empty]

### Virtual WAN values

- `vwan_associated_routetable_resource_id`: The resource ID of the route table to associate with the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]
- `vwan_connection_enabled`: Whether to create a connection to a Virtual WAN. [optional - default false]
- `vwan_connection_name`: The name of the connection to the Virtual WAN. [optional - leave empty to use calculated name]
- `vwan_hub_resource_id`: The resource ID of the hub to connect to. [optional - but required if vwan\_connection\_enabled is `true`]
- `vwan_propagated_routetables_labels`: A list of labels of route tables to propagate to the virtual network. [optional - leave empty to use `["default"]`]
- `vwan_propagated_routetables_resource_ids`: A list of resource IDs of route tables to propagate to the virtual network. [optional - leave empty to use `defaultRouteTable` on hub]

### Tags

- `tags`: A map of tags to apply to the virtual network, e.g. `{ mytag = "myvalue", mytag2 = "myvalue2" }`. [optional - default empty]

Type:

```hcl
map(object({
    name                = string
    address_space       = list(string)
    location            = string
    resource_group_name = string

    hub_network_resource_id         = optional(string, "")
    hub_peering_enabled             = optional(string, false)
    hub_peering_name_tohub          = optional(string, "")
    hub_peering_name_fromhub        = optional(string, "")
    hub_peering_use_remote_gateways = optional(bool, true)

    mesh_peering_enabled                 = optional(bool, false)
    mesh_peering_allow_forwarded_traffic = optional(bool, false)

    other_peerings = optional(map(object({
      remote_network_resource_id            = string
      name_inbound                          = optional(string, "")
      name_outbound                         = optional(string, "")
      outbound_only                         = optional(bool, false)
      allow_forwarded_traffic_inbound       = optional(bool, true)
      allow_forwarded_traffic_outbound      = optional(bool, true)
      allow_gateway_transit_inbound         = optional(bool, false)
      allow_gateway_transit_outbound        = optional(bool, false)
      allow_virtual_network_access_inbound  = optional(bool, true)
      allow_virtual_network_access_outbound = optional(bool, true)
      use_remote_gateways_inbound           = optional(bool, false)
      use_remote_gateways_outbound          = optional(bool, false)
    })), {})

    resource_group_creation_enabled = optional(bool, true)
    resource_group_lock_enabled     = optional(bool, true)
    resource_group_lock_name        = optional(string, "")
    resource_group_tags             = optional(map(string), {})

    vwan_associated_routetable_resource_id   = optional(string, "")
    vwan_connection_enabled                  = optional(bool, false)
    vwan_connection_name                     = optional(string, "")
    vwan_hub_resource_id                     = optional(string, "")
    vwan_propagated_routetables_labels       = optional(list(string), [])
    vwan_propagated_routetables_resource_ids = optional(list(string), [])

    tags = optional(map(string), {})
  }))
```

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_disable_telemetry"></a> [disable\_telemetry](#input\_disable\_telemetry)

Description: To disable tracking, we have included this variable with a simple boolean flag.  
The default value is `false` which does not disable the telemetry.  
If you would like to disable this tracking, then simply set this value to true and this module will not create the telemetry tracking resources and therefore telemetry tracking will be disabled.

For more information, see the [wiki](https://aka.ms/lz-vending/tf/telemetry)

E.g.

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  # ... other module variables

  disable_telemetry = true
}
```

Type: `bool`

Default: `false`

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

If disabled, supply the `subscription_id` variable to use an existing subscription instead.

> **Note**: When the subscription is destroyed, this module will try to remove the NetworkWatcherRG resource group using `az cli`.
> This requires the `az cli` tool be installed and authenticated.
> If the command fails for any reason, the provider will attempt to cancel the subscription anyway.

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

E.g.

- For CustomerLed and FieldLed, e.g. MCA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}`
- For PartnerLed, e.g. MPA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}`
- For Legacy EA - `/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/enrollmentAccounts/{enrollmentAccountName}`

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

## Resources

The following resources are used by this module:

- [azapi_resource.telemetry_root](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_subscription_id"></a> [subscription\_id](#output\_subscription\_id)

Description: The subscription\_id is the Azure subscription id that resources have been deployed into.

### <a name="output_subscription_resource_id"></a> [subscription\_resource\_id](#output\_subscription\_resource\_id)

Description: The subscription\_resource\_id is the Azure subscription resource id that resources have been deployed into

<!-- markdownlint-enable -->
<!-- markdownlint-disable MD041 -->
## Telemetry
<!-- markdownlint-enable -->

When you deploy one or more modules using the landing zone vending module, Microsoft can identify the installation of said module with the deployed Azure resources.
Microsoft can correlate these resources used to support the software.
Microsoft collects this information to provide the best experiences with their products and to operate their business.
The telemetry is collected through customer usage attribution.
The data is collected and governed by Microsoft's privacy policies.

If you don't wish to send usage data to Microsoft, details on how to turn it off can be found [here](https://github.com/Azure/terraform-azurerm-lz-vending/wiki/Telemetry).

## Contributing

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

See [DEVELOPER.md](https://github.com/Azure/terraform-azurerm-lz-vending/blob/main/DEVELOPER.md).

## Trademarks

This project may contain trademarks or logos for projects, products, or services.
Authorized use of Microsoft trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
<!-- END_TF_DOCS -->