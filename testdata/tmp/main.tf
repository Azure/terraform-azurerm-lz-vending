data "azurerm_client_config" "current" {}

provider "azurerm" {
  features {}
}

module "lz_vending" {
  source = "../../"
  #source                                 = "Azure/lz-vending/azurerm"
  #version                                = "= 3.3.0"
  subscription_id                        = "d5ffd04f-25c8-4494-a5de-4e1c707bf600"
  network_watcher_resource_group_enabled = true
  location                               = "uksouth"
  disable_telemetry                      = true
}
