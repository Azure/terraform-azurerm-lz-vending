package integration

// Integration tests apply to the root module and test the deployment of resources
// in specific scenarios.

import (
	"testing"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../"
)

// TestIntegrationHubAndSpoke tests the resource plan when creating a new subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpoke(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	terraformOptions.Vars = v

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 5, "expected 5 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.subscription[0].azapi_resource.subscription_alias",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
	}
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

// TestIntegrationVwan tests the resource plan when creating a new subscription,
// with a new virtual network and vwan connection to a supplied vhub.
func TestIntegrationVwan(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["vwan_hub_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualHubs/testhub"
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = true
	v["virtual_network_vwan_connection_enabled"] = true
	terraformOptions.Vars = v

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 4, "expected 4 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.subscription[0].azapi_resource.subscription_alias",
		"module.virtualnetwork[0].azapi_resource.vhubconnection[\"vhubcon-1b4db7eb-4057-5ddf-91e0-36dec72071f5\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
	}
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

// TestIntegrationHubAndSpokeExistingSubscription tests the resource plan when supplying an existing subscription,
// with a new virtual network with peerings to a supplied hub network.
func TestIntegrationHubAndSpokeExistingSubscription(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["hub_network_resource_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/testvnet"
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	terraformOptions.Vars = v

	// Create plan and ensure only a single resource is created.
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 4, "expected 4 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
	}
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		// subscription variables
		"subscription_alias_billing_scope": "/providers/Microsoft.Billing/billingAccounts/test-billing-account",
		"subscription_alias_display_name":  "test-subscription-alias",
		"subscription_alias_name":          "test-subscription-alias",
		"subscription_alias_workload":      "Production",

		// virtualnetwork variables
		"virtual_network_address_space":       []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_location":            "northeurope",
		"virtual_network_name":                "testvnet",
		"virtual_network_resource_group_name": "testrg",
	}
}
