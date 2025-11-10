# Test wrapper for the budget module
module "budget" {
  source = "../../modules/budget"

  budget_name          = var.budget_name
  budget_scope         = var.budget_scope
  budget_amount        = var.budget_amount
  budget_time_grain    = var.budget_time_grain
  budget_time_period   = var.budget_time_period
  budget_notifications = var.budget_notifications
}
