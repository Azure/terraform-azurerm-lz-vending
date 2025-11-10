# Test wrapper for the networksecuritygroup module
module "networksecuritygroup" {
  source = "../../modules/networksecuritygroup"

  name           = var.name
  location       = var.location
  parent_id      = var.parent_id
  tags           = var.tags
  security_rules = var.security_rules
}
