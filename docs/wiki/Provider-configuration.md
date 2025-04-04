<!-- markdownlint-disable MD041 -->
Prior to v5.0, this module used both the [AzureRM][azurerm_provider] and then[AzAPI][azapi_provider] providers.

After v5.0 the module has been refactored to use only the AzAPI provider.
This was done because of the design of the AzureRM v4.0 provider and its mandatory requirement to be supplied with a subscription id at init time.
This does not make sense for a module that is designed to create subscriptions!
Also, this has the benefit of simplifying the module design.

Also see [required permissions](Permissions) for more information.

## AzureRM provider

The AzureRM provider is used to manage the following resources:

- Subscription alias ([`azurerm_subscription`][azurerm_subscription])
- Subscription management group association ([`azurerm_management_group_subscription_association`][azurerm_management_group_subscription_association])

## AzAPI provider

The AzAPI provider is used to manage all resources.

If you use az cli authentication, it is possible to not specify a subscription id for the provider using the `--allow-no-subscriptions` flag of [az login](https://docs.microsoft.com/cli/azure/reference-index?view=azure-cli-latest#az-login).
A subscription id is not required as this module makes use of the `parent_id` property to deploy resources at different scopes.
This is also possible using the GitHub Actions [Azure/login](https://github.com/marketplace/actions/azure-login) action.
However, if you must specify a subscription id, then use an existing subscription id, one that has not been created by this module.

For more information, including how to authenticate the provider, see the [AzureRM provider documentation][azapi_provider_docs].

[azapi_provider]: https://registry.terraform.io/providers/azure/azapi/latest
[azapi_provider_docs]: https://registry.terraform.io/providers/azure/azapi/latest/docs
[azurerm_management_group_subscription_association]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group_subscription_association
[azurerm_provider]: https://registry.terraform.io/providers/hashicorp/azurerm/latest
[azurerm_subscription]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subscription
