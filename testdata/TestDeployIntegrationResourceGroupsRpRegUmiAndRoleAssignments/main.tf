data "azurerm_client_config" "current" {}

provider "azurerm" {
  features {}
}

module "lz_vending" {
  source                                          = "../../"
  subscription_id                                 = var.subscription_id
  location                                        = "westeurope"
  umi_enabled                                     = true
  umi_name                                        = "umi-${var.random_hex}"
  umi_resource_group_name                         = "rg-umi-${var.random_hex}"
  disable_telemetry                               = true
  resource_group_creation_enabled                 = true
  subscrption_register_resource_providers_enabled = true

  resource_groups = {
    rg1 = {
      name     = "rg-${var.random_hex}"
      location = "westeurope"
    }
  }
  umi_role_assignments = {
    rg1 = {
      definition     = "Storage Blob Data Contributor"
      relative_scope = "/resourceGroups/rg-${var.random_hex}"
    }
  }
  role_assignments = {
    rg1 = {
      definition     = "Storage Blob Data Contributor"
      relative_scope = "/resourceGroups/rg-${var.random_hex}"
      principal_id   = data.azurerm_client_config.current.object_id
    }
  }
}
