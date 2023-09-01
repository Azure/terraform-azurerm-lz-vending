variables {
  subscription_alias_enabled = true
  subscription_alias_name    = "my-subscription-alias"
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
  subscription_display_name  = "test-subscription-alias"
  subscription_workload      = "Production"
}

provider "azurerm" {
  features {}
}

run "basic_valid" {
  command = plan

  assert {
    condition = azurerm_subscription.this[0].alias == var.subscription_alias_name
    error_message = "Subscription alias name is not correct"
  }
}
