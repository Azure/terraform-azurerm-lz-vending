package resourcegroup

import (
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

const (
	moduleDir = "../../modules/resourcegroup"
)

// TestNetworkWatcherRg tests creation of a NetworkwatcherRG resource group.
func TestNetworkWatcherRg(t *testing.T) {

	v := getMockInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	check.InPlan(test.PlanStruct).That(`azapi_resource.rg`).Key("name").HasValue("NetworkWatcherRG").ErrorIsNil(t)
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getMockInputVariables() map[string]any {
	return map[string]any{
		"location":            "eastus",
		"resource_group_name": "NetworkWatcherRG",
		"subscription_id":     "00000000-0000-0000-0000-000000000000",
	}
}
