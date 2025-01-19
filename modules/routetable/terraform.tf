terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = ">= 1.11.0"
    }
  }
  required_version = ">= 1.5.0"
}
