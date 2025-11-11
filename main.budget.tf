# The budget module creates budgets from the data
# supplied in the var.budgets variable
module "budget" {
<<<<<<< HEAD
  source               = "./modules/budget"
  for_each             = { for k, v in var.budgets : k => v if var.budget_enabled }
  budget_name          = each.value.name
  budget_scope         = each.value.resource_group_key != null ? module.resourcegroup[each.value.resource_group_key].resource_group_resource_id : "${local.subscription_resource_id}${each.value.relative_scope}"
  budget_amount        = each.value.amount
  budget_time_grain    = each.value.time_grain
  budget_notifications = each.value.notifications
  budget_time_period = {
    end_date   = each.value.time_period_end
    start_date = each.value.time_period_start
  }
}
=======
  source   = "./modules/budget"
  for_each = { for k, v in var.budgets : k => v if var.budget_enabled }

  budget_amount     = each.value.amount
  budget_name       = each.key
  budget_scope      = "${local.subscription_resource_id}${each.value.relative_scope}"
  budget_time_grain = each.value.time_grain
  budget_time_period = {
    end_date   = each.value.time_period_end
    start_date = each.value.time_period_start
  }
  budget_notifications = each.value.notifications

  depends_on = [
    module.resourcegroup,
    module.subscription,
    module.usermanagedidentity,
    module.virtualnetwork,
  ]
}
>>>>>>> e2d7552 (fix line endings)
