# Landing zone user managed identity submodule

## Overview

Creates a user managed identity and (optionally) a resource group in the supplied subscription and creates role assignments for the identity at the supplied scopes, with the supplied role definitions.

Can also configure federated identity credentials to support OpenID Connect (OIDC) authentication, typically for use in GitHub Actions/Terraform cloud workflows.

Outputs useful values for use in other modules, e.g. assigning the identity to another Azure resource.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "umi" {
  source  = "Azure/lz-vending/azurerm/modules/usermanagedidentity"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  location            = "eastus"
  name                = "umi-1"
  resource_group_name = "rg-umi"
  subscription_id     = "00000000-0000-0000-0000-000000000000"

  # optional creation of federated identity credentials
  # for GitHub Actions
  federated_credentials_github = {
    gh1 = {
      organization = "Azure"
      repository   = "terraform-azurerm-lz-vending"
      entity       = "branch"
      value        = "main"
    }
  }

  # optional creation of federated identity credentials
  # for Terraform Cloud
  federated_credentials_terraform_cloud = {
    tfc1 = {
      organization = "myorg"
      project      = "myproject"
      workspace    = "myworkspace"
      run_phase    = "apply"
    }
  }

  # optional creation of federated identity credentials
  # for advanced scenarios
  federated_credentials_advance = {
    adv1 = {
      name               = "custom"
      audience           = "custom"
      issuer_url         = "https://custom"
      subject_identifier = "custom"
    }
  }
}
```
