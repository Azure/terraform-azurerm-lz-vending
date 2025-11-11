# Tests for the subscription module
# These tests cover subscription alias creation and validation scenarios

run "valid_subscription_alias_create" {
  command = plan

  variables {
    subscription_alias_enabled = true
    subscription_alias_name    = "test-subscription-alias"
    subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name  = "Test Subscription"
    subscription_workload      = "Production"
  }

  assert {
    condition     = var.subscription_alias_enabled == true
    error_message = "Subscription alias should be enabled"
  }

  assert {
    condition     = var.subscription_alias_name == "test-subscription-alias"
    error_message = "Subscription alias name should match"
  }
}

run "valid_subscription_with_management_group" {
  command = plan

  variables {
    subscription_alias_enabled                        = true
    subscription_alias_name                           = "test-subscription-alias-mg"
    subscription_billing_scope                        = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name                         = "Test Subscription with MG"
    subscription_workload                             = "Production"
    subscription_management_group_id                  = "mg-test"
    subscription_management_group_association_enabled = true
  }

  assert {
    condition     = var.subscription_alias_enabled == true
    error_message = "Subscription alias should be enabled"
  }

  assert {
    condition     = var.subscription_management_group_association_enabled == true
    error_message = "Management group association should be enabled"
  }
}

# Validation test: invalid billing scope
run "invalid_billing_scope" {
  command = plan

  variables {
    subscription_alias_enabled = true
    subscription_alias_name    = "test-subscription"
    subscription_billing_scope = "invalid-billing-scope" # Must start with /providers/Microsoft.Billing/billingAccounts/
    subscription_display_name  = "Test Subscription"
    subscription_workload      = "Production"
  }

  expect_failures = [
    var.subscription_billing_scope
  ]
}

# Validation test: invalid workload (must be Production or DevTest)
run "invalid_workload" {
  command = plan

  variables {
    subscription_alias_enabled = true
    subscription_alias_name    = "test-subscription"
    subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name  = "Test Subscription"
    subscription_workload      = "InvalidWorkload" # Must be Production or DevTest
  }

  expect_failures = [
    var.subscription_workload
  ]
}

# Validation test: invalid management group ID with invalid characters
run "invalid_management_group_id_chars" {
  command = plan

  variables {
    subscription_alias_enabled                        = true
    subscription_alias_name                           = "test-subscription"
    subscription_billing_scope                        = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name                         = "Test Subscription"
    subscription_workload                             = "Production"
    subscription_management_group_id                  = "mg@invalid!" # Invalid characters (!, @)
    subscription_management_group_association_enabled = true
  }

  expect_failures = [
    var.subscription_management_group_id
  ]
}

# Validation test: invalid management group ID with excessive length
run "invalid_management_group_id_length" {
  command = plan

  variables {
    subscription_alias_enabled                        = true
    subscription_alias_name                           = "test-subscription"
    subscription_billing_scope                        = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name                         = "Test Subscription"
    subscription_workload                             = "Production"
    subscription_management_group_id                  = "a123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789" # >90 chars
    subscription_management_group_association_enabled = true
  }

  expect_failures = [
    var.subscription_management_group_id
  ]
}

# Validation test: invalid tag value (>256 characters)
run "invalid_tag_value" {
  command = plan

  variables {
    subscription_alias_enabled = true
    subscription_alias_name    = "test-subscription"
    subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name  = "Test Subscription"
    subscription_workload      = "Production"
    subscription_tags = {
      test-tag = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" # >256 chars
    }
  }

  expect_failures = [
    var.subscription_tags
  ]
}

# Validation test: invalid tag name with special characters
run "invalid_tag_name" {
  command = plan

  variables {
    subscription_alias_enabled = true
    subscription_alias_name    = "test-subscription"
    subscription_billing_scope = "/providers/Microsoft.Billing/billingAccounts/1234567/enrollmentAccounts/123456"
    subscription_display_name  = "Test Subscription"
    subscription_workload      = "Production"
    subscription_tags = {
      "invalid<tag>" = "value" # Tag name contains <> which are not allowed
    }
  }

  expect_failures = [
    var.subscription_tags
  ]
}
