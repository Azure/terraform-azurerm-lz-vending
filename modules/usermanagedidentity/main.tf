resource "azapi_resource" "rg" {
  count = var.resource_group_creation_enabled ? 1 : 0

  type      = "Microsoft.Resources/resourceGroups@2022-09-01"
  location  = var.location
  name      = var.resource_group_name
  parent_id = "/subscriptions/${var.subscription_id}"
  tags      = var.resource_group_tags
}

resource "azapi_resource" "rg_lock" {
  count = var.resource_group_lock_enabled && var.resource_group_creation_enabled ? 1 : 0

  type = "Microsoft.Authorization/locks@2020-05-01"
  body = {
    properties = {
      level = "CanNotDelete"
    }
  }
  name      = coalesce(var.resource_group_lock_name, "lock-${one(azapi_resource.rg).name}")
  parent_id = one(azapi_resource.rg).id

  depends_on = [
    azapi_resource.rg,
    azapi_resource.umi,
    azapi_resource.umi_federated_credential_github_branch,
    azapi_resource.umi_federated_credential_github_tag,
    azapi_resource.umi_federated_credential_github_environment,
    azapi_resource.umi_federated_credential_github_pull_request,
    azapi_resource.umi_federated_credential_terraform_cloud,
    azapi_resource.umi_federated_credential_advanced,
  ]
}

resource "azapi_resource" "umi" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  body      = {}
  location  = var.location
  name      = var.name
  parent_id = var.resource_group_creation_enabled ? one(azapi_resource.rg).id : "/subscriptions/${var.subscription_id}/resourceGroups/${var.resource_group_name}"
  response_export_values = [
    "properties.principalId",
    "properties.clientId",
    "properties.tenantId"
  ]
  tags = var.tags
}

resource "azapi_resource" "umi_federated_credential_github_branch" {
  for_each = { for k, v in var.federated_credentials_github : k => v if v.entity == "branch" }

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = each.value.enterprise_slug != null ? "${local.github_actions_issuer}/${each.value.enterprise_slug}" : local.github_actions_issuer
      subject   = "repo:${each.value.organization}/${each.value.repository}:ref:refs/heads/${each.value.value}"
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-branch-${each.value.value}")
  parent_id = azapi_resource.umi.id
}

resource "azapi_resource" "umi_federated_credential_github_tag" {
  for_each = { for k, v in var.federated_credentials_github : k => v if v.entity == "tag" }

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = each.value.enterprise_slug != null ? "${local.github_actions_issuer}/${each.value.enterprise_slug}" : local.github_actions_issuer
      subject   = "repo:${each.value.organization}/${each.value.repository}:ref:refs/tags/${each.value.value}"
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-tag-${each.value.value}")
  parent_id = azapi_resource.umi.id
}

resource "azapi_resource" "umi_federated_credential_github_environment" {
  for_each = { for k, v in var.federated_credentials_github : k => v if v.entity == "environment" }

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = each.value.enterprise_slug != null ? "${local.github_actions_issuer}/${each.value.enterprise_slug}" : local.github_actions_issuer
      subject   = "repo:${each.value.organization}/${each.value.repository}:environment:${each.value.value}"
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-environment-${each.value.value}")
  parent_id = azapi_resource.umi.id
}

resource "azapi_resource" "umi_federated_credential_github_pull_request" {
  for_each = { for k, v in var.federated_credentials_github : k => v if v.entity == "pull_request" }

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = each.value.enterprise_slug != null ? "${local.github_actions_issuer}/${each.value.enterprise_slug}" : local.github_actions_issuer
      subject   = "repo:${each.value.organization}/${each.value.repository}:pull_request"
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = coalesce(each.value.name, "github-${each.value.organization}-${each.value.repository}-pull-request")
  parent_id = azapi_resource.umi.id
}

resource "azapi_resource" "umi_federated_credential_terraform_cloud" {
  for_each = var.federated_credentials_terraform_cloud

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = ["api://AzureADTokenExchange"]
      issuer    = "https://app.terraform.io"
      subject   = "organization:${each.value.organization}:project:${each.value.project}:workspace:${each.value.workspace}:run_phase:${each.value.run_phase}"
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = coalesce(each.value.name, "terraformcloud-${each.value.organization}-${each.value.project}-${each.value.workspace}-${each.value.run_phase}")
  parent_id = azapi_resource.umi.id
}

resource "azapi_resource" "umi_federated_credential_advanced" {
  for_each = var.federated_credentials_advanced

  type = "Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31"
  body = {
    properties = {
      audiences = each.value.audiences
      issuer    = each.value.issuer_url
      subject   = each.value.subject_identifier
    }
  }
  locks     = [azapi_resource.umi.id] # Concurrent Federated Identity Credentials writes under the same managed identity are not supported
  name      = each.value.name
  parent_id = azapi_resource.umi.id
}
