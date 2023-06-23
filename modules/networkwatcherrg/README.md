<!-- BEGIN_TF_DOCS -->
# ALZ landing zone virtual network submodule

## Overview

Creates multiple virtual networks in the supplied subscription.
Optionally:

- Creates bi-directional peering and/or a virtual WAN connection
- Creates peerings between the virtual networks (mesh peering)

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "virtualnetwork" {
  source  = "Azure/lz-vending/azurerm/modules/virtualnetwork"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  virtual_networks = {
    vnet1 = {
      name                = "myvnet"
      address_space       = ["192.168.0.0/24", "10.0.0.0/24"]
      location            = "westeurope"
      resource_group_name = "myrg"
    },
    vnet2 = {
      name                = "myvnet2"
      address_space       = ["192.168.1.0/24", "10.0.1.0/24"]
      location            = "northeurope"
      resource_group_name = "myrg2"
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

Description: The Azure region to deploy resources into.

Type: `string`

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The ID of the subscription to deploy resources into. E.g. `00000000-0000-0000-0000-000000000000`

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_tags"></a> [tags](#input\_tags)

Description: A mapping of tags to assign to the resource.

Type: `map(string)`

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.network_watcher_rg](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

No outputs.

<!-- markdownlint-enable -->

<!-- END_TF_DOCS -->