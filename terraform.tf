terraform {
  required_version = "~> 1.10"
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.2"
    }
  }
}
