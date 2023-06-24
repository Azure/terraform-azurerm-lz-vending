<!-- BEGIN_TF_DOCS -->
# ALZ landing zone user managed identity submodule

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

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.3.0)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (>= 1.0.0)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_location"></a> [location](#input\_location)

Description: The location of the user managed identity

Type: `string`

### <a name="input_name"></a> [name](#input\_name)

Description: The name of the user managed identity

Type: `string`

### <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name)

Description: The name of the resource group in which to create the user managed identity

Type: `string`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription id

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_federated_credentials_advanced"></a> [federated\_credentials\_advanced](#input\_federated\_credentials\_advanced)

Description: allow the caller to configure federated credentials by supplying the values verbatim

Type:

```hcl
map(object({
    name               = string
    subject_identifier = string
    audience           = optional(string, "api://AzureADTokenExchange")
    issuer_url         = string
  }))
```

Default: `{}`

### <a name="input_federated_credentials_github"></a> [federated\_credentials\_github](#input\_federated\_credentials\_github)

Description: allow the caller to easily configure federated credentials for GitHub Actions

Type:

```hcl
map(object({
    name         = optional(string, "")
    organization = string
    repository   = string
    entity       = string
    value        = optional(string, "")
  }))
```

Default: `{}`

### <a name="input_federated_credentials_terraform_cloud"></a> [federated\_credentials\_terraform\_cloud](#input\_federated\_credentials\_terraform\_cloud)

Description: allow the caller to easily configure federated credentials for Terraform Cloud

Type:

```hcl
map(object({
    name         = optional(string, "")
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
```

Default: `{}`

### <a name="input_resource_group_creation_enabled"></a> [resource\_group\_creation\_enabled](#input\_resource\_group\_creation\_enabled)

Description: Whether to create the supplied resource group

Type: `bool`

Default: `true`

### <a name="input_resource_group_lock_enabled"></a> [resource\_group\_lock\_enabled](#input\_resource\_group\_lock\_enabled)

Description: Whether to enable resource group lock

Type: `bool`

Default: `true`

### <a name="input_resource_group_lock_name"></a> [resource\_group\_lock\_name](#input\_resource\_group\_lock\_name)

Description: The name of the resource group lock, if blank will be set to `lock-<resource_group_name>`

Type: `string`

Default: `""`

### <a name="input_resource_group_tags"></a> [resource\_group\_tags](#input\_resource\_group\_tags)

Description: The tags to apply to the resource group, if we create it.

Type: `map(string)`

Default: `{}`

### <a name="input_tags"></a> [tags](#input\_tags)

Description: The tags to apply to the user managed identity

Type: `map(string)`

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.rg](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.rg_lock](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credentials](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_client_id"></a> [client\_id](#output\_client\_id)

Description: The client id of the user managed identity

### <a name="output_principal_id"></a> [principal\_id](#output\_principal\_id)

Description: The object id of the user managed identity

### <a name="output_tenant_id"></a> [tenant\_id](#output\_tenant\_id)

Description: The tenant id of the user managed identity

### <a name="output_umi_id"></a> [umi\_id](#output\_umi\_id)

Description: The resource id of the user managed identity

<!-- markdownlint-enable -->

<!-- END_TF_DOCS -->