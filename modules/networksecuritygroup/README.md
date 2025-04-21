<!-- BEGIN_TF_DOCS -->
# Landing zone network security group submodule

## Overview

Creates multiple network security groups in the supplied subscription.
Optionally:

- Creates security rules within the network security groups.

## Notes

See [README.md](https://github.com/Azure/terraform-azurerm-lz-vending#readme) in the parent module for more information.

## Example

See documentation for optional parameters.

```terraform
module "networksecuritygroup" {
  source  = "./"

  subscription_id = "00000000-0000-0000-0000-000000000000"
  resource_group_name = "test-rg"
  name                = "test-nsg1"
  location            = var.location

  security_rules = {
    "rule01" = {
      name                       = "test-r1"
      access                     = "Deny"
      destination_address_prefix = "*"
      destination_port_range     = "80-88"
      direction                  = "Outbound"
      priority                   = 100
      protocol                   = "Tcp"
      source_address_prefix      = "*"
      source_port_range          = "*"
    }
    "rule02" = {
      name                       = "test-r2"
      access                     = "Allow"
      destination_address_prefix = "*"
      destination_port_ranges    = ["80", "443"]
      direction                  = "Inbound"
      priority                   = 200
      protocol                   = "Tcp"
      source_address_prefix      = "*"
      source_port_range          = "*"
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


## Optional Inputs


## Outputs

<!-- markdownlint-enable -->
<!-- END_TF_DOCS -->