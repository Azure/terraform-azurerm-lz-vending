<!-- markdownlint-disable MD041 -->
Here is a simple example of deploying a landing zone with a hub & spoke peering to a hub network:

```terraform
resource "azurerm_resource_group" "example" {
  name     = "rg-hub"
  location = "northeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "hubvnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.0.0/23"]
}

module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  location = "northeurope"

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "mysub"
  subscription_alias_name    = "mysub"
  subscription_workload      = "DevTest"

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks = {
    vnet1 = {
      name                    = "spoke"
      address_space           = ["192.168.1.0/24"]
      resource_group_name     = "rg-networking"
      hub_peering_enabled     = true
      hub_network_resource_id = azurerm_virtual_network.example.id
    }
  }
}
```

Back to [Examples](Examples)
