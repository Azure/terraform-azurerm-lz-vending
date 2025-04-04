mock_provider "azapi" {}
mock_provider "time" {}

variables {
  subscription_alias_enabled = true
  subscription_alias_name    = "test-subscription-alias"
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
  subscription_display_name  = "test-subscription-alias"
  subscription_workload      = "Production"
  subscription_tags = {
    test-tag  = "test-value"
    test-tag2 = "test-value2"
  }
}

run "existing_with_management_group" {
  command = plan

  variables {
    subscription_management_group_id                  = "00000000-0000-0000-0000-000000000000"
    subscription_management_group_association_enabled = true
    subscription_alias_enabled                        = false
    subscription_id                                   = "00000000-0000-0000-0000-000000000000"
    subscription_update_existing                      = true
  }

  override_data {
    target = data.azapi_resource_list.subscriptions[0]
    values = {
      output = {
        value = [
          {
            subscriptionId = "00000000-0000-0000-0000-000000000000"
          }
        ]
      }
    }
  }
  override_data {
    target = data.azapi_resource_list.subscription_management_group_association[0]
    values = {
      output = {
        value = []
      }
    }
  }
}
