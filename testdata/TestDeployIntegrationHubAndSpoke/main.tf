data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "hub" {
  #ts:skip=AC_AZURE_0389 skip resource lock check
  name     = "${var.virtual_networks["primary"].name}-hub"
  location = var.location
}

resource "azurerm_virtual_network" "hub" {
  #ts:skip=AC_AZURE_0356 skip NSG subnet check
  name                = "${var.virtual_networks["primary"].name}-hub"
  location            = azurerm_resource_group.hub.location
  resource_group_name = azurerm_resource_group.hub.name
  address_space       = ["192.168.10.0/23"]
}

locals {
  virtual_network_primary_merged = merge(var.virtual_networks["primary"], {
    hub_network_resource_id = azurerm_virtual_network.hub.id
  })
  virtual_networks_merged = {
    primary = local.virtual_network_primary_merged
  }
}

module "lz_vending" {
  source = "../../"

  location = var.location

  # subscription variables
  subscription_alias_enabled = var.subscription_alias_enabled
  subscription_billing_scope = var.subscription_billing_scope
  subscription_display_name  = var.subscription_display_name
  subscription_alias_name    = var.subscription_alias_name
  subscription_workload      = var.subscription_workload

  # virtual network variables
  virtual_network_enabled = var.virtual_network_enabled
  virtual_networks        = local.virtual_networks_merged

  # role assignment
  role_assignment_enabled = var.role_assignment_enabled
  role_assignments = [
    {
      principal_id   = data.azurerm_client_config.current.object_id
      definition     = "Storage Blob Data Contributor"
      relative_scope = ""
    }
  ]
}
