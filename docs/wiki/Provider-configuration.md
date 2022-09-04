<!-- markdownlint-disable MD041 -->
The module makes use of both [AzureRM][azurerm_provider] and [AzAPI][azapi_provider] providers.

Also see [required permissions](Permissions) for more information.

## AzureRM provider

The AzureRM provider is used to manage the following resources:

- Subscription alias ([`azurerm_subscription`][azurerm_subscription])
- Subscription management group association ([`azurerm_management_group_subscription_association`][azurerm_management_group_subscription_association])

### AzureRM configuration

The AzureRM provider must be configured with a `subscription_id`.
This must be an EXISTING subscription id, one that has not been created by this module.

The `tenant_id` property must also be set, and should reflect the tenant that any created subscriptions will be associated with.

For more information, including how to authenticate the provider, see the [AzureRM provider documentation][azurerm_provider_docs].

## AzAPI provider

The AzAPI provider is used to manage the following resources:

- Resource locks
- Role assignments
- Virtual network
- Virtual network peering
- Virtual WAN hub virtual network connection

If you use az cli authentication, it is possible to not specify a subscription id for the provider using the `--allow-no-subscriptions` flag of [az login](https://docs.microsoft.com/cli/azure/reference-index?view=azure-cli-latest#az-login).
A subscription id is not required as this module makes use of the `parent_id` property to deploy resources at different scopes.
This is also possible using the GitHub Actions [Azure/login](https://github.com/marketplace/actions/azure-login) action.
However, if you must specify a subscription id, then use an existing subscription id, one that has not been created by this module.

For more information, including how to authenticate the provider, see the [AzureRM provider documentation][azapi_provider_docs].

[azapi_provider]: https://registry.terraform.io/providers/azure/azapi/latest
[azapi_provider_docs]: https://registry.terraform.io/providers/azure/azapi/latest/docs
[azurerm_management_group_subscription_association]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group_subscription_association
[azurerm_provider_docs]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs
[azurerm_provider]: https://registry.terraform.io/providers/hashicorp/azurerm/latest
[azurerm_subscription]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subscription
