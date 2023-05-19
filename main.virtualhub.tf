# Virtual networking submodule, disabled by default
# Will create a vnet, and optionally peerings and a virtual hub connection
module "virtualhub" {
  source = "./modules/virtualhub"
  count  = var.intent_based_routing_enabled ? 1 : 0

  virtual_hubs = var.virtual_hubs
}
