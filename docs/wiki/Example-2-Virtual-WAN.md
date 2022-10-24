<!-- markdownlint-disable MD041 -->
Here is a simple example of deploying a landing zone with a Virtual WAN connection:

```terraform
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"
}

module "lz_vending" {
  source  = "Azure/lz-vending/azurerm"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  location = each.value.location

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "mylz"
  subscription_alias_name    = "mylz"
  subscription_workload      = Production

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks = {
    vnet1 = {
      name                    = "spoke"
      address_space           = ["192.168.1.0/24"]
      resource_group_name     = "rg-networking"
      vwan_connection_enabled = true
      vwan_hub_resource_id    = azurerm_virtual_hub.example.id
    }
  }
}
```

Back to [Examples](Examples)
