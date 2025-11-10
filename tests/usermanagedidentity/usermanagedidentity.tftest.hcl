# Tests for the usermanagedidentity module
# Converts the tests from tests/usermanagedidentity/usermanagedidentity_test.go

# Test 1: Basic user managed identity
run "basic_user_managed_identity" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
  }

  assert {
    condition     = var.name == "test"
    error_message = "UMI name should be test"
  }

  assert {
    condition     = var.location == "westeurope"
    error_message = "UMI location should be westeurope"
  }
}

# Test 2: UMI with GitHub federated credentials
run "umi_with_github_credentials" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    federated_credentials_github = {
      gh1 = {
        organization = "my-organization"
        repository   = "my-repository"
        entity       = "branch"
        value        = "my-branch"
      }
      gh2 = {
        organization = "my-organization2"
        repository   = "my-repository2"
        entity       = "pull_request"
      }
    }
  }

  assert {
    condition     = length(var.federated_credentials_github) == 2
    error_message = "Should have 2 GitHub federated credentials"
  }

  assert {
    condition     = var.federated_credentials_github["gh1"].entity == "branch"
    error_message = "First credential should be for branch"
  }

  assert {
    condition     = var.federated_credentials_github["gh2"].entity == "pull_request"
    error_message = "Second credential should be for pull_request"
  }
}

# Test 3: UMI with Terraform Cloud federated credentials
run "umi_with_terraform_cloud_credentials" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    federated_credentials_terraform_cloud = {
      tfc1 = {
        organization = "my-organization"
        project      = "my-repository"
        workspace    = "my-workspace"
        run_phase    = "plan"
      }
      tfc2 = {
        organization = "my-organization"
        project      = "my-repository"
        workspace    = "my-workspace"
        run_phase    = "apply"
      }
    }
  }

  assert {
    condition     = length(var.federated_credentials_terraform_cloud) == 2
    error_message = "Should have 2 Terraform Cloud federated credentials"
  }

  assert {
    condition     = var.federated_credentials_terraform_cloud["tfc1"].run_phase == "plan"
    error_message = "First credential should be for plan phase"
  }

  assert {
    condition     = var.federated_credentials_terraform_cloud["tfc2"].run_phase == "apply"
    error_message = "Second credential should be for apply phase"
  }
}

# Test 4: UMI with advanced federated credentials
run "umi_with_advanced_credentials" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    federated_credentials_advanced = {
      adv1 = {
        name               = "myadvancedcred1"
        subject_identifier = "field:value"
        issuer_url         = "https://test"
      }
      adv2 = {
        name               = "myadvancedcred2"
        subject_identifier = "field:value"
        issuer_url         = "https://test"
      }
    }
  }

  assert {
    condition     = length(var.federated_credentials_advanced) == 2
    error_message = "Should have 2 advanced federated credentials"
  }

  assert {
    condition     = var.federated_credentials_advanced["adv1"].issuer_url == "https://test"
    error_message = "Advanced credential should have correct issuer URL"
  }
}

# Test 5: Validation test - invalid Terraform Cloud run_phase
run "invalid_terraform_cloud_run_phase" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    federated_credentials_terraform_cloud = {
      tfc1 = {
        organization = "my-organization"
        project      = "my-repository"
        workspace    = "my-workspace"
        run_phase    = "check" # Invalid - must be 'plan' or 'apply'
      }
    }
  }

  expect_failures = [
    var.federated_credentials_terraform_cloud
  ]
}

# Test 6: Validation test - invalid GitHub credentials (missing value for branch)
run "invalid_github_credentials" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    federated_credentials_github = {
      gh1 = {
        organization = "my-organization"
        repository   = "my-repository"
        entity       = "branch"
        # Missing 'value' field - required for branch entity
      }
    }
  }

  expect_failures = [
    var.federated_credentials_github
  ]
}
