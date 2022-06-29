package subscription

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/azureutils"
	"github.com/Azure/terraform-azurerm-alz-landing-zone/tests/utils"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/matryer/try.v1"
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

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 20*time.Second, 6)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)
	_, err = uuid.Parse(sid)
	require.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// cancel the newly created sub - removed as azurerm_subscription
	// now handles this for us
	// defer func() {
	// 	err := cancelSubscription(t, u)
	// 	terraformOptions.Logger.Logf(t, "cannot cancel subscription: %v", err)
	// }()
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

	_, err = terraform.ApplyAndIdempotentE(t, terraformOptions)
	assert.NoError(t, err)

	// defer terraform destroy, but wrap in a try.Do to retry a few times
	// due to eventual consistency of the subscription aliases API
	defer utils.TerraformDestroyWithRetry(t, terraformOptions, 30*time.Second, 10)

	sid, err := terraform.OutputE(t, terraformOptions, "subscription_id")
	assert.NoError(t, err)

	u, err := uuid.Parse(sid)
	assert.NoErrorf(t, err, "subscription id %s is not a valid uuid", sid)

	// cancel the newly created sub - removed as azurerm_subscription
	// now handles this for us
	// defer func() {
	// 	err := cancelSubscription(t, u)
	// 	terraformOptions.Logger.Logf(t, "cannot cancel subscription: %v", err)
	// }()

	err = isSubscriptionInManagementGroup(t, u, v["subscription_management_group_id"].(string))
	assert.NoError(t, err)

	// removed as azurerm_management_group_subscription_association handles this for us
	//tid := os.Getenv("AZURE_TENANT_ID")
	// if err := setSubscriptionManagementGroup(u, tid); err != nil {
	// 	t.Logf("could not move subscription to management group %s: %s", tid, err)
	// }
}

// cancelSubscription cancels the supplied Azure subscription.
// it retries a few times as the subscription api is eventually consistent.
// func cancelSubscription(t *testing.T, id uuid.UUID) error {
// 	const (
// 		max      = 4
// 		delaysec = 20
// 	)

// 	if exists, err := subscriptionExists(id); err != nil || !exists {
// 		return fmt.Errorf("subscription %s does not exist or cannot successfully check, %s", id, err)
// 	}

// 	client, err := azureutils.NewSubscriptionClient()
// 	if err != nil {
// 		return fmt.Errorf("cannot create subscription client, %s", err)
// 	}
// 	ctx := context.Background()
// 	err = try.Do(func(attempt int) (bool, error) {
// 		_, err := client.Cancel(ctx, id.String(), nil)
// 		if err != nil {
// 			t.Logf("subscription id %s cancel failed, attempt %d/%d: %v", id, attempt, max, err)
// 			time.Sleep(delaysec * time.Second)
// 		}
// 		return attempt < max, err
// 	})
// 	if err != nil {
// 		return fmt.Errorf("cannot cancel subscription %s, %v", id, err)
// 	}
// 	return nil
// }

// subscriptionExists checks if the supplied subscription exists.
func subscriptionExists(id uuid.UUID) (bool, error) {
	client, err := azureutils.NewSubscriptionsClient()
	if err != nil {
		return false, fmt.Errorf("cannot create subscriptions client, %s", err)
	}
	ctx := context.Background()
	if _, err := client.Get(ctx, id.String(), nil); err != nil {
		return false, fmt.Errorf("cannot get subscription, %s", err)
	}
	return true, nil
}

// isSubscriptionInManagementGroup returns true if the subscription is a management group.
func isSubscriptionInManagementGroup(t *testing.T, id uuid.UUID, mg string) error {
	const (
		max      = 8
		delaysec = 20
	)

	if exists, err := subscriptionExists(id); err != nil || !exists {
		return fmt.Errorf("subscription %s does not exist, or could not successfully check, %s", id, err)
	}

	client, err := azureutils.NewManagementGroupSubscriptionsClient()
	if err != nil {
		return fmt.Errorf("cannot create mg subscriptions client, %s", err)
	}

	var mgopts armmanagementgroups.ManagementGroupSubscriptionsClientGetSubscriptionOptions
	cc := "no-cache"
	mgopts.CacheControl = &cc

	err = try.Do(func(attempt int) (bool, error) {
		_, err := client.GetSubscription(context.Background(), mg, id.String(), &mgopts)
		if err != nil {
			t.Logf("failed to get subscription %s in management group %s, attempt %d/%d", id.String(), mg, attempt, max)
			time.Sleep(delaysec * time.Second)
		}
		return attempt < max, err
	})
	if err != nil {
		return fmt.Errorf("failed determine if subscription %s in management group %s: %v", id.String(), mg, err)
	}
	return nil
}

// setSubscriptionManagementGroup moves the subscription to the management group.
// func setSubscriptionManagementGroup(id uuid.UUID, mg string) error {
// 	client, err := azureutils.NewManagementGroupSubscriptionsClient()
// 	if err != nil {
// 		return fmt.Errorf("cannot create mg subscriptions client, %s", err)
// 	}
// 	cc := "no-cache"
// 	opts := armmanagementgroups.ManagementGroupSubscriptionsClientCreateOptions{
// 		CacheControl: &cc,
// 	}
// 	if _, err := client.Create(context.Background(), mg, id.String(), &opts); err != nil {
// 		return fmt.Errorf("cannot create subscription %s in management group %s, %s", id.String(), mg, err)
// 	}
// 	return nil
// }

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
