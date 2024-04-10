locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(azurerm_subscription.this[0].subscription_id, jsondecode(azapi_resource.subscription[0].output).properties.subscriptionId, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = coalesce(local.subscription_id_alias, var.subscription_id)
}

locals {
  # Check if subscription is vended.
  is_subscription_vended = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) ? contains(jsondecode(data.azapi_resource_list.subscriptions[0].output).value[*].subscriptionId, local.subscription_id) : true
  # Check for drift between subscription and target management group.
  is_subscription_associated_to_management_group = (var.subscription_management_group_association_enabled && var.subscription_use_azapi) && local.is_subscription_vended ? contains(jsondecode(data.azapi_resource_list.subscription_management_group_association[0].output).value[*].id, "/providers/Microsoft.Management/managementGroups/${var.subscription_management_group_id}/subscriptions/${local.subscription_id}") : true
}

locals {
  # Transform subscription budgets to be able to use them with the API.
  transformed_budgets = {
    for key, budget in var.subscription_budgets :
    key => {
      amount    = budget.amount
      timeGrain = budget.time_grain
      timePeriod = {
        endDate   = budget.time_period_end
        startDate = budget.time_period_start
      }
      notifications = {
        for key, notification in budget.notifications :
        key => {
          enabled       = notification.enabled
          operator      = notification.operator
          threshold     = notification.threshold
          thresholdType = notification.threshold_type
          contactEmails = notification.contact_emails
          contactRoles  = notification.contact_roles
          contactGroups = notification.contact_groups
          locale        = notification.locale
        }
      }
    }
  }
}
