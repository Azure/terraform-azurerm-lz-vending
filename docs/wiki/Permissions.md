<!-- markdownlint-disable MD041 -->
## Permissions required

This module uses two providers, `AzAPI` and `AzureRM`.
We recommend that you use the same identity for both providers.
See [provider configuration](Provider-configuration) for more information.

### Subscription sub-module

This sub-module manages the following resources:

- [`azurerm_subscription`][azurerm_subscription]
- [`azurerm_management_group_subscription_association`][azurerm_management_group_subscription_association]

The identity used must have permission to:

- Create subscriptions using the `Microsoft.Subscription/aliases` resource.
See the [documentation][programatically_create_subscription] for details.

> **Note:** The following process explains how to [assign EA roles to SPNs][assign_ea_roles_to_spns].

- Manage the subscription's management group using the `Microsoft.Management/managementGroups` resource.
For a detailed explanation of the permissions required, see the [documentation][moving_management_groups].

> **Note**: the identity that creates the subscription will have `Owner` permissions assigned by default.
> If you instead supply an existing subscription id, you must ensure that the identity of the provider has the `Owner` permissions assigned.

### Virtual network sub-module

This sub-module manages the following resources using the AzAPI provider:

- [`Microsoft.Network/virtualHubs/hubVirtualNetworkConnections`][hubVirtualNetworkConnections]
- [`Microsoft.Network/virtualNetworks/virtualNetworkPeerings`][virtualnetworkpeerings]
- [`Microsoft.Network/virtualNetworks`][virtualNetworks]
- [`Microsoft.Resources/resourceGroups`][resourceGroups]

These resources are deployed into the new or the supplied subscription.
The identity of the AzAPI provider must have permission to create these resources.

#### Hub virtual network peering

The identity assigned to the AzAPI provider must also have the following permissions on hub networks to create virtual network peerings.
We recommend that you create a custom role in order to maintain the least privilege principle.

| Action | Name |
| - | - |
| `Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write` | Required to create a peering from the supplied hub network. |
| `Microsoft.Network/virtualNetworks/peer/action` | Required to create a peering from the supplied hub network. |
| `Microsoft.Network/virtualNetworks/virtualNetworkPeerings/read` | Read a virtual network peering |
| `Microsoft.Network/virtualNetworks/virtualNetworkPeerings/delete` | Delete a virtual network peering |

See the [documentation](https://learn.microsoft.com/azure/virtual-network/virtual-network-manage-peering?tabs=peering-portal#permissions) for more information.

#### Azure vWAN hub virtual network connection

The identity assigned to the AzAPI provider must also have the following permissions on hub networks to create virtual network connections.
We recommend that you create a custom tole in order to maintain the least privilege principle.

> TBC

### Role assignments sub-module

This sub-module manages the following resources using the AzureRM provider:

- [`azurerm_role_assignment`][azurerm_role_assignment]

The role assignments are deployed into either the new or the supplied subscription, at subscription or child scopes.
The identity of the AzureRM provider must have permission to create these resources, typically this means having the `Owner` or `User Access Administrator` roles.

[assign_ea_roles_to_spns]: https://docs.microsoft.com/azure/cost-management-billing/manage/assign-roles-azure-service-principals
[azurerm_management_group_subscription_association]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group_subscription_association
[azurerm_role_assignment]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/role_assignment
[azurerm_subscription]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subscription
[hubVirtualNetworkConnections]: https://docs.microsoft.com/azure/templates/microsoft.network/virtualhubs/hubvirtualnetworkconnections?tabs=bicep
[moving_management_groups]: https://docs.microsoft.com/azure/governance/management-groups/overview#moving-management-groups-and-subscriptions
[programatically_create_subscription]: https://docs.microsoft.com/azure/cost-management-billing/manage/programmatically-create-subscription
[resourceGroups]: https://docs.microsoft.com/azure/templates/microsoft.resources/resourceGroups?tabs=bicep
[virtualnetworkpeerings]: https://docs.microsoft.com/azure/templates/microsoft.network/virtualnetworks/virtualnetworkpeerings?tabs=bicep
[virtualNetworks]: https://docs.microsoft.com/azure/templates/microsoft.network/virtualnetworks?tabs=bicep
