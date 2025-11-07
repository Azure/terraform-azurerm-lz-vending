# The budget module creates budgets from the data
# supplied in the var.budgets variable
module "budget" {
  source               = "./modules/budget"
  for_each             = { for k, v in var.budgets : k => v if var.budget_enabled }
  budget_name          = coalesce(each.value.name, each.key)
  budget_scope         = each.value.resource_group_key != null ? module.resourcegroup[each.value.resource_group_key].resource_group_resource_id : "${local.subscription_resource_id}${each.value.relative_scope}"
  budget_amount        = each.value.amount
  budget_time_grain    = each.value.time_grain
  budget_notifications = each.value.notifications
  budget_time_period = {
    end_date   = each.value.time_period_end
    start_date = each.value.time_period_start
  }
}
