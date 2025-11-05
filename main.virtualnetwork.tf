# Virtual networking submodule, disabled by default
# Will create a vnet, and optionally peerings and a virtual hub connection
module "virtualnetwork" {
  source                  = "./modules/virtualnetwork"
  count                   = var.virtual_network_enabled ? 1 : 0
  subscription_id         = local.subscription_id
  virtual_networks        = local.virtual_networks
  location                = var.location
  enable_telemetry        = !var.disable_telemetry
  ipam_pool_id_by_vnet    = var.ipam_pool_id_by_vnet
  ipam_pool_prefix_length = var.ipam_pool_prefix_length

  depends_on = [
    module.resourcegroup,
    module.routetable,
    module.networksecuritygroup
  ]
}
