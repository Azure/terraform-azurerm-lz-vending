#!/usr/bin/env pwsh
# Quick test runner script for Terraform tests

param(
    [Parameter()]
    [ValidateSet('plan', 'deploy', 'integration', 'modules', 'all')]
    [string]$TestType = 'plan',

    [Parameter()]
    [string]$SubscriptionId = $env:AZURE_SUBSCRIPTION_ID
)

$ErrorActionPreference = "Stop"

function Write-TestHeader {
    param([string]$Message)
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host $Message -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
}

function Run-TerraformTest {
    param(
        [string]$Path,
        [string]$Filter = ""
    )

    Push-Location $Path
    try {
        Write-Host "Initializing Terraform..." -ForegroundColor Yellow
        terraform init -upgrade

        if ($Filter) {
            Write-Host "Running tests: $Filter" -ForegroundColor Yellow
            terraform test -filter=$Filter
        } else {
            Write-Host "Running all tests in $Path" -ForegroundColor Yellow
            terraform test
        }

        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Tests passed" -ForegroundColor Green
        } else {
            Write-Host "✗ Tests failed" -ForegroundColor Red
            exit 1
        }
    } finally {
        Pop-Location
    }
}

# Validate subscription ID for deployment tests
if ($TestType -in @('deploy', 'all') -and -not $SubscriptionId) {
    Write-Error "AZURE_SUBSCRIPTION_ID environment variable must be set for deployment tests"
    exit 1
}

# Set subscription ID if provided
if ($SubscriptionId) {
    $env:AZURE_SUBSCRIPTION_ID = $SubscriptionId
    Write-Host "Using subscription: $SubscriptionId" -ForegroundColor Green
}

# Get the terraform-tests directory
$TestsRoot = $PSScriptRoot

switch ($TestType) {
    'plan' {
        Write-TestHeader "Running Plan Tests (VirtualNetwork)"
        Run-TerraformTest -Path "$TestsRoot/virtualnetwork" -Filter "tests/virtualnetwork_basic.tftest.hcl"

        Write-TestHeader "Running Integration Tests"
        Run-TerraformTest -Path "$TestsRoot/integration" -Filter "tests/integration.tftest.hcl"
    }

    'deploy' {
        Write-TestHeader "Running Deployment Tests (VirtualNetwork)"
        Write-Host "WARNING: This will create real Azure resources!" -ForegroundColor Yellow
        Write-Host "Press Ctrl+C within 5 seconds to cancel..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5

        Run-TerraformTest -Path "$TestsRoot/virtualnetwork" -Filter "tests/virtualnetwork_deploy.tftest.hcl"
    }

    'integration' {
        Write-TestHeader "Running Integration Tests Only"
        Run-TerraformTest -Path "$TestsRoot/integration" -Filter "tests/integration.tftest.hcl"
    }

    'modules' {
        Write-TestHeader "Running Module Tests"

        Write-Host "`nTesting subscription module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/subscription"

        Write-Host "`nTesting resourcegroup module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/resourcegroup"

        Write-Host "`nTesting budget module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/budget"

        Write-Host "`nTesting networksecuritygroup module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/networksecuritygroup"

        Write-Host "`nTesting resourceprovider module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/resourceprovider"

        Write-Host "`nTesting usermanagedidentity module..." -ForegroundColor Cyan
        Run-TerraformTest -Path "$TestsRoot/usermanagedidentity"
    }

    'all' {
        Write-TestHeader "Running ALL Tests"
        Write-Host "WARNING: Deployment tests will create real Azure resources!" -ForegroundColor Yellow
        Write-Host "Press Ctrl+C within 5 seconds to cancel..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5

        Write-TestHeader "1/9: VirtualNetwork Plan Tests"
        Run-TerraformTest -Path "$TestsRoot/virtualnetwork" -Filter "tests/virtualnetwork_basic.tftest.hcl"

        Write-TestHeader "2/9: Integration Tests"
        Run-TerraformTest -Path "$TestsRoot/integration" -Filter "tests/integration.tftest.hcl"

        Write-TestHeader "3/9: Subscription Module Tests"
        Run-TerraformTest -Path "$TestsRoot/subscription"

        Write-TestHeader "4/9: Resource Group Module Tests"
        Run-TerraformTest -Path "$TestsRoot/resourcegroup"

        Write-TestHeader "5/9: Budget Module Tests"
        Run-TerraformTest -Path "$TestsRoot/budget"

        Write-TestHeader "6/9: Network Security Group Module Tests"
        Run-TerraformTest -Path "$TestsRoot/networksecuritygroup"

        Write-TestHeader "7/9: Resource Provider Module Tests"
        Run-TerraformTest -Path "$TestsRoot/resourceprovider"

        Write-TestHeader "8/9: User Managed Identity Module Tests"
        Run-TerraformTest -Path "$TestsRoot/usermanagedidentity"

        Write-TestHeader "9/9: Deployment Tests"
        Run-TerraformTest -Path "$TestsRoot/virtualnetwork" -Filter "tests/virtualnetwork_deploy.tftest.hcl"
    }
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "All tests completed successfully!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green
