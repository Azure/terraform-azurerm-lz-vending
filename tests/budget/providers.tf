terraform {
  required_version = "~> 1.9"
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.2"
    }
  }
}

provider "azapi" {
  # Configuration will be inherited from environment
}
