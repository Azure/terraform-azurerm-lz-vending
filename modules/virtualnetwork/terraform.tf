terraform {
  required_version = ">= 1.0.0"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = ">= 0.3.0"
    }
  }
}
