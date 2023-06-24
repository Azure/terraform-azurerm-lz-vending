# process the federated credentials for GitHub
locals {
  federated_credentials_github_branch = [
    for v in var.federated_credentials_github : {
      subject_identifier = "repo:${v.organization}/${v.repository}:ref:refs/heads/${v.value}"
      name               = coalesce(v.name, "github-${v.organization}-${v.repository}-branch-${v.value}")
    }
    if v.entity == "branch"
  ]

  federated_credentials_github_tag = [
    for v in var.federated_credentials_github : {
      subject_identifier = "repo:${v.organization}/${v.repository}:ref:refs/tags/${v.value}"
      name               = coalesce(v.name, "github-${v.organization}-${v.repository}-tag-${v.value}")
    }
    if v.entity == "tag"
  ]

  federated_credentials_github_environment = [
    for v in var.federated_credentials_github : {
      subject_identifier = "repo:${v.organization}/${v.repository}:environment:${v.value}"
      name               = coalesce(v.name, "github-${v.organization}-${v.repository}-environment-${v.value}")
    }
    if v.entity == "environment"
  ]

  federated_credentials_github_pull_request = [
    for v in var.federated_credentials_github : {
      subject_identifier = "repo:${v.organization}/${v.repository}:pull_request"
      name               = coalesce(v.name, "github-${v.organization}-${v.repository}-pull_request")
    }
    if v.entity == "pull_request"
  ]

  # combine all the above into a single list
  federated_credentials_github = [
    for cred in concat(
      local.federated_credentials_github_branch,
      local.federated_credentials_github_tag,
      local.federated_credentials_github_environment,
      local.federated_credentials_github_pull_request,
      ) : {
      name               = cred.name
      subject_identifier = cred.subject_identifier
      audience           = "api://AzureADTokenExchange"
      issuer_url         = "https://token.actions.githubusercontent.com"
    }
  ]
}

# Process federated credentials for Terraform Cloud
locals {
  federated_credentials_terraform_cloud = [
    for v in var.federated_credentials_terraform_cloud : {
      name               = coalesce(v.name, "terraformcloud-${v.organization}-${v.project}-${v.workspace}-${v.run_phase}")
      subject_identifier = "organization:${v.organization}:project:${v.project}:workspace:${v.workspace}:run_phase:${v.run_phase}"
      audience           = "api://AzureADTokenExchange"
      issuer_url         = "https://app.terraform.io"
    }
  ]
}

# Combine all the federated credentials into a single set to use in the resource for_each
locals {
  federated_credentials_set = toset(concat(
    local.federated_credentials_github,
    local.federated_credentials_terraform_cloud,
    tolist(var.federated_credentials_advanced)
  ))
}

output "federated_credentials_set" {
  value = local.federated_credentials_set
}
