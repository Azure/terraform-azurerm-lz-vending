package subscription

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-lz-vending/tests/azureutils"
	"github.com/Azure/terraform-azurerm-lz-vending/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, "")
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()

	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	if err != nil {
		t.Fatalf("Cannot generate valid input variables, %s", err)
	}

	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		t.Logf("cannot cancel subscription: %v", err)
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)
	u, err = uuid.Parse(sid)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)
}

// TestDeploySubscriptionAliasManagementGroupValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasManagementGroupValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	testDir := "testdata/" + t.Name()
	tmp, cleanup, err := utils.CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	require.NoErrorf(t, err, "failed to copy module to temp: %v", err)
	defer cleanup()
	terraformOptions := utils.GetDefaultTerraformOptions(t, tmp)

	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	v["subscription_billing_scope"] = billingScope
	v["subscription_management_group_id"] = v["subscription_alias_name"]
	v["subscription_management_group_association_enabled"] = true
	terraformOptions.Vars = v

	require.NoErrorf(t, utils.CreateTerraformProvidersFile(tmp), "Unable to create providers.tf: %v", err)
	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	// Defer the cleanup of the subscription alias to the end of the test.
	// Should be run after the Terraform destroy.
	// We don't know the sub ID yet, so use zeros for now and then
	// update it after the apply.
	u := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	defer func() {
		err := azureutils.CancelSubscription(t, &u)
		t.Logf("cannot cancel subscription: %v", err)
	}()

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)
	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)

	u, err = uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	err = azureutils.IsSubscriptionInManagementGroup(t, u, v["subscription_management_group_id"].(string))
	assert.NoErrorf(t, err, "subscription %s is not in management group %s", sid, v["subscription_management_group_id"].(string))

	// removed as azurerm_management_group_subscription_association handles this for us
	// tid := os.Getenv("AZURE_TENANT_ID")
	// if err := setSubscriptionManagementGroup(u, tid); err != nil {
	// 	t.Logf("could not move subscription to management group %s: %s", tid, err)
	// }
}

// getValidInputVariables returns a set of valid input variables that can be used and modified for testing scenarios.
func getValidInputVariables(billingScope string) (map[string]interface{}, error) {
	r, err := utils.RandomHex(4)
	if err != nil {
		return nil, fmt.Errorf("cannot generate random hex, %s", err)
	}
	name := fmt.Sprintf("testdeploy-%s", r)
	return map[string]interface{}{
		"subscription_alias_name":    name,
		"subscription_display_name":  name,
		"subscription_billing_scope": billingScope,
		"subscription_workload":      "DevTest",
		"subscription_alias_enabled": true,
	}, nil
}
