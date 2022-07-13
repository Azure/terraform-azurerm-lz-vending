package integration

// Integration tests apply to the root module and test the deployment of resources
// in specific scenarios.

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
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

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 6, "expected 6 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
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

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 5, "expected 5 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.virtualnetwork[0].azapi_resource.vhubconnection[\"vhubcon-1b4db7eb-4057-5ddf-91e0-36dec72071f5\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

// TestIntegrationSubscriptionAndRoleAssignmentOnly tests the resource plan when creating a new subscription,
// with a role assignments, but no networking.
// This tests that the depends_on property of the roleassignments module is working
// when a dependent resource is disabled through the use of count.
func TestIntegrationSubscriptionAndRoleAssignmentOnly(t *testing.T) {
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	v := getMockInputVariables()
	v["subscription_alias_enabled"] = true
	v["virtual_network_enabled"] = false
	v["role_assignment_enabled"] = true
	v["role_assignments"] = []interface{}{
		map[string]interface{}{
			"principal_id":   "00000000-0000-0000-0000-000000000000",
			"definition":     "Owner",
			"relative_scope": "",
		},
	}
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 2, "expected 2 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.subscription[0].azurerm_subscription.this[0]",
		"module.roleassignment[\"7f69efa3-575a-5f8b-a989-c3978b92b58a\"].azurerm_role_assignment.this",
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
	v["subscription_alias_enabled"] = false
	v["subscription_id"] = "00000000-0000-0000-0000-000000000000"
	v["virtual_network_enabled"] = true
	v["virtual_network_peering_enabled"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	assert.NoError(t, err)
	assert.Lenf(t, plan.ResourcePlannedValuesMap, 5, "expected 5 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.virtualnetwork[0].azapi_resource.peering[\"inbound\"]",
		"module.virtualnetwork[0].azapi_resource.peering[\"outbound\"]",
		"module.virtualnetwork[0].azapi_resource.rg",
		"module.virtualnetwork[0].azapi_resource.vnet",
		"module.virtualnetwork[0].azapi_update_resource.vnet",
	}
	for _, v := range resources {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, v)
	}
}

// TestIntegrationWithYaml tests the use of the module with a for_each loop
// using YAML files as input.
func TestIntegrationWithYaml(t *testing.T) {
	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	plan, err := terraform.InitAndPlanAndShowWithStructE(t, terraformOptions)
	require.NoErrorf(t, err, "failed to generate plan: %v", err)

	assert.Lenf(t, plan.ResourcePlannedValuesMap, 21, "expected 21 resources to be created, but got %d", len(plan.ResourcePlannedValuesMap))
	resources := []string{
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_update_resource.vnet",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_resource.vnet",
		"module.alz_landing_zone[\"%s\"].module.virtualnetwork[0].azapi_resource.rg",
		"module.alz_landing_zone[\"%s\"].module.subscription[0].azurerm_subscription.this[0]",
		"module.alz_landing_zone[\"%s\"].module.subscription[0].azurerm_management_group_subscription_association.this[0]",
		"module.alz_landing_zone[\"%s\"].module.roleassignment[\"8a6eec3e-78d9-5ff3-89cc-b144ae761a9a\"].azurerm_role_assignment.this",
		"module.alz_landing_zone[\"%s\"].module.roleassignment[\"7f69efa3-575a-5f8b-a989-c3978b92b58a\"].azurerm_role_assignment.this",
	}
	lzs := []string{
		"landing_zone_1.yaml",
		"landing_zone_2.yaml",
		"landing_zone_3.yaml",
	}
	for _, r := range resources {
		for _, lz := range lzs {
			res := fmt.Sprintf(r, lz)
			terraform.AssertPlannedValuesMapKeyExists(t, plan, res)
		}
	}
}

func getMockInputVariables() map[string]interface{} {
	return map[string]interface{}{
		// subscription variables
		"subscription_billing_scope": "/providers/Microsoft.Billing/billingAccounts/0000000/enrollmentAccounts/000000",
		"subscription_display_name":  "test-subscription-alias",
		"subscription_alias_name":    "test-subscription-alias",
		"subscription_workload":      "Production",

		// virtualnetwork variables
		"virtual_network_address_space":       []string{"10.1.0.0/24", "172.16.1.0/24"},
		"virtual_network_location":            "northeurope",
		"virtual_network_name":                "testvnet",
		"virtual_network_resource_group_name": "testrg",
	}
}
