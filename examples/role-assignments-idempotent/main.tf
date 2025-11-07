terraform {
  required_version = "~> 1.12"
  required_providers {
    azapi = {
      source  = "Azure/azapi"
      version = "~> 2.5"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azapi_client_config" "current" {}

module "lz-vending" {
  source                  = "../../"
  subscription_id         = data.azapi_client_config.current.subscription_id
  location                = "swedencentral"
  role_assignment_enabled = true
  role_assignments = {
    ra1 = {
      definition      = "/providers/Microsoft.Authorization/roleDefinitions/48e5e92e-a480-4e71-aa9c-2778f4c13781" # Azure Batch Job Submitter
      relative_scope  = ""
      principal_id    = data.azapi_client_config.current.object_id
      use_random_uuid = true
    }
    ra2 = {
      definition               = "Storage Blob Data Contributor"
      resource_group_scope_key = "rg1"
      principal_id             = data.azapi_client_config.current.object_id
    }
  }
  disable_telemetry               = true
  resource_group_creation_enabled = true
  resource_groups = {
    rg1 = {
      name = "rg-vending-002"
    }
  }
  umi_enabled = true
  user_managed_identities = {
    umi1 = {
      name               = "umi-vending-001"
      resource_group_key = "rg1"
      role_assignments = {
        stg_blob_rg = {
          definition               = "Storage Blob Data Contributor"
          resource_group_scope_key = "rg1"
        }
        owner_sub = {
          use_random_uuid = true
          definition      = "Owner"
        }
      }
    }
  }
}
