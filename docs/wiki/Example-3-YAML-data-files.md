<!-- markdownlint-disable MD041 -->
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
virtual_networks:
  vnet1:
    name: my-vnet
    address_space:
      - "10.0.0.0/24"
    resource_group_name: my-rg
role_assignments:
  my_assignment_1:
    principal_id: 00000000-0000-0000-0000-000000000000
    definition: Owner
    relative_scope: ''
  my_assignment_2:
    principal_id: 11111111-1111-1111-1111-111111111111
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
module "lz_vending" {
  source   = "Azure/lz-vending/azurerm"
  version  = "<version>" # change this to your desired version, https://www.terraform.io/language/expressions/version-constraints
  for_each = local.landing_zone_data_map

  location = each.value.location

  # subscription variables
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/${each.value.billing_enrollment_account}"
  subscription_display_name  = each.value.name
  subscription_alias_name    = each.value.name
  subscription_workload      = each.value.workload

  network_watcher_resource_group_enabled = true

  # management group association variables
  subscription_management_group_association_enabled = true
  subscription_management_group_id                  = each.value.management_group_id

  # virtual network variables
  virtual_network_enabled = true
  virtual_networks        = each.value.virtual_networks

  # role assignment variables
  role_assignment_enabled = true
  role_assignments        = each.value.role_assignments
}
```

We have provided a working example of this in the [testdata/TestIntegrationWithYaml](https://github.com/Azure/terraform-azurerm-lz-vending/tree/main/testdata/TestIntegrationWithYaml) directory.

Back to [Examples](Examples)
