terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.7.0"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  #ts:skip=AC_AZURE_0389 skip resource lock check
  name     = "testdeploy-${var.random_hex}"
  location = "northeurope"
}

module "roleassignment_test" {
  source                       = "../../"
  role_assignment_principal_id = data.azurerm_client_config.current.object_id
  role_assignment_definition   = var.role_definition
  role_assignment_scope        = azurerm_resource_group.test.id
}
