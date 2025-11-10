# Integration tests wrapper for root module

module "lz_vending" {
  source = "../../"

  location                           = var.location
  subscription_alias_enabled         = var.subscription_alias_enabled
  subscription_display_name          = var.subscription_display_name
  subscription_alias_name            = var.subscription_alias_name
  subscription_workload              = var.subscription_workload
  subscription_tags                  = var.subscription_tags
  subscription_billing_scope         = var.subscription_billing_scope
  subscription_id                    = var.subscription_id
  resource_group_creation_enabled    = var.resource_group_creation_enabled
  resource_groups                    = var.resource_groups
  virtual_network_enabled            = var.virtual_network_enabled
  virtual_networks                   = var.virtual_networks
  role_assignment_enabled            = var.role_assignment_enabled
  role_assignments                   = var.role_assignments
  route_table_enabled                = var.route_table_enabled
  route_tables                       = var.route_tables
  enable_telemetry                   = var.enable_telemetry
}

variable "location" {
  type    = string
  default = "northeurope"
}

variable "subscription_alias_enabled" {
  type    = bool
  default = false
}

variable "subscription_display_name" {
  type    = string
  default = ""
}

variable "subscription_alias_name" {
  type    = string
  default = ""
}

variable "subscription_workload" {
  type    = string
  default = "Production"
}

variable "subscription_tags" {
  type    = map(string)
  default = {}
}

variable "subscription_billing_scope" {
  type    = string
  default = ""
}

variable "subscription_id" {
  type    = string
  default = "00000000-0000-0000-0000-000000000000"
}

variable "resource_group_creation_enabled" {
  type    = bool
  default = false
}

variable "resource_groups" {
  type    = map(any)
  default = {}
}

variable "virtual_network_enabled" {
  type    = bool
  default = false
}

variable "virtual_networks" {
  type    = map(any)
  default = {}
}

variable "role_assignment_enabled" {
  type    = bool
  default = false
}

variable "role_assignments" {
  type    = map(any)
  default = {}
}

variable "route_table_enabled" {
  type    = bool
  default = false
}

variable "route_tables" {
  type    = map(any)
  default = {}
}

variable "enable_telemetry" {
  type    = bool
  default = true
}

output "subscription_resource_id" {
  value = try(module.lz_vending.subscription_resource_id, null)
}

output "virtual_network_resource_ids" {
  value = try(module.lz_vending.virtual_network_resource_ids, {})
}
