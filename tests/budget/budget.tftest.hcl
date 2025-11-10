# Tests for the budget module
# Converts the test from tests/budget/budget_test.go

run "budget_scope_subscription" {
  command = plan

  variables {
    budget_name       = "budget"
    budget_scope      = "/subscriptions/00000000-0000-0000-0000-000000000000"
    budget_amount     = 1000
    budget_time_grain = "Monthly"
    budget_time_period = {
      start_date = "2024-01-01T00:00:00Z"
      end_date   = "2025-01-01T00:00:00Z"
    }
    budget_notifications = {
      notification1 = {
        enabled        = true
        operator       = "GreaterThanOrEqualTo"
        threshold      = 50
        threshold_type = "Actual"
        contact_emails = ["email1@example.com", "email2@example.com"]
      }
      notification2 = {
        enabled        = true
        operator       = "GreaterThan"
        threshold      = 75
        threshold_type = "Actual"
        contact_roles  = ["role1", "role2"]
      }
    }
  }

  assert {
    condition     = var.budget_name == "budget"
    error_message = "Budget name should match"
  }

  assert {
    condition     = var.budget_amount == 1000
    error_message = "Budget amount should be 1000"
  }

  assert {
    condition     = var.budget_time_grain == "Monthly"
    error_message = "Budget time grain should be Monthly"
  }

  assert {
    condition     = length(var.budget_notifications) == 2
    error_message = "Should have 2 budget notifications"
  }
}
