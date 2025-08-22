# Telemetry is collected by creating an empty ARM deployment with a specific name
# If you want to disable telemetry, you can set the disable_telemetry variable to true
locals {
  # telem_arm_subscription_template is the ARM template content for the telemetry deployment.
  telem_arm_subscription_template = {
    "$schema"      = "https://schema.management.azure.com/schemas/2018-05-01/subscriptionDeploymentTemplate.json#"
    contentVersion = "1.0.0.0"
    parameters     = {}
    variables      = {}
    resources      = []
    outputs = {
      telemetry = {
        type  = "String"
        value = "For more information, see https://aka.ms/lz-vending/tf/telemetry"
      }
    }
  }
  telem_root_arm_deployment_name = substr(
    format(
      "pid-%s_%s_%s",
      local.telem_root_puid,
      local.module_version,
      local.telem_root_bitfield_hex,
    ),
    0,
    64
  )
  # Calculate the denary value of the bit fields
  telem_root_bitfield_denary = (
    local.telem_root_subscription_alias_enabled +
    local.telem_root_subscription_management_group_association_enabled +
    local.telem_root_subscription_tags_enabled +
    local.telem_root_virtual_network_enabled +
    local.telem_root_virtual_network_peering_enabled +
    local.telem_root_virtual_network_vwan_connection_enabled +
    local.telem_root_vwan_advanced_routing_enabled +
    local.telem_root_role_assignment_enabled
  )
  # Convert the denary value to hexadecimal and pad with zeros to the left to a length of 8 characters.
  telem_root_bitfield_hex = format("%08x", local.telem_root_bitfield_denary)
  # root_module_puid is the UUID that identifies the root module in the telemetry ARM deployment.
  telem_root_puid = "50a8a460-d517-4b11-b86c-6de447806b67"
  # roleassignment telemetry bit fields
  telem_root_role_assignment_enabled = var.role_assignment_enabled ? 65536 : 0
  # subscription telemetry bit fields
  telem_root_subscription_alias_enabled                        = var.subscription_alias_enabled ? 1 : 0
  telem_root_subscription_management_group_association_enabled = var.subscription_management_group_association_enabled ? 2 : 0
  telem_root_subscription_tags_enabled                         = length(var.subscription_tags) > 0 ? 4 : 0
  # virtualnetwork telemetry bit fields
  telem_root_virtual_network_enabled                 = var.virtual_network_enabled ? 256 : 0
  telem_root_virtual_network_peering_enabled         = var.virtual_network_enabled && anytrue([for k, v in var.virtual_networks : v.hub_peering_enabled]) ? 512 : 0
  telem_root_virtual_network_vwan_connection_enabled = var.virtual_network_enabled && anytrue([for k, v in var.virtual_networks : v.vwan_connection_enabled]) ? 1024 : 0
  telem_root_vwan_advanced_routing_enabled           = var.virtual_network_enabled && anytrue([for k, v in var.virtual_networks : length(v.vwan_propagated_routetables_labels) > 0 || length(v.vwan_propagated_routetables_resource_ids) > 0 || v.vwan_associated_routetable_resource_id != null]) ? 4096 : 0
}
