<!-- BEGIN_TF_DOCS -->
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

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.10)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 2.2)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD024 -->
## Required Inputs

The following input variables are required:

### <a name="input_location"></a> [location](#input\_location)

Description: The location of the user-assigned managed identity

Type: `string`

### <a name="input_name"></a> [name](#input\_name)

Description: The name of the user managed identity

Type: `string`

### <a name="input_parent_id"></a> [parent\_id](#input\_parent\_id)

Description: The ID of the parent resource to which this user-assigned managed identity.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_federated_credentials_advanced"></a> [federated\_credentials\_advanced](#input\_federated\_credentials\_advanced)

Description: Configure federated identity credentials, using OpenID Connect, for use scenarios outside GitHub Actions and Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `subject_identifier` - The subject of the token.
- `issuer_url` - the URL of the token issuer, should begin with `https://`
- `audiences` - (optional) a set of strings containing the token audiences, defaults to `["api://AzureADTokenExchange"]`.

Type:

```hcl
map(object({
    name               = string
    subject_identifier = string
    audiences          = optional(set(string), ["api://AzureADTokenExchange"])
    issuer_url         = string
  }))
```

Default: `{}`

### <a name="input_federated_credentials_github"></a> [federated\_credentials\_github](#input\_federated\_credentials\_github)

Description: Configure federated identity credentials, using OpenID Connect, for use in GitHub actions.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the GitHub organization, e.g. `Azure` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `repository` - the name of the GitHub respository, e.g. `terraform-azurerm-lz-vending` in `https://github.com/Azure/terraform-azurerm-lz-vending`.
- `entity` - one of 'environment', 'pull\_request', 'tag', or 'branch'
- `enterprise_slug` - the name of the GitHub Enterprise, e.g. `my-enterprise`. This is optional and only valid when using an enterprise.
- `value` - identifies the `entity` type, e.g. `main` when using entity is `branch`. Should be blank when `entity` is `pull_request`.

Type:

```hcl
map(object({
    name            = optional(string)
    organization    = string
    repository      = string
    entity          = string
    enterprise_slug = optional(string)
    value           = optional(string)
  }))
```

Default: `{}`

### <a name="input_federated_credentials_terraform_cloud"></a> [federated\_credentials\_terraform\_cloud](#input\_federated\_credentials\_terraform\_cloud)

Description: Configure federated identity credentials, using OpenID Connect, for use in Terraform Cloud.

The may key is arbitrary and only used for the `for_each` in the resource declaration.

The map value is an object with the following attributes:

- `name` - the name of the federated credential resource, the last segment of the Azure resource id.
- `organization` - the name of the Terraform Cloud organization.
- `project` - the name of the Terraform Cloud project.
- `workspace` - the name of the Terraform Cloud workspace.
- `run_phase` - one of `plan`, or `apply`.

Type:

```hcl
map(object({
    name         = optional(string)
    organization = string
    project      = string
    workspace    = string
    run_phase    = string
  }))
```

Default: `{}`

### <a name="input_tags"></a> [tags](#input\_tags)

Description: The tags to apply to the user-assigned managed identity

Type: `map(string)`

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.umi](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_advanced](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_github_branch](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_github_environment](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_github_pull_request](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_github_tag](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource.umi_federated_credential_terraform_cloud](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

The following outputs are exported:

### <a name="output_client_id"></a> [client\_id](#output\_client\_id)

Description: The client id of the user managed identity

### <a name="output_principal_id"></a> [principal\_id](#output\_principal\_id)

Description: The object id of the user managed identity

### <a name="output_resource_id"></a> [resource\_id](#output\_resource\_id)

Description: The resource id of the user managed identity

### <a name="output_tenant_id"></a> [tenant\_id](#output\_tenant\_id)

Description: The tenant id of the user managed identity

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->