data "azurerm_client_config" "current" {}

provider "azurerm" {
  features {}
}

module "lz_vending" {
  source          = "../../"
  subscription_id = var.subscription_id
  location        = "westeurope"
  umi_enabled     = true
  user_managed_identities = {
    "default" = {
      name               = "umi-${var.random_hex}"
      resource_group_key = "rg1"
      role_assignments = {
        "blob" = {
          definition     = "Storage Blob Data Contributor"
          relative_scope = "/resourceGroups/rg-${var.random_hex}"
        }
      }
    }
  }
  disable_telemetry                                = true
  resource_group_creation_enabled                  = true
  subscription_register_resource_providers_enabled = true

  resource_groups = {
    rg1 = {
      name     = "rg-${var.random_hex}"
      location = "westeurope"
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

module "lz_vending_rg_existing" {
  source          = "../../"
  subscription_id = var.subscription_id
  location        = "westeurope"
  umi_enabled     = true
  user_managed_identities = {
    "default" = {
      name                         = "umi-0${var.random_hex}"
      resource_group_name_existing = "rg-${var.random_hex}"
      role_assignments = {
        "blob" = {
          definition     = "Storage Blob Data Contributor"
          relative_scope = "/resourceGroups/rg-${var.random_hex}"
        }
      }
    }
  }
  disable_telemetry                                = true
  resource_group_creation_enabled                  = false
  subscription_register_resource_providers_enabled = true

  depends_on = [module.lz_vending]
}
