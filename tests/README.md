# Terraform Native Tests

This directory contains Terraform native tests that replace the previous Go/Terratest-based tests. These tests use Terraform's built-in `terraform test` command introduced in Terraform 1.6+.

## Directory Structure

```
terraform-tests/
├── virtualnetwork/          # Tests for the virtualnetwork module
│   ├── main.tf             # Test wrapper configuration
│   ├── providers.tf        # Provider configuration
│   ├── virtualnetwork_basic.tftest.hcl    # Basic plan tests (7 scenarios)
│   └── virtualnetwork_deploy.tftest.hcl   # Deployment tests (3 scenarios)
├── integration/            # Integration tests for the root module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   └── integration.tftest.hcl    # Integration test scenarios (6 scenarios)
├── subscription/           # Tests for the subscription module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   ├── variables.tf       # Variable declarations
│   └── subscription.tftest.hcl    # Subscription tests (8 scenarios: 2 positive, 6 validation)
├── resourcegroup/          # Tests for the resourcegroup module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   ├── variables.tf       # Variable declarations
│   └── resourcegroup.tftest.hcl   # Resource group tests (1 scenario)
├── budget/                 # Tests for the budget module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   ├── variables.tf       # Variable declarations
│   └── budget.tftest.hcl  # Budget tests (1 scenario)
├── networksecuritygroup/   # Tests for the networksecuritygroup module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   ├── variables.tf       # Variable declarations
│   └── networksecuritygroup.tftest.hcl    # NSG tests (8 scenarios)
├── resourceprovider/       # Tests for the resourceprovider module
│   ├── main.tf            # Test wrapper configuration
│   ├── providers.tf       # Provider configuration
│   ├── variables.tf       # Variable declarations
│   └── resourceprovider.tftest.hcl    # Resource provider tests (1 scenario)
└── usermanagedidentity/    # Tests for the usermanagedidentity module
    ├── main.tf            # Test wrapper configuration
    ├── providers.tf       # Provider configuration
    ├── variables.tf       # Variable declarations
    └── usermanagedidentity.tftest.hcl    # UMI tests (6 scenarios: 4 positive, 2 validation)
```

**Note**: Role assignment tests already exist in `modules/roleassignment/tests/unit/roleassignments.tftest.hcl` with mock providers.

## Running Tests

### Prerequisites

1. **Terraform 1.9+**: These tests require Terraform 1.9.0 or later
   ```powershell
   terraform version
   ```

2. **Azure Authentication**: For deployment tests, configure Azure authentication:
   ```powershell
   az login
   # Or use service principal
   $env:ARM_CLIENT_ID = "your-client-id"
   $env:ARM_CLIENT_SECRET = "your-client-secret"
   $env:ARM_TENANT_ID = "your-tenant-id"
   $env:ARM_SUBSCRIPTION_ID = "your-subscription-id"
   ```

### Running Plan Tests (No Azure Required)

Plan tests only validate the Terraform configuration without deploying resources:

```powershell
# Test virtual network module
cd terraform-tests/virtualnetwork
terraform init
terraform test -filter=tests/virtualnetwork_basic.tftest.hcl

# Test integration scenarios
cd ../integration
terraform init
terraform test -filter=tests/integration.tftest.hcl

# Test subscription module (includes validation tests)
cd ../subscription
terraform init
terraform test

# Test other modules
cd ../resourcegroup && terraform init && terraform test
cd ../budget && terraform init && terraform test
cd ../networksecuritygroup && terraform init && terraform test
cd ../resourceprovider && terraform init && terraform test
cd ../usermanagedidentity && terraform init && terraform test
```

### Running Deployment Tests (Requires Azure)

Deployment tests actually create resources in Azure. **These incur costs and require cleanup.**

```powershell
# Set environment variables
$env:AZURE_SUBSCRIPTION_ID = "your-subscription-id"
$env:AZURE_TENANT_ID = "your-tenant-id"

# Run deployment tests for virtualnetwork
cd terraform-tests/virtualnetwork
terraform init
terraform test -filter=tests/virtualnetwork_deploy.tftest.hcl
```

**Note**: Deployment tests assume resource groups already exist. You may need to create them first or modify the tests to create them.

### Running All Tests

```powershell
# Run all tests in a directory
cd terraform-tests/virtualnetwork
terraform test

# Or run specific test files
terraform test tests/virtualnetwork_basic.tftest.hcl
```

