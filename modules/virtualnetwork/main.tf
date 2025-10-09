# module.virtual_networks uses the Azure Verified Module to create
# as many virtual networks as is required by the var.virtual_networks input variable
module "virtual_networks" {
  for_each        = var.virtual_networks
  source          = "Azure/avm-res-network-virtualnetwork/azurerm"
  version         = "0.14.1" # AVM version with IPAM support
  subscription_id = var.subscription_id

  name                    = each.value.name
  address_space           = lookup(each.value, "address_space", null)
  resource_group_name     = lookup(each.value, "resource_group_name", null)
  location                = coalesce(lookup(each.value, "location", null), var.location)
  flow_timeout_in_minutes = lookup(each.value, "flow_timeout_in_minutes", null)

  # DDoS protection plan configuration
  ddos_protection_plan = lookup(each.value, "ddos_protection_plan_id", null) == null ? null : {
    id     = lookup(each.value, "ddos_protection_plan_id", null)
    enable = true
  }

  # DNS servers configuration
  dns_servers = length(lookup(each.value, "dns_servers", [])) == 0 ? null : {
    dns_servers = lookup(each.value, "dns_servers", [])
  }

  # Keep existing static subnet map (existing behaviour)
  subnets = lookup(each.value, "subnets", null)

  # -------------------------
  # IPAM-related inputs (now supported in v0.14.1)
  # -------------------------
  # Toggle IPAM allocation for this VNet (false = keep using static address_space)
  enable_ipam = lookup(each.value, "enable_ipam", false)

  # Full resource id of the Azure Virtual Network Manager (Network Manager)
  ipam_network_manager_id = lookup(each.value, "ipam_network_manager_id", null)

  # Full resource id of the IPAM pool used to allocate the VNet address space
  ipam_vnet_pool_id = lookup(each.value, "ipam_vnet_pool_id", null)

  # Per-subnet IPAM allocations — pass the AVM-shaped list/map (or null)
  ipam_subnet_allocations = lookup(each.value, "ipam_subnet_allocations", null)

  # If you want to attach to an existing VNet instead of creating a new one:
  existing_vnet_id = lookup(each.value, "existing_vnet_id", null)

  # Providers: AVM IPAM requires azapi for certain resources — ensure azapi provider is configured in your repo
  providers = {
    azurerm = azurerm
    azapi   = azapi
  }

  tags             = lookup(each.value, "tags", null)
  enable_telemetry = var.enable_telemetry
}

# module.peering_hub_outbound uses the peering submodule from theAzure Verified Module
# to create the outboud peering from the spoke to the hub network when specified
module "peering_hub_outbound" {
  for_each        = { for k, v in local.hub_peering_map : k => v if v.peering_direction != local.peering_direction_fromhub }
  source          = "Azure/avm-res-network-virtualnetwork/azurerm//modules/peering"
  version         = "0.8.1"
  subscription_id = var.subscription_id

  virtual_network = {
    "resource_id" = each.value["outbound"].this_resource_id,
  }
  remote_virtual_network = {
    "resource_id" = each.value["outbound"].remote_resource_id,
  }
  name                         = each.value.outbound.name
  allow_forwarded_traffic      = each.value.outbound.options.allow_forwarded_traffic
  allow_gateway_transit        = each.value.outbound.options.allow_gateway_transit
  allow_virtual_network_access = each.value.outbound.options.allow_virtual_network_access
  use_remote_gateways          = each.value.outbound.options.use_remote_gateways
  create_reverse_peering       = false

  depends_on = [module.virtual_networks]
}

# module.peering_hub_inbound uses the peering submodule from theAzure Verified Module
# to create the inbound peering from the hub network to the spoke network when specified
module "peering_hub_inbound" {
  for_each        = { for k, v in local.hub_peering_map : k => v if v.peering_direction != local.peering_direction_tohub }
  source          = "Azure/avm-res-network-virtualnetwork/azurerm//modules/peering"
  version         = "0.8.1"
  subscription_id = var.subscription_id

  virtual_network = {
    "resource_id" = each.value["inbound"].this_resource_id,
  }
  remote_virtual_network = {
    "resource_id" = each.value["inbound"].remote_resource_id,
  }
  name                         = each.value.inbound.name
  allow_forwarded_traffic      = each.value.inbound.options.allow_forwarded_traffic
  allow_gateway_transit        = each.value.inbound.options.allow_gateway_transit
  allow_virtual_network_access = each.value.inbound.options.allow_virtual_network_access
  use_remote_gateways          = each.value.inbound.options.use_remote_gateways
  create_reverse_peering       = false

  depends_on = [module.virtual_networks]
}

# module.peering_mesh uses the peering submodule from theAzure Verified Module
# to create the peering from the local and remote virtual networks as specified
module "peering_mesh" {
  for_each        = { for i in local.virtual_networks_mesh_peering_list : "${i.source_key}-${i.destination_key}" => i }
  source          = "Azure/avm-res-network-virtualnetwork/azurerm//modules/peering"
  version         = "0.8.1"
  subscription_id = var.subscription_id

  virtual_network = {
    "resource_id" = each.value.this_resource_id,
  }
  remote_virtual_network = {
    "resource_id" = each.value.remote_resource_id,
  }
  name                         = each.value.name
  allow_forwarded_traffic      = each.value.allow_forwarded_traffic
  allow_gateway_transit        = false
  allow_virtual_network_access = true
  use_remote_gateways          = false
  create_reverse_peering       = false

  depends_on = [module.virtual_networks]
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
resource "azapi_resource" "vhubconnection" {
  for_each = { for k, v in var.virtual_networks : k => v if v.vwan_connection_enabled && !v.vwan_security_configuration.routing_intent_enabled }

  type = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2022-07-01"
  body = {
    properties = local.vhubconnection_body_properties[each.key]
  }
  name      = coalesce(each.value.vwan_connection_name, "vhc-${uuidv5("url", module.virtual_networks[each.key].resource_id)}")
  parent_id = each.value.vwan_hub_resource_id

  depends_on = [module.virtual_networks]
}

# azapi_resource.vhubconnection creates a virtual wan hub connection between the spoke and the supplied vwan hub.
# This resource is used when routing intent is enabled on the vwan security configuration,
# as the routing configuration is then ignored.
resource "azapi_resource" "vhubconnection_routing_intent" {
  for_each = { for k, v in var.virtual_networks : k => v if v.vwan_connection_enabled && v.vwan_security_configuration.routing_intent_enabled }

  type = "Microsoft.Network/virtualHubs/hubVirtualNetworkConnections@2022-07-01"
  body = {
    properties = local.vhubconnection_body_properties[each.key]
  }
  name      = coalesce(each.value.vwan_connection_name, "vhc-${uuidv5("url", module.virtual_networks[each.key].resource_id)}")
  parent_id = each.value.vwan_hub_resource_id

  depends_on = [module.virtual_networks]

  lifecycle {
    ignore_changes = [
      body.properties.routingConfiguration,
    ]
  }
}
