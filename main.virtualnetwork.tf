# Virtual networking submodule, disabled by default
# Will create a vnet, and optionally peerings and a virtual hub connection
module "virtualnetwork" {
  source = "./modules/virtualnetwork"
  count  = var.virtual_network_enabled ? 1 : 0

  # Required variables
  subscription_id                     = local.subscription_id
  virtual_network_address_space       = var.virtual_network_address_space
  virtual_network_location            = coalesce(var.location, var.virtual_network_location)
  virtual_network_name                = var.virtual_network_name
  virtual_network_resource_group_name = var.virtual_network_resource_group_name

  # Optional variables
  hub_network_resource_id                                  = var.hub_network_resource_id
  virtual_network_peering_enabled                          = var.virtual_network_peering_enabled
  virtual_network_use_remote_gateways                      = var.virtual_network_use_remote_gateways
  virtual_network_vwan_connection_enabled                  = var.virtual_network_vwan_connection_enabled
  virtual_network_vwan_propagated_routetables_labels       = var.virtual_network_vwan_propagated_routetables_labels
  virtual_network_vwan_propagated_routetables_resource_ids = var.virtual_network_vwan_propagated_routetables_resource_ids
  virtual_network_vwan_routetable_resource_id              = var.virtual_network_vwan_routetable_resource_id
  vwan_hub_resource_id                                     = var.vwan_hub_resource_id
}
