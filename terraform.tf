terraform {
  required_version = "~> 1.3"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "~> 1.4"
    }
  }
}
