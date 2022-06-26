resource "azapi_resource" "rg" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  parent_id = "/subscriptions/fd4423da-0d2c-480f-a55e-5453db019939"
  name      = "deploytest-6754asd3"
  location  = "westeurope"
}

resource "azurerm_public_ip" "ip" {
  name                         = "deploytest-6754asd3-pip"
  location                     = "westeurope"
  resource_group_name          = azapi_resource.rg.name
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}
