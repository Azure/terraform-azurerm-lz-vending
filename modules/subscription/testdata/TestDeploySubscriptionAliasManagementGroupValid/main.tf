terraform {
  required_providers {
    azapi = {
      source  = "azure/azapi"
      version = ">= 0.3.0"
    }
  }
}

variable "subscription_billing_scope" {
  type = string
}

variable "subscription_management_group_id" {
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

variable "subscription_management_group_association_enabled" {
  type = bool
}

variable "subscription_alias_enabled" {
  type = bool
}

resource "azapi_resource" "mg" {
  type      = "Microsoft.Management/managementGroups@2021-04-01"
  parent_id = "/"
  name      = var.subscription_management_group_id
}

module "subscription_test" {
  source                                            = "../../"
  subscription_alias_name                           = var.subscription_alias_name
  subscription_display_name                         = var.subscription_display_name
  subscription_workload                             = var.subscription_workload
  subscription_management_group_id                  = azapi_resource.mg.name
  subscription_billing_scope                        = var.subscription_billing_scope
  subscription_management_group_association_enabled = var.subscription_management_group_association_enabled
  subscription_alias_enabled                        = var.subscription_alias_enabled
}

output "subscription_id" {
  value = module.subscription_test.subscription_id
}
