<!-- markdownlint-disable MD041 -->
Here are some example configurations that demonstrate the module usage and integration with the [ALZ Terraform module][alz_tf_module].

## Example 1 - Hub & Spoke

Here is a simple example of deploying a hub network and then creating a new landing zone with virtual network peering to the hub:

```terraform
resource "azurerm_resource_group" "example" {
  name     = "rg-hub"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "hubvnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.0.0/23"]
}

module "alz_landing_zone" {
  source = "..."

  location = var.location

  # subscription variables
  subscription_alias_enabled = var.subscription_alias_enabled
  subscription_billing_scope = var.subscription_billing_scope
  subscription_display_name  = var.subscription_display_name
  subscription_alias_name    = var.subscription_alias_name
  subscription_workload      = var.subscription_workload

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_name                = "spoke"
  virtual_network_resource_group_name = "rg-networking"

  # virtual network peering
  virtual_network_peering_enabled     = true
  virtual_network_use_remote_gateways = false
  hub_network_resource_id             = azurerm_virtual_network.example.id
}
```

## Example 2 - Virtual WAN

```terraform
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"
}

module "landing_zone" {
  source = "..."

  location = each.value.location

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
  subscription_display_name  = "mylz"
  subscription_alias_name    = "mylz"
  subscription_workload      = Production

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = ["192.168.1.0/24"]
  virtual_network_name                = "spoke"
  virtual_network_resource_group_name = "rg-networking"

  # virtual network vwan connection
  virtual_network_vwan_connection_enabled = true
  vwan_hub_resource_id                    = azurerm_virtual_hub.example.id
}
```

## Example 3 - Using YAML files for module input

Due to the flexibility provided by Terraform, we can use YAML or JSON files to define the module's input.
Together with `for_each`, this provides a way of scaling the module to multiple landing zones.

Given a directory of YAML files structured as follows:

```yaml
---
name: lz1
workload: Production
location: northeurope
billing_enrollment_account: 123456
management_group_id: Corp
vnet_address_space:
  - "10.0.1.0/24"
  - "192.168.1.0/24"
role_assignments:
  - principal_id: 00000000-0000-0000-0000-000000000000
    definition: Owner
    relative_scope: ''
  - principal_id: 11111111-1111-1111-1111-111111111111
    definition: Reader
    relative_scope: ''
```

Using some HCL, we can search a directory for files that match a pattern and use them as input.
We then create a map of the data in these files:

```terraform
locals {
  # landing_zone_data_dir is the directory containing the YAML files for the landing zones.
  landing_zone_data_dir = "${path.root}/data"

  # landing_zone_files is the list of landing zone YAML files to be processed
  landing_zone_files = fileset(local.landing_zone_data_dir, "landing_zone_*.yaml")

  # landing_zone_data_map is the decoded YAML data stored in a map
  landing_zone_data_map = {
    for f in local.landing_zone_files :
    f => yamldecode(file("${local.landing_zone_data_dir}/${f}"))
  }
}
```

Finally, we use a for_each on the module:

```terraform
# The landing zone module will be called once per landing_zone_*.yaml file
# in the data directory.
module "landing_zone" {
  source   = "..."
  for_each = local.landing_zone_data_map

  location = each.value.location

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/${each.value.billing_enrollment_account}"
  subscription_display_name  = each.value.name
  subscription_alias_name    = each.value.name
  subscription_workload      = each.value.workload

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = each.value.management_group_id

  # virtual network variables
  virtual_network_enabled             = true
  virtual_network_address_space       = each.value.vnet_address_space
  virtual_network_name                = "spoke"
  virtual_network_resource_group_name = "rg-networking"

  # role assignment variables
  role_assignment_enabled = true
  role_assignments        = each.value.role_assignments
}
```

We have provided a working example of this in the [testdata/TestIntegrationWithYaml](https://github.com/Azure/terraform-azurerm-alz-landing-zone/tree/main/testdata/TestIntegrationWithYaml) directory.

[comment]: # (Link labels below, please sort a-z, thanks!)

[comment]: # (Link labels below, please sort a-z, thanks!)

[alz_tf_module]: https://aka.ms/alz/tf
