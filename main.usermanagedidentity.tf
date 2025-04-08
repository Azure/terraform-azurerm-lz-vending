module "usermanagedidentity" {
  source          = "./modules/usermanagedidentity"
  subscription_id = local.subscription_id

  for_each = { for umi_k, umi_v in var.user_managed_identities : umi_k => umi_v if var.umi_enabled }

  name     = each.value.name
  location = coalesce(each.value.location, var.location)
  tags     = each.value.tags

  resource_group_creation_enabled = var.resource_group_creation_enabled
  resource_group_name             = each.value.resource_group_name
  resource_group_lock_enabled     = each.value.resource_group_lock_enabled
  resource_group_lock_name        = each.value.resource_group_lock_name
  resource_group_tags             = each.value.resource_group_tags

  federated_credentials_advanced        = each.value.federated_credentials_advanced
  federated_credentials_github          = each.value.federated_credentials_github
  federated_credentials_terraform_cloud = each.value.federated_credentials_terraform_cloud

  depends_on = [
    module.resourcegroup
  ]
}
