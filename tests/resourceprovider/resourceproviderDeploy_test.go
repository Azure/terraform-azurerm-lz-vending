package resourceprovider

import (
	"os"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

func TestDeploySubscriptionDeployExistingWithRpFeatureRegistration(t *testing.T) {

	utils.PreCheckDeployTests(t)

	v := make(map[string]any)
	v["subscription_id"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
	v["resource_provider"] = "Microsoft.PowerBI"
	v["features"] = []string{"DailyPrivateLinkServicesForPowerBI"}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(2).ErrorIsNil(t)
}
