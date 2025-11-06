#!/usr/bin/env pwsh

[CmdletBinding()]
param(
  [Parameter(Position = 0, Mandatory = $true)]
  [string]$Target
)

Set-StrictMode -Version 3.0
$ErrorActionPreference = "Stop"

function Show-Usage {
  Write-Host "Usage: avm <make target>"
}

# Default values for environment variables
$CONTAINER_RUNTIME = if ($env:CONTAINER_RUNTIME) { $env:CONTAINER_RUNTIME } else { "docker" }
$CONTAINER_IMAGE = if ($env:CONTAINER_IMAGE) { $env:CONTAINER_IMAGE } else { "mcr.microsoft.com/azterraform:avm-latest" }
$CONTAINER_PULL_POLICY = if ($env:CONTAINER_PULL_POLICY) { $env:CONTAINER_PULL_POLICY } else { "always" }
$MAKEFILE_REF = if ($env:MAKEFILE_REF) { $env:MAKEFILE_REF } else { "main" }
$PORCH_REF = if ($env:PORCH_REF) { $env:PORCH_REF } else { "main" }

# Check if container runtime is available
if (-not (Get-Command $CONTAINER_RUNTIME -ErrorAction SilentlyContinue) -and -not $env:AVM_IN_CONTAINER) {
  Write-Error "Error: $CONTAINER_RUNTIME is not installed. Please install $CONTAINER_RUNTIME first."
  exit 1
}

# Check if AZURE_CONFIG_DIR is set, if not, set it to ~/.azure
$AZURE_CONFIG_DIR = if ($env:AZURE_CONFIG_DIR) {
  $env:AZURE_CONFIG_DIR
}
else {
  if ($IsWindows) {
    Join-Path $env:USERPROFILE ".azure"
  }
  else {
    Join-Path $env:HOME ".azure"
  }
}

# Check if AZURE_CONFIG_DIR exists, if it does, mount it to the container
$AZURE_CONFIG_MOUNT = $null
$AZURE_CONFIG_MOUNT_PATH = $null
if (Test-Path $AZURE_CONFIG_DIR) {
  $AZURE_CONFIG_MOUNT = "-v"
  $AZURE_CONFIG_MOUNT_PATH = "${AZURE_CONFIG_DIR}:/home/runtimeuser/.azure"
}

# New: allow overriding TUI behavior with PORCH_FORCE_TUI and PORCH_NO_TUI environment variables.
# - If PORCH_FORCE_TUI is set, force TUI and interactive mode (even in GH Actions).
# - If PORCH_NO_TUI is set, explicitly disable TUI.
# - Otherwise, fallback to previous behavior: enable TUI only when not in GitHub Actions and NO_COLOR is not set.
$TUI = $null
$DOCKER_INTERACTIVE = $null
if ($env:PORCH_FORCE_TUI -and $env:PORCH_FORCE_TUI -ne "") {
  $TUI = "--tui"
  $DOCKER_INTERACTIVE = "-it"
  $env:FORCE_COLOR = "1"
}
elseif ($env:PORCH_NO_TUI -and $env:PORCH_NO_TUI -ne "") {
  # Explicitly disable TUI and interactive flags
  $TUI = $null
  $DOCKER_INTERACTIVE = $null
}
else {
  # If we are not in GitHub Actions and NO_COLOR is not set, we want to use TUI and interactive mode
  if (-not $env:GITHUB_RUN_ID -and -not $env:NO_COLOR) {
    $TUI = "--tui"
    $DOCKER_INTERACTIVE = "-it"
    $env:FORCE_COLOR = "1"
  }
}

# if PORCH_BASE_URL is set, we want to add it to the make command
$PORCH_BASE_URL_MAKE_ADD = $null
if ($env:PORCH_BASE_URL) {
  $PORCH_BASE_URL_MAKE_ADD = "PORCH_BASE_URL=$($env:PORCH_BASE_URL)"
}

# Check if we are running in a container
# If we are then just run make directly
if (-not $env:AVM_IN_CONTAINER) {
  # Build the docker command arguments
  $dockerArgs = @(
    "run"
    "--pull", $CONTAINER_PULL_POLICY
    "--rm"
  )

  # Add user parameter only on Unix-like systems
  if (-not $IsWindows) {
    try {
      $userId = & id -u
      $groupId = & id -g
      $dockerArgs += @("--user", "${userId}:${groupId}")
    }
    catch {
      Write-Warning "Could not determine user/group ID, running without --user parameter"
    }
  }

  if ($DOCKER_INTERACTIVE) {
    $dockerArgs += $DOCKER_INTERACTIVE
  }

  $dockerArgs += @(
    "-v", "$(Get-Location):/src"
  )

  if ($AZURE_CONFIG_MOUNT -and $AZURE_CONFIG_MOUNT_PATH) {
    $dockerArgs += @($AZURE_CONFIG_MOUNT, $AZURE_CONFIG_MOUNT_PATH)
  }

  # Add environment variables
  $envVars = @(
    "ARM_CLIENT_ID",
    "ARM_OIDC_REQUEST_TOKEN",
    "ARM_OIDC_REQUEST_URL",
    "ARM_SUBSCRIPTION_ID",
    "ARM_TENANT_ID",
    "ARM_USE_OIDC",
    "AVM_EXAMPLE",
    "CONFTEST_APRL_URL",
    "CONFTEST_AVMSEC_URL",
    "CONFTEST_EXCEPTIONS_URL",
    "FORCE_COLOR",
    "GITHUB_TOKEN",
    "GREPT_URL",
    "MPTF_URL",
    "NO_COLOR",
    "PORCH_LOG_LEVEL",
    "TEST_TYPE",
    "TFLINT_CONFIG_URL"
  )

  foreach ($envVar in $envVars) {
    $envValue = [System.Environment]::GetEnvironmentVariable($envVar)
    if ($null -ne $envValue -and $envValue -ne "") {
      $dockerArgs += @("-e", $envVar)
    }
  }

  # Add TF_IN_AUTOMATION
  $dockerArgs += @("-e", "TF_IN_AUTOMATION=1")

  # Add TF_VAR_ environment variables
  Get-ChildItem env: | Where-Object { $_.Name -like "TF_VAR_*" } | ForEach-Object {
    $dockerArgs += @("-e", "$($_.Name)=$($_.Value)")
  }

  # Add AVM_ environment variables
  Get-ChildItem env: | Where-Object { $_.Name -like "AVM_*" } | ForEach-Object {
    $dockerArgs += @("-e", "$($_.Name)=$($_.Value)")
  }

  $dockerArgs += $CONTAINER_IMAGE
  $dockerArgs += "make"

  if ($TUI) {
    $dockerArgs += "TUI=$TUI"
  }

  $dockerArgs += "MAKEFILE_REF=$MAKEFILE_REF"

  if ($PORCH_BASE_URL_MAKE_ADD) {
    $dockerArgs += $PORCH_BASE_URL_MAKE_ADD
  }

  $dockerArgs += "PORCH_REF=$PORCH_REF"
  $dockerArgs += $Target

  & $CONTAINER_RUNTIME @dockerArgs
}
else {
  # Build make command arguments
  $makeArgs = @()

  if ($TUI) {
    $makeArgs += "TUI=$TUI"
  }

  $makeArgs += "MAKEFILE_REF=$MAKEFILE_REF"

  if ($PORCH_BASE_URL_MAKE_ADD) {
    $makeArgs += $PORCH_BASE_URL_MAKE_ADD
  }

  $makeArgs += "PORCH_REF=$PORCH_REF"
  $makeArgs += $Target

  & make @makeArgs
}
