terraform {
  required_version = "~> 1.3"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "~> 1.4"
    }
    modtm = {
      source  = "azure/modtm"
      version = "~> 0.3"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}
