<!-- BEGIN_TF_DOCS -->
# Landing zone route table submodule

## Overview

Creates multiple route tables in the supplied subscription.
Optionally:

- Creates routes within the route tables

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "routetable" {
  source  = "Azure/lz-vending/azurerm/modules/routetable"
  version = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints

  subscription_id = "00000000-0000-0000-0000-000000000000"
  route_tables = {
    rt1 = {
      name                   = "myroutetable"
      address_prefix         = ["192.168.0.0/24"]
      next_hop_in_ip_address = "192.168.0.5"
      next_hop_type          = "VirtualAppliance"
    },
    rt2 = {
      name           = "myroutetable2"
      address_prefix = "GatewayManager"
      next_hop_type  = "Internet"
    }
  }
}
```

## Documentation
<!-- markdownlint-disable MD033 -->

## Requirements

The following requirements are needed by this module:

- <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) (>= 1.5.0)

- <a name="requirement_azapi"></a> [azapi](#requirement\_azapi) (>= 1.11.0)

## Modules

No modules.

<!-- markdownlint-disable MD013 -->
## Required Inputs

The following input variables are required:

### <a name="input_subscription_id"></a> [subscription\_id](#input\_subscription\_id)

Description: The subscription ID of the subscription to create the virtual network in.

Type: `string`

## Optional Inputs

The following input variables are optional (have default values):

### <a name="input_route_tables"></a> [route\_tables](#input\_route\_tables)

Description: A map defining route tables and their associated routes to be created.
  - `name` (required): The name of the route table.
  - `location` (required): The location of the resource group.
  - `resource_group_name` (required): The name of the resource group.
  - `tags` (optional): A map of key-value pairs for tags associated with the route table.
  - `routes` (optional): A map defining routes for the route table. Each route object has the following properties:
      - `name` (required): The name of the route.
      - `address_prefix` (required): The address prefix for the route.
      - `next_hop_type` (required): The type of next hop for the route.
      - `next_hop_in_ip_address` (required): The next hop IP address for the route.

Type:

```hcl
map(object({
    name                = string
    location            = string
    resource_group_name = string
    tags                = optional(map(string))

    routes = optional(map(object({
      name                   = string
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = string
    })))
  }))
```

Default: `{}`

## Resources

The following resources are used by this module:

- [azapi_resource.route_table](https://registry.terraform.io/providers/azure/azapi/latest/docs/resources/resource) (resource)

## Outputs

No outputs.

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->