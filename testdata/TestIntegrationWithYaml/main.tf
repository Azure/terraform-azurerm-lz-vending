# The landing zone module will be called once per landing_zone_*.yaml file
# in the data directory.
module "alz_landing_zone" {
  source   = "../../"
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