## Test Types

### 1. Plan Tests (`command = plan`)
- Validate Terraform configuration without deploying
- Check resource counts, attributes, and relationships
- Fast and don't require Azure credentials
- Safe to run in CI/CD pipelines

### 2. Apply Tests (`command = apply`)
- Actually deploy resources to Azure
- Validate real-world behavior
- **Require valid Azure credentials and subscription**
- **Incur Azure costs**
- Should be run selectively

## Migrating from Terratest

### Key Differences

| Aspect | Terratest (Go) | Terraform Native |
|--------|---------------|------------------|
| Language | Go | HCL (Terraform) |
| Test Runner | `go test` | `terraform test` |
| Setup | Requires Go toolchain | Only Terraform CLI |
| Assertions | Go test assertions | HCL `assert` blocks |
| Resource Inspection | `terraform.Output()` | Direct module output access |
| Plan Inspection | JSON parsing in Go | Native plan validation |

### Migration Notes

1. **Test Structure**:
   - Go test functions → `run` blocks in `.tftest.hcl` files
   - Each `run` block is an independent test scenario

2. **Variables**:
   - Go variable maps → `variables` block in run blocks
   - Type-safe variable declarations in wrapper `main.tf`

3. **Assertions**:
   ```go
   // Old (Terratest/Go)
   assert.Equal(t, 2, len(vnets))
   ```
   ```hcl
   # New (Terraform native)
   assert {
     condition     = length(keys(module.virtualnetwork.virtual_network_resource_ids)) == 2
     error_message = "Expected 2 virtual networks"
   }
   ```

4. **Module Testing**:
   - Create a wrapper `main.tf` that calls the module under test
   - Use `module.<name>.<output>` to access module outputs
   - Pass variables through to the module

5. **Deployment Tests**:
   - Use `command = apply` instead of `command = plan`
   - Terraform handles resource creation and cleanup
   - Use `depends_on` or multiple run blocks for setup/teardown

## Best Practices

1. **Start with Plan Tests**: Write plan tests first to validate configuration logic
2. **Minimize Apply Tests**: Only test deployment where necessary (e.g., idempotency, real API validation)
3. **Use Unique Names**: Generate unique resource names to avoid conflicts
4. **Clean Up**: Terraform test automatically destroys resources after `apply` tests
5. **Organize Tests**: Group related tests in the same `.tftest.hcl` file
6. **Document Assumptions**: Note any prerequisites (e.g., existing resource groups)

## Examples

### Basic Plan Test
```hcl
run "valid_vnets" {
  command = plan

  variables {
    subscription_id  = "00000000-0000-0000-0000-000000000000"
    virtual_networks = {
      primary = {
        name                = "test-vnet"
        address_space       = ["10.0.0.0/16"]
        location            = "eastus"
        resource_group_name = "test-rg"
      }
    }
  }

  assert {
    condition     = length(keys(module.virtualnetwork.virtual_network_resource_ids)) == 1
    error_message = "Expected 1 virtual network"
  }
}
```

### Deployment Test with Setup
```hcl
run "setup_resource_group" {
  command = apply

  module {
    source = "./setup"
  }
}

run "deploy_vnet" {
  command = apply

  variables {
    resource_group_name = run.setup_resource_group.resource_group_name
  }
}
```

## Troubleshooting

### Tests Fail to Initialize
```powershell
# Clean and reinitialize
Remove-Item -Recurse -Force .terraform
terraform init
```

### Provider Configuration Issues
- Ensure `providers.tf` exists in the test directory
- Check provider version constraints match module requirements

### Module Not Found
- Verify relative paths in `source = "../../module-path"`
- Ensure you're running from the correct directory

### Azure Authentication Failures
```powershell
# Verify authentication
az account show
# Or check environment variables
$env:ARM_SUBSCRIPTION_ID
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Terraform Tests

on: [pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: "1.9.0"

      - name: Run Plan Tests
        run: |
          cd terraform-tests/virtualnetwork
          terraform init
          terraform test -filter=tests/virtualnetwork_basic.tftest.hcl
```

## Further Reading

- [Terraform Test Documentation](https://developer.hashicorp.com/terraform/language/tests)
- [Azure Verified Modules Testing](https://azure.github.io/Azure-Verified-Modules/contributing/terraform/testing/)
- [Terraform Testing Best Practices](https://developer.hashicorp.com/terraform/tutorials/configuration-language/test)
