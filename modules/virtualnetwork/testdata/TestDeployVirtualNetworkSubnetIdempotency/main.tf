resource "azapi_resource" "subnet" {
  parent_id = module.virtualnetwork_test.virtual_network_resource_id
  type      = "Microsoft.Network/virtualNetworks/subnets@2021-08-01"
  name      = "vnet-subnet"
  body = jsonencode({
    properties = {
      addressPrefix = "10.1.0.0/26"
    }
  })
}

module "virtualnetwork_test" {
  source                                = "../../"
  subscription_id                       = var.subscription_id
  virtual_network_address_space         = var.virtual_network_address_space
  virtual_network_location              = var.virtual_network_location
  virtual_network_resource_group_name   = var.virtual_network_resource_group_name
  virtual_network_name                  = var.virtual_network_name
  virtual_network_resource_lock_enabled = var.virtual_network_resource_lock_enabled
}
