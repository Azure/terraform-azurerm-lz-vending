<!-- markdownlint-disable MD041 -->
In this page we go through the rationale for selecting the [AzAPI][github_azapi] Terraform provider,
rather than the more common [AzureRM][github_azurerm] provider from Hashicorp.

## A bit about AzureRM

Before we can deploy resources using AzureRM, we must declare a provider.
When we configure the provider, we set arguments or environment variables depending on how we want to [authenticate][hashicorp_azurerm_auth_to_azure].
However, we must always use a `subscription_id` and a `tenant_id`.
Even when authenticating using Azure CLI, we must set the subscription context with `az account set --subscription="<subscription_id>"`.

```terraform
provider "azurerm" {
  features {}

  # ... other authentication arguments, depending on configuration,
  # but typically we set the following either here or
  # using environment variables:

  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
}
```

This configuration couples a provider declaration with a subscription.
If you want to deploy resources to a different subscription, you must create a new provider.

## AzureRM in a multi-subscription environment

In an environment where we are deploying to multiple subscriptions, we must create a provider for each subscription.
However, since we are creating the subscription using Terraform, the subscription id is not known until after the apply.
This presents a problem in that the Terraform providers must be configured during the init phase.

We then have a circular dependency.
We cannot create the provider block for the new subscription as we do not know the subscription id until after we have deployed the resource.

## AzAPI to the rescue

The AzAPI provider allows us to specify a `parent_id`, this is the deployment scope of the resource.
This can be any scope in Azure and is not bound to a specific subscription.

Therefore the provider can be used to deploy resources to any subscription in the tenant, without having to declare a new provider block.

## Drawbacks of AzAPI

One of the drawbacks of AzAPI is that we have to be very careful with idempotency.
Some Azure resource providers accept properties in the PUT, but do not return those properties upon GET.
This creates idempotency issues where Terraform will detect that the resource has changed and try to remediate the missing properties.

[comment]: # (Link labels below, please sort a-z, thanks!)

[github_azapi]: https://github.com/Azure/terraform-provider-azapi
[github_azurerm]: https://github.com/hashicorp/terraform-provider-azurerm
[hashicorp_azurerm_auth_to_azure]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure
