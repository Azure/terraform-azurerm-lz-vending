package resourcegroup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionAliasCreateValid tests the validation functions with valid data,
// then creates a plan and compares the input variables to the planned values.
func TestDeployNetworkWatcherRg(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)

	v, err := getValidInputVariables()
	require.NoError(t, err)

	// delete the resource group if it already exists
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	sid, err := uuid.Parse(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	require.NoError(t, err)

	t.Logf("Getting resource groups in subscription %s", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	rgs, err := azureutils.ListResourceGroup(ctx, sid)
	require.NoError(t, err)

	for _, rg := range rgs {
		if *rg.Name == v["resource_group_name"].(string) {
			t.Logf("Deleting resource group %s", *rg.Name)
			err := azureutils.DeleteResourceGroup(ctx, *rg.Name, uuid.MustParse(os.Getenv("AZURE_SUBSCRIPTION_ID")))
			require.NoError(t, err)
		}
	}

	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(1).ErrorIsNil(t)
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables() (map[string]any, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]any{
		"subscription_id":     os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"location":            "eastus",
		"resource_group_name": name,
	}, nil
}
