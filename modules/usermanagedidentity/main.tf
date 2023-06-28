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
    azapi_resource.umi_federated_credentials_github_branch,
    azapi_resource.umi_federated_credentials_github_tag,
    azapi_resource.umi_federated_credentials_github_environment,
    azapi_resource.umi_federated_credentials_github_pull_request,
    azapi_resource.umi_federated_credentials_terraform_cloud,
    azapi_resource.umi_federated_credentials_advanced,
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

resource "azapi_resource" "umi_federated_credentials_github_branch" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = { for k, v in var.federated_credentials_github : k => v if v.entity == "branch" }
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-branch-${each.value.value}")
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://token.actions.githubusercontent.com"
      subject   = "repo:${each.value.organization}/${each.value.repository}:ref:refs/heads/${each.value.value}"
    }
  })
}

resource "azapi_resource" "umi_federated_credentials_github_tag" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = { for k, v in var.federated_credentials_github : k => v if v.entity == "tag" }
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-tag-${each.value.value}")
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://token.actions.githubusercontent.com"
      subject   = "repo:${each.value.organization}/${each.value.repository}:ref:refs/tags/${each.value.value}"
    }
  })
}

resource "azapi_resource" "umi_federated_credentials_github_environment" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = { for k, v in var.federated_credentials_github : k => v if v.entity == "environment" }
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-environment-${each.value.value}")
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://token.actions.githubusercontent.com"
      subject   = "repo:${each.value.organization}/${each.value.repository}:environment:${each.value.value}"
    }
  })
}

resource "azapi_resource" "umi_federated_credentials_github_pull_request" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = { for k, v in var.federated_credentials_github : k => v if v.entity == "pull_request" }
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-pull-request")
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://token.actions.githubusercontent.com"
      subject   = "repo:${each.value.organization}/${each.value.repository}:pull_request"
    }
  })
}

resource "azapi_resource" "umi_federated_credentials_terraform_cloud" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = var.federated_credentials_terraform_cloud
  name      = coalesce(each.value.name, "terraformcloud-${each.value.organization}-${each.value.project}-${each.value.workspace}-${each.value.run_phase}")
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://app.terraform.io"
      subject   = "organization:${each.value.organization}:project:${each.value.project}:workspace:${each.value.workspace}:run_phase:${each.value.run_phase}"
    }
  })
}

resource "azapi_resource" "umi_federated_credentials_advanced" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  for_each  = var.federated_credentials_advanced
  name      = each.value.name
  parent_id = azapi_resource.umi.id
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  body = jsonencode({
    properties = {
      audiences = each.value.audiences
      issuer    = each.value.issuer_url
      subject   = each.value.subject_identifier
    }
  })
}
