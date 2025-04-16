terraform {
  required_version = "~> 1.10"
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = "~> 2.2"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
  }
}
