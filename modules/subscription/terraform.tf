terraform {
  required_version = ">= 1.4.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.7"
    }
    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.13"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
  }
}
