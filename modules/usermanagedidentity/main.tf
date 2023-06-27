resource "azapi_resource" "rg" {
  count     = var.resource_group_creation_enabled ? 1 : 0
  type      = "Microsoft.Resources/resourceGroups@2022-09-01"
  parent_id = "/subscriptions/${var.subscription_id}"
  name      = var.resource_group_name
  location  = var.location
  tags      = var.resource_group_tags
}

resource "azapi_resource" "rg_lock" {
  count     = var.resource_group_lock_enabled && var.resource_group_creation_enabled ? 1 : 0
  type      = "Microsoft.Authorization/locks@2020-05-01"
  parent_id = azapi_resource.rg[0].id
  name      = coalesce(var.resource_group_lock_name, "lock-${azapi_resource.rg[0].name}")
  body = jsonencode({
    properties = {
      level = "CanNotDelete"
    }
  })
  depends_on = [
    azapi_resource.rg,
    azapi_resource.umi,
    azapi_resource.umi_federated_credentials
  ]
}

resource "azapi_resource" "umi" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name      = var.name
  parent_id = var.resource_group_creation_enabled ? azapi_resource.rg[0].id : "/subscriptions/${var.subscription_id}/resourceGroups/${var.resource_group_name}"
  body      = jsonencode({})
  location  = var.location
  tags      = var.tags
  response_export_values = [
    "properties.principalId",
    "properties.clientId",
    "properties.tenantId"
  ]
}

resource "azapi_resource" "umi_federated_credentials" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = local.federated_credentials_map
  name      = each.value.name
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = [each.value.audience]
      issuer    = each.value.issuer_url
      subject   = each.value.subject_identifier
    }
  })
}
