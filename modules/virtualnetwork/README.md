<!-- BEGIN_TF_DOCS -->
# ALZ landing zone virtual network submodule

## Overview

Creates a virtual network in the supplied subscription.
Optionally, created bi-directional peering and/or a virtual WAN connection.

## Notes

See [README.md](../../README.md) in the parent module for more information.

## Example

```terraform
module "virtualnetwork" {
  source  = "Azure/alz-landing-zone/azurerm/modules/virtualnetwork"
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
| <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id) | The subscription ID of the subscription to create the virtual network in. | `string` | n/a | yes |
| <a name="input_virtual_network_address_space"></a> [virtual\_network\_address\_space](#input\_virtual\_network\_address\_space) | The address space of the virtual network, supplied as multiple CIDR blocks, e.g. `["10.0.0.0/16","172.16.0.0/12"]`. | `list(string)` | n/a | yes |
| <a name="input_virtual_network_location"></a> [virtual\_network\_location](#input\_virtual\_network\_location) | The location of the virtual network. | `string` | n/a | yes |
| <a name="input_virtual_network_name"></a> [virtual\_network\_name](#input\_virtual\_network\_name) | The name of the virtual network. | `string` | n/a | yes |
| <a name="input_virtual_network_resource_group_name"></a> [virtual\_network\_resource\_group\_name](#input\_virtual\_network\_resource\_group\_name) | The name of the resource group to create the virtual network in. | `string` | n/a | yes |
| <a name="input_hub_network_resource_id"></a> [hub\_network\_resource\_id](#input\_hub\_network\_resource\_id) | The resource ID of the virtual network in the hub to which the created virtual network will be peered.<br><br>    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualNetworks/my-vnet`<br><br>    Leave blank to create the virtual network without peering. | `string` | `""` | no |
| <a name="input_virtual_network_peering_enabled"></a> [virtual\_network\_peering\_enabled](#input\_virtual\_network\_peering\_enabled) | Whether to enable peering with the supplied hub virtual network.<br>    Enables a hub & spoke networking topology.<br><br>    If enabled the `hub_network_resource_id` must also be suppled. | `bool` | `false` | no |
| <a name="input_virtual_network_use_remote_gateways"></a> [virtual\_network\_use\_remote\_gateways](#input\_virtual\_network\_use\_remote\_gateways) | Enables the use of remote gateways for the virtual network.<br><br>    Applies to both hub and spoke (vnet peerings) as well as virtual WAN connections. | `bool` | `true` | no |
| <a name="input_virtual_network_vwan_connection_enabled"></a> [virtual\_network\_vwan\_connection\_enabled](#input\_virtual\_network\_vwan\_connection\_enabled) | Whether to enable connection with supplied vwan hub.<br>    Enables a vwan networking topology.<br><br>    If enabled the `vwan_hub_resource_id` must also be suppled. | `bool` | `false` | no |
| <a name="input_virtual_network_vwan_propagated_routetables_labels"></a> [virtual\_network\_vwan\_propagated\_routetables\_labels](#input\_virtual\_network\_vwan\_propagated\_routetables\_labels) | The list of virtual WAN labels to advertise the routes to.<br><br>    Leave blank to use the `default` label. | `list(string)` | `[]` | no |
| <a name="input_virtual_network_vwan_propagated_routetables_resource_ids"></a> [virtual\_network\_vwan\_propagated\_routetables\_resource\_ids](#input\_virtual\_network\_vwan\_propagated\_routetables\_resource\_ids) | The list of route table resource ids to advertise routes to.<br><br>    Leave blank to use the `defaultRouteTable.<br>` | `list(string)` | `[]` | no |
| <a name="input_virtual_network_vwan_routetable_resource_id"></a> [virtual\_network\_vwan\_routetable\_resource\_id](#input\_virtual\_network\_vwan\_routetable\_resource\_id) | The resource ID of the virtual network route table to use for the virtual network.<br><br>    Leave blank to use the `defaultRouteTable`.<br><br>    E.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-vhub/hubRouteTables/defaultRouteTable | `string` | `""` | no |
| <a name="input_vwan_hub_resource_id"></a> [vwan\_hub\_resource\_id](#input\_vwan\_hub\_resource\_id) | The resource ID of the vwan hub to which the virtual network will be connected.<br><br>    E.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Network/virtualHubs/my-hub`<br><br>    Leave blank to create a virtual network without a vwan hub connection. | `string` | `""` | no |

## Resources

| Name | Type |
|------|------|
| [azapi_resource.peering](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |
| [azapi_resource.rg](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |
| [azapi_resource.vhubconnection](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |
| [azapi_resource.vnet](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) | resource |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_virtual_network_resource_id"></a> [virtual\_network\_resource\_id](#output\_virtual\_network\_resource\_id) | The created virtual network resource ID |

<!-- markdownlint-enable -->

<!-- END_TF_DOCS -->
