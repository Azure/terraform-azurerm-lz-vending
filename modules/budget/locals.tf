locals {
  notifications = {
    for key, notification in var.budget_notifications :
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
