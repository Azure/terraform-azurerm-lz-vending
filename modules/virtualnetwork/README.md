<!-- BEGIN_TF_DOCS -->
# ALZ landing zone virtual network submodule

## Overview

Creates a virtual network in the supplied subscription.
Optionally, created bi-directional peering and/or a virtual WAN connection.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "virtualnetwork" {
  source  = "Azure/lz-vending/azurerm/modules/virtualnetwork"
  version = "~> 0.1.0"

  subscription_id                     = "00000000-0000-0000-0000-000000000000"
  virtual_network_name                = "my-virtual-network"
  virtual_network_resource_group_name = "my-network-rg"
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_location            = "eastus"

  virtual_network_peering_enabled = true
  hub_network_resource_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-hub-network-rg/providers/Microsoft.Network/virtualNetworks/my-hub-network"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.0.0)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (>= 0.3.0)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription ID of the subscription to create the virtual network in.

Type: `string`

### <a name="input_virtual_network_address_space"></a> [virtual\_network\_address\_space](#input\_virtual\_network\_address\_space)

Description: The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/16","172.16.0.0/12"]`.

Type: `list(string)`

### <a name="input_virtual_network_location"></a> [virtual\_network\_location](#input\_virtual\_network\_location)

Description: The location of the virtual network.

Type: `string`

### <a name="input_virtual_network_name"></a> [virtual\_network\_name](#input\_virtual\_network\_name)

Description: The name of the virtual network.

Type: `string`

### <a name="input_virtual_network_resource_group_name"></a> [virtual\_network\_resource\_group\_name](#input\_virtual\_network\_resource\_group\_name)

Description: The name of the resource group to create the virtual network in.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_hub_network_resource_id"></a> [hub\_network\_resource\_id](#input\_hub\_network\_resource\_id)

Description: The resource ID of the virtual network in the hub to which the created virtual network will be peered.  
The module will fully establish the peering by creating both sides of the peering connection.

You must also set `virtual_network_peering_enabled = true`.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`

Leave blank and set `virtual_network_peering_enabled = false` (the default) to create the virtual network without peering.

Type: `string`

Default: `""`

### <a name="input_virtual_network_peering_enabled"></a> [virtual\_network\_peering\_enabled](#input\_virtual\_network\_peering\_enabled)

Description: Whether to enable peering with the supplied hub virtual network.  
Enables a hub & spoke networking topology.

If enabled the `hub_network_resource_id` must also be suppled.

Type: `bool`

Default: `false`

### <a name="input_virtual_network_resource_lock_enabled"></a> [virtual\_network\_resource\_lock\_enabled](#input\_virtual\_network\_resource\_lock\_enabled)

Description: Enables the deployment of resource locks to the virtual network's resource group.  
Currently only `CanNotDelete` locks are supported.

Type: `bool`

Default: `true`

### <a name="input_virtual_network_use_remote_gateways"></a> [virtual\_network\_use\_remote\_gateways](#input\_virtual\_network\_use\_remote\_gateways)

Description: Enables the use of remote gateways for the virtual network.

Applies to hub and spoke (vnet peerings).

Type: `bool`

Default: `true`

### <a name="input_virtual_network_vwan_associated_routetable_resource_id"></a> [virtual\_network\_vwan\_associated\_routetable\_resource\_id](#input\_virtual\_network\_vwan\_associated\_routetable\_resource\_id)

Description: The resource ID of the virtual network route table to use for the virtual network.

Leave blank to use the `defaultRouteTable`.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable`

Type: `string`

Default: `""`

### <a name="input_virtual_network_vwan_connection_enabled"></a> [virtual\_network\_vwan\_connection\_enabled](#input\_virtual\_network\_vwan\_connection\_enabled)

Description: The resource ID of the vwan hub to which the virtual network will be connected.  
E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`

You must also set `virtual_network_vwan_connection_enabled = true`.

Leave blank to and set `virtual_network_vwan_connection_enabled = false` (the default) to create a virtual network without a vwan hub connection.

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

### <a name="input_vwan_hub_resource_id"></a> [vwan\_hub\_resource\_id](#input\_vwan\_hub\_resource\_id)

Description: The resource ID of the vwan hub to which the virtual network will be connected.

E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`

Leave blank to create a virtual network without a vwan hub connection.

Type: `string`

Default: `""`

## Resources

The following resources are used by this module:

- [azapi_resource.peering](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.rg](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.rg_lock](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.vhubconnection](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.vnet](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_update_resource.vnet](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/update_resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_virtual_network_resource_id"></a> [virtual\_network\_resource\_id](#output\_virtual\_network\_resource\_id)

Description: The created virtual network resource ID

<!-- markdownlint-enable -->

<!-- END_TF_DOCS -->