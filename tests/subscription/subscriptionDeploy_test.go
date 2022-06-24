package subscription

import (
	"os"
	"testing"
	"time"

	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeploySubscriptionAliasValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	if err != nil {
		t.Fatalf("Cannot generate valid input variables, %s", err)
	}

	terraformOptions := utils.GetDefaultTerraformOptions(t, moduleDir)
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)
	u, err := uuid.Parse(sid)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// cancel the newly created sub
	if err := cancelSubscription(t, u); err != nil {
		t.Logf("could not cancel subscription: %v", err)
	} else {
		t.Logf("subscription %s cancelled", sid)
	}
}

// TestDeploySubscriptionAliasManagementGroupValid tests the deployment of a subscription alias
// with valid input variables.
func TestDeploySubscriptionAliasManagementGroupValid(t *testing.T) {
	utils.PreCheckDeployTests(t)

	dir := utils.GetTestDir(t)
	dir += "/testdata/" + t.Name()
	terraformOptions := utils.GetDefaultTerraformOptions(t, dir)
	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
	v, err := getValidInputVariables(billingScope)
	require.NoError(t, err)
	v["subscription_alias_billing_scope"] = billingScope
	v["subscription_alias_management_group_id"] = v["subscription_alias_name"]
	terraformOptions.Vars = v

	_, err = terraform.InitAndPlanE(t, terraformOptions)
	require.NoError(t, err)

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)

	u, err := uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// cancel the newly created sub
	defer func() {
		err := cancelSubscription(t, u)
		terraformOptions.Logger.Logf(t, "cannot cancel subscription: %v", err)
	}()

	err = isSubscriptionInManagementGroup(t, u, v["subscription_alias_management_group_id"].(string))
	assert.NoError(t, err)

	tid := os.Getenv("AZURE_TENANT_ID")

	if err := setSubscriptionManagementGroup(u, tid); err != nil {
		t.Logf("could not move subscription to management group %s: %s", tid, err)
	}
}

// Creating an alias for an existing subscription is not currently supported.
// Need use case data to justify the effort in testing support.
//
// // TestDeploySubscriptionAliasExistingSubscription tests the creation
// // of a subscription alias for an existing subscription
// func TestDeploySubscriptionAliasExistingSubscription(t *testing.T) {
// 	utils.PreCheckDeployTests(t)

// 	billingScope := os.Getenv("AZURE_BILLING_SCOPE")
// 	v, err := getValidInputVariables(billingScope)
// 	if err != nil {
// 		t.Fatalf("Cannot generate valid input variables, %s", err)
// 	}

// 	existingSub, err := uuid.Parse(os.Getenv("AZURE_EXISTING_SUBSCRIPTION_ID"))
// 	if err != nil {
// 		t.Fatalf("Cannot parse AZURE_EXISTING_SUBSCRIPTION_ID as uuid, %s", err)
// 	}
// 	v["subscription_id"] = existingSub.String()

// 	terraformOptions := utils.GetDefaultTerraformOptions(v)

// 	_, err = terraform.InitAndPlanE(t, terraformOptions)
// 	require.NoError(t, err)

// 	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
// 	defer terraform.Destroy(t, terraformOptions)
// 	require.NoError(t, err)

// 	sid := terraform.Output(t, terraformOptions, "subscription_id")
// 	_, err = uuid.Parse(sid)
// 	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)
// 	// DO NOT CANCEL THIS SUBSCRIPTION
// }
