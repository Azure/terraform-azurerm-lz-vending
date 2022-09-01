<!-- markdownlint-disable MD041 -->
Here is a simple example of deploying a landing zone with a hub & spoke peering to a hub network:

```terraform
resource "azurerm_resource_group" "example" {
  name     = "rg-hub"
  location = var.location
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

  location = var.location

  # subscription variables
  subscription_alias_enabled = var.subscription_alias_enabled
  subscription_billing_scope = var.subscription_billing_scope
  subscription_display_name  = var.subscription_display_name
  subscription_alias_name    = var.subscription_alias_name
  subscription_workload      = var.subscription_workload

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_name                = "spoke"
  virtual_network_resource_group_name = "rg-networking"

  # virtual network peering
  virtual_network_peering_enabled     = true
  virtual_network_use_remote_gateways = false
  hub_network_resource_id             = azurerm_virtual_network.example.id
}
```

Back to [Examples](Examples)
