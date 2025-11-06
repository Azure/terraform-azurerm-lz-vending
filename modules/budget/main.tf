resource "azapi_resource" "budget" {
  name      = var.budget_name
  parent_id = var.budget_scope
  type      = "Microsoft.Consumption/budgets@2021-10-01"
  body = {
    properties = {
      amount        = var.budget_amount
      category      = "Cost"
      notifications = local.notifications
      timeGrain     = var.budget_time_grain
      timePeriod = {
        endDate   = var.budget_time_period.end_date
        startDate = var.budget_time_period.start_date
      }
    }
  }
}
