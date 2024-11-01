variable "subscription_billing_scope" {
  type = string
}

variable "subscription_alias_name" {
  type = string
}

variable "subscription_display_name" {
  type = string
}

variable "subscription_workload" {
  type = string
}

variable "subscription_alias_enabled" {
  type = bool
}

variable "subscription_use_azapi" {
  type = bool
}

variable "subscription_dfc_contact_enabled" {
  type = bool
}

variable "subscription_dfc_contact" {
  type = object({
    emails                = optional(string, "")
    phone                 = optional(string, "")
    alert_notifications   = optional(string, "Off")
    notifications_by_role = optional(list(string), [])
  })
}

module "subscription_test" {
  source                           = "../../"
  subscription_alias_name          = var.subscription_alias_name
  subscription_display_name        = var.subscription_display_name
  subscription_workload            = var.subscription_workload
  subscription_billing_scope       = var.subscription_billing_scope
  subscription_alias_enabled       = var.subscription_alias_enabled
  subscription_use_azapi           = var.subscription_use_azapi
  subscription_dfc_contact_enabled = var.subscription_dfc_contact_enabled
  subscription_dfc_contact         = var.subscription_dfc_contact
}

output "subscription_id" {
  value = module.subscription_test.subscription_id
}
