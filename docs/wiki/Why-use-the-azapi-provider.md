<!-- markdownlint-disable MD041 -->
In this page we go thtough th erationale for selecting the [AzAPI][github_azapi] Terraform provider,
rather than the more common [AzureRM][github_azurerm] provider from Hashicorp.

## A bit about AzureRM

Before we can deploy resources ugin AzureRM, we must declare a provider.
When we configure the provider, we set arguments or environment variables depending on how we want to [authenticate][hashicorp_azurerm_auth_to_azure].
However, we must always use a `subscription_id` and a `tenant_id`.
Even when authenticating using Azure CLI, we must set the subscription context with `az account set --subscription="<subscription_id>"`.

```terraform
provider "azurerm" {
  features {}

  # ... other authentication arguments, depending on configuraiton,
  # but typically we set the following either here or
  # using environment variables:

  subscription_id = "00000000-0000-0000-0000-000000000000"
  tenant_id       = "00000000-0000-0000-0000-000000000000"
}
```

This configuration couples a provider declaration with a subscripiton.
If you want to deploy resources to a different subscription, you must create a new provider.

[comment]: # (Link labels below, please sort a-z, thanks!)

[github_azapi]: https://github.com/Azure/terraform-provider-azapi
[github_azurerm]: https://github.com/hashicorp/terraform-provider-azurerm
[hashicorp_azurerm_auth_to_azure]: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#authenticating-to-azure
