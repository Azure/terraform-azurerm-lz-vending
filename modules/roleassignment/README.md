<!-- BEGIN_TF_DOCS -->
# Landing zone role assignment submodule

## Overview

Creates a role assignment at subscription or lower scope.
Module is designed to be instantiated many times, once per role assignment.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

```terraform
module "roleassignment" {
  source          = "Azure/lz-vending/azurerm/modules/roleassignment"
  version         = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints
  role_definition = "Owner"
  scope           = "/subscriptions/00000000-0000-0000-0000-000000000000"
  principal_id    = "00000000-0000-0000-0000-000000000000"
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (~> 1.10)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (~> 2.2)

- <a name="requirement_time"></a> [time](#requirement\_time) (~> 0.9)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD024 -->
## Required Inputs

The following input variables are required:

### <a name="input_role_assignment_definition"></a> [role\_assignment\_definition](#input\_role\_assignment\_definition)

Description: Either the role definition resource id, e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c`.  
Or, the role definition name, e.g. `Contributor`.

Type: `string`

### <a name="input_role_assignment_principal_id"></a> [role\_assignment\_principal\_id](#input\_role\_assignment\_principal\_id)

Description: The principal (object) ID of the role assignment.  
Note, for a service principal, this is not the application id.

Can be user, group or service principal.

Type: `string`

### <a name="input_role_assignment_scope"></a> [role\_assignment\_scope](#input\_role\_assignment\_scope)

Description: The scope of the role assignment.

Must begin with `/subscriptions/{subscription-id}` to avoid accidentally creating a role assignment at higher scopes.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_role_assignment_condition"></a> [role\_assignment\_condition](#input\_role\_assignment\_condition)

Description: (Optional) The condition that limits the resources that the role can be assigned to.

Type: `string`

Default: `null`

### <a name="input_role_assignment_condition_version"></a> [role\_assignment\_condition\_version](#input\_role\_assignment\_condition\_version)

Description: The version of the condition. Possible values are `null`, 1.0 or 2.0. If `null` then `role_assignment_condition` will also be null.

Type: `string`

Default: `null`

### <a name="input_role_assignment_definition_lookup_enabled"></a> [role\_assignment\_definition\_lookup\_enabled](#input\_role\_assignment\_definition\_lookup\_enabled)

Description: Whether to look up the role definition resource id from the role definition name.  
If disabled, the `role_assignment_definition` must be a role definition resource id.

Type: `bool`

Default: `true`

### <a name="input_role_assignment_principal_type"></a> [role\_assignment\_principal\_type](#input\_role\_assignment\_principal\_type)

Description: Required when using attribute based access control (ABAC).  
The type of principal. Can be `User`, `Group`, `ServicePrincipal`, `Device`, or `ForeignGroup`.

Type: `string`

Default: `null`

## Resources

The following resources are used by this module:

- [azapi_resource.this](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) (resource)
- [azapi_resource_list.role_definitions](https://registry.terraform.io/providers/Azure/azapi/latest/docs/data-sources/resource_list) (data source)

## Outputs

The following outputs are exported:

### <a name="output_role_assignment_name"></a> [role\_assignment\_name](#output\_role\_assignment\_name)

Description: The Azure name (uuid) of the created role assignment.

### <a name="output_role_assignment_resource_id"></a> [role\_assignment\_resource\_id](#output\_role\_assignment\_resource\_id)

Description: The Azure resource id of the created role assignment.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->