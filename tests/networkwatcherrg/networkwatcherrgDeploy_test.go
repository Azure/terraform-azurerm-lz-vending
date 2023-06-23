package networkwatcherrg

import (
	"os"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionAliasCreateValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestDeployNetworkWatcherRg(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)

	v := getValidInputVariables()
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

// getMockInputVariables returns a set of mock input variables that can be used and modified for testing scenarios.
func getValidInputVariables() map[string]any {
	return map[string]any{
		"location":        "eastus",
		"subscription_id": os.Getenv("AZURE_SUBSCRIPTION_ID"),
	}
}
