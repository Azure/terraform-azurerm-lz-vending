module "lz_vending" {
  source                     = "../../"
  subscription_alias_enabled = true
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/7690848/enrollmentAccounts/321984"
  subscription_alias_name    = "ratest"
  subscription_display_name  = "ratest"
  subscription_workload      = "DevTest"
  location                   = "swedencentral"
  role_assignment_enabled    = true
  role_assignments = {
    "owner" = {
      principal_type            = "User"
      definition_lookup_enabled = true
      definition                = "Owner"
      principal_id              = "780161f8-01ab-4b08-bf7e-2d76de974d84"
      use_random_uuid           = true
    }
  }

  resource_group_creation_enabled = true
  resource_groups = {
    "umi" = {
      name     = "rg-umi"
      location = "swedencentral"
    }
  }

  umi_enabled = true
  user_managed_identities = {
    "umi1" = {
      resource_group_creation_enabled = false
      location                        = "swedencentral"
      resource_group_name             = "rg-umi"
      name                            = "umi1"
      role_assignments = {
        "test" = {
          definition_lookup_enabled = true
          definition                = "Reader"
          use_random_uuid           = true
        }
      }
    }
  }
  disable_telemetry = true
}
