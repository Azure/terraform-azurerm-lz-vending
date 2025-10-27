module "usermanagedidentity" {
  source   = "./modules/usermanagedidentity"
  for_each = { for umi_k, umi_v in var.user_managed_identities : umi_k => umi_v if var.umi_enabled }

  name     = each.value.name
  location = coalesce(each.value.location, var.location)
  tags     = each.value.tags
  parent_id = coalesce(
    can(module.resourcegroup[each.value.resource_group_key].resource_group_resource_id) ? module.resourcegroup[each.value.resource_group_key].resource_group_resource_id : null,
    each.value.resource_group_name_existing != null ? "${local.subscription_resource_id}/resourceGroups/${each.value.resource_group_name_existing}" : null
  )
  umi_skip_validation                   = local.umi_skip_validation
  federated_credentials_advanced        = each.value.federated_credentials_advanced
  federated_credentials_github          = each.value.federated_credentials_github
  federated_credentials_terraform_cloud = each.value.federated_credentials_terraform_cloud
}
