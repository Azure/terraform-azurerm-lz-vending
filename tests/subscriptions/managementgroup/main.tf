terraform {
  required_version = ">= 1.2.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.10.0"
    }
    azapi = {
      source  = "azure/azapi"
      version = ">= 0.3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

variable "subscription_alias_billing_scope" {
  type = string
}

variable "subscription_alias_management_group_id" {
  type = string
}

variable "subscription_alias_name" {
  type = string
}

variable "subscription_alias_display_name" {
  type = string
}

data "azurerm_client_config" "default" {}

resource "azurerm_management_group" "test" {
  name                       = var.subscription_alias_management_group_id
  display_name               = var.subscription_alias_management_group_id
  parent_management_group_id = data.azurerm_client_config.default.tenant_id
}

module "lz_test" {
  source                                 = "../../../../"
  location                               = "northeurope"
  subscription_alias_name                = var.subscription_alias_name
  subscription_alias_display_name        = var.subscription_alias_display_name
  subsciption_alias_workload             = "DevTest"
  subscription_alias_management_group_id = azurerm_management_group.test.id
  subscription_alias_billing_scope       = var.subscription_alias_billing_scope
}
