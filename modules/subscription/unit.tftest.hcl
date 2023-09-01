# HINT: make sure to run `terraform init` in this directory before running this test

# Default variable values (can be overridden inside the `run` block)
variables {
  subscription_alias_enabled = true
  subscription_alias_name    = "my-subscription-alias"
  subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000"
  subscription_display_name  = "test-subscription-alias"
  subscription_workload      = "Production"
}

# Stop Terraform moaning about the provider not being configured
provider "azurerm" {
  features {}
}

run "basic_valid" {
  command = plan

  assert {
    condition     = azurerm_subscription.this[0].alias == var.subscription_alias_name
    error_message = "Subscription alias name is not correct"
  }

  assert {
    condition     = azurerm_subscription.this[0].billing_scope_id == var.subscription_billing_scope
    error_message = "Subscription billing scope not correct"
  }

  assert {
    condition     = azurerm_subscription.this[0].subscription_name == var.subscription_display_name
    error_message = "Subscription name is not correct"
  }

  assert {
    condition     = azurerm_subscription.this[0].workload == var.subscription_workload
    error_message = "Subscription workload is not correct"
  }
}
