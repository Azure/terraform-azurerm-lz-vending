<!-- markdownlint-disable MD041 -->
It may be desireable to use this module to deploy resources into an existing subscription. This example shows how to do that by supplying the subscription ID as a module input.

In this example we deploy a virtual network to an existing landing zone subscription and manage the management group association.

The management group association is managed by the AzureRM provider, the provider passed to the module must use an identity that has permissions to manage the management group subscription membership.

The use of the AzAPI provider means we do not need a distinct provider declaration for the LZ subscription.
The identity used by the AzAPI provider must have permissions to create the resources in the LZ subscription.

```terraform
module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  location = "northeurope"

  # subscription variables
  subscription_id = "00000000-0000-0000-0000-000000000000"

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = "mymg"

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.2.0/24"]
  virtual_network_name                = "spoke"
  virtual_network_resource_group_name = "rg-networking"
}
```

Back to [Examples](Examples)
