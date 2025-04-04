package subscription

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/Azure/terratest-terraform-fluent/check"
	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var billingScope = os.Getenv("AZURE_BILLING_SCOPE")
var tenantId = os.Getenv("AZURE_TENANT_ID")

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
// We also test RP registration here.
// This test uses the azapi provider.
func TestDeploySubscriptionAliasValid(t *testing.T) {
	t.Parallel()

	utils.PreCheckDeployTests(t)

	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	test, err := setuptest.Dirs(moduleDir, "").WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()

	check.InPlan(test.PlanStruct).NumberOfResourcesEquals(4).ErrorIsNil(t)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		if err != nil {
			t.Logf("cannot cancel subscription: %v", err)
		}
	}()

	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	sid, err := test.Output("subscription_id").GetValue()
	assert.NoError(t, err)
	sids, ok := sid.(string)
	assert.True(t, ok, "subscription_id is not a string")
	u, err = uuid.Parse(sids)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)
}

// TestDeploySubscriptionAliasManagementGroupValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasManagementGroupValid(t *testing.T) {
	t.Parallel()
	utils.PreCheckDeployTests(t)

	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	v["subscription_billing_scope"] = billingScope
	v["subscription_management_group_id"] = v["subscription_alias_name"]
	v["subscription_management_group_association_enabled"] = true

	testDir := filepath.Join("testdata", t.Name())
	test, err := setuptest.Dirs(moduleDir, testDir).WithVars(v).InitPlanShowWithPrepFunc(t, utils.AzureRmAndRequiredProviders)
	require.NoError(t, err)
	defer test.Cleanup()
	require.NoError(t, err)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		if err != nil {
			t.Logf("cannot cancel subscription: %v", err)
		}
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer test.DestroyRetry(setuptest.DefaultRetry) //nolint:errcheck
	test.ApplyIdempotent().ErrorIsNil(t)

	sid, err := terraform.OutputE(t, test.Options, "subscription_id")
	assert.NoError(t, err)

	u, err = uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// err = azureutils.IsSubscriptionInManagementGroup(t, u, v["subscription_management_group_id"].(string))
	// assert.NoErrorf(t, err, "subscription %s is not in management group %s", sid, v["subscription_management_group_id"].(string))

	if err := azureutils.SetSubscriptionManagementGroup(u, tenantId); err != nil {
		t.Logf("cannot move subscription to tenant root group: %v", err)
	}
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]any, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]any{
		"subscription_alias_name":    name,
		"subscription_display_name":  name,
		"subscription_billing_scope": billingScope,
		"subscription_use_azapi":     false,
		"subscription_workload":      "DevTest",
		"subscription_alias_enabled": true,
	}, nil
}
