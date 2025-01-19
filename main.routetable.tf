# route table submodule, disabled by default
# Will create a route table, and optionally routes
module "routetable" {
  source          = "./modules/routetable"
  count           = var.route_table_enabled ? 1 : 0
  subscription_id = local.subscription_id

  route_tables = var.route_tables
  depends_on = [
    module.resourcegroup,
  ]
}
