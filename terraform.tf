terraform {
  required_version = "~> 1.8"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "~> 2.2"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.116, < 5.0"
    }
  }
}
